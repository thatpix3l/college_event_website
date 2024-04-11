package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/thatpix3l/cew/src/gen_sql"
	"github.com/thatpix3l/cew/src/utils"
	"golang.org/x/crypto/bcrypt"

	pg "github.com/go-jet/jet/v2/postgres"
	t "github.com/thatpix3l/cew/src/gen_sql/college_event_website/cew/table"
)

var tokenSecret = []byte(uuid.NewString())

// State accessible by all handlers.
type globalState struct {
	Pool *pgxpool.Pool
	Db   *sql.DB
}

// API package's access to the global state.
var GlobalState = globalState{}

// State accessible by a single handler.
type LocalState struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	ParsedForm     bool // already parsed request's form data?
}

// Parse form if never attempted.
func (hs HandlerState) ParseForm() error {
	if hs.Local.ParsedForm {
		return nil
	}

	if err := hs.Local.Request.ParseForm(); err != nil {
		return utils.ErrPrep(err, "unable to parse form")
	}

	hs.Local.ParsedForm = true

	return nil
}

func (ls LocalState) RespondHtml(component templ.Component, status ...int) error {

	// HTML strings representation of component.
	htmlStr, err := templ.ToGoHTML(ls.Request.Context(), component)
	if err != nil {
		return err
	}

	// HTML content type.
	ls.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Use optional status code if not empty.
	if len(status) > 0 {
		ls.ResponseWriter.WriteHeader(status[0])
	}

	// Write HTML string to client (by default, use status code 200)
	if _, err := ls.ResponseWriter.Write([]byte(htmlStr)); err != nil {
		return err
	}
	return nil
}

// Hash plaintext password and store back into form for later usage.
func (hs HandlerState) HashPasswordInput() error {

	// Parse form data if needed.
	hs.ParseForm()

	// Get plaintext password
	passwordRaw, err := hs.FormGet("PasswordPlaintext")
	if err != nil {
		return utils.ErrPrep(err, "unable to get PasswordPlaintext")
	}

	// Hash password
	passwordHashBuf, err := bcrypt.GenerateFromPassword([]byte(passwordRaw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Store back into form data
	hs.Local.Request.Form["PasswordHash"] = []string{string(passwordHashBuf)}

	return nil

}

// Current state for each handler.
type HandlerState struct {
	Global *globalState
	Local  *LocalState
}

type HandlerFunc func(hs HandlerState) error

type HandlerFuncMiddleware func(hs HandlerState, next http.Handler) error

// Convert this package's custom handler function signature into Go's stdlib version.
func StdHttpFunc(path string, method string, handler HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {

		// Create handler state.
		hs := HandlerState{
			Global: &GlobalState,
			Local: &LocalState{
				Request:        req,
				ResponseWriter: rw,
			},
		}

		// Run handler.
		if err := handler(hs); err != nil {
			hs.Local.ResponseWriter.WriteHeader(http.StatusInternalServerError)
			log.Println(utils.ErrPrep(err, "path "+path+", method "+method))
		}

	}
}

// List of configured generic handle funcs to be used in any generic router.
var HandleFuncs = make(map[string]map[string]func(http.ResponseWriter, *http.Request))

// Add given handler func to list of configured handler funcs to be used in any generic router.
func addHandlerFunc(path string, method string, fn HandlerFunc) error {

	method = strings.ToUpper(method)

	if HandleFuncs[path] == nil {
		HandleFuncs[path] = map[string]func(http.ResponseWriter, *http.Request){}
	}

	HandleFuncs[path][method] = StdHttpFunc(path, method, fn)

	return nil
}

// Decoder for HTML form data that ignores unknown keys.
var decoder = func() *schema.Decoder {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	return d
}()

// Run Jet SQL statement; store in output pointer if not nil.
func runQuery(hs HandlerState, stmt pg.Statement, output any) error {

	if output == nil {
		// If no destination, execute query and don't store result
		if _, err := stmt.Exec(hs.Global.Db); err != nil {
			return err
		}

	} else {
		// Otherwise, execute query and store result
		if err := stmt.QueryContext(hs.Local.Request.Context(), hs.Global.Db, output); err != nil {
			return err
		}
	}

	return nil
}

// Copy all required fields from the HandlerState's Form to destination struct
func (hs HandlerState) ToParams(dest any) error {

	// Build form if not already exists
	if err := hs.ParseForm(); err != nil {
		return err
	}

	// Verify that each field from the destination struct exists in the flattened Form
	destStruct := reflect.Indirect(reflect.ValueOf(dest))
	for i := 0; i < destStruct.NumField(); i++ {

		field := destStruct.Type().Field(i) // struct field
		fieldName := field.Name             // struct field name

		// Check if the field is a primary key
		tags := strings.Split(field.Tag.Get("sql"), ",")
		isPrimaryKey := false
		for _, tag := range tags {
			if tag == "primary_key" {
				isPrimaryKey = true
				break
			}
		}

		// Error if field does not exist from Form and is not a primary key.
		if _, ok := hs.Local.Request.Form[fieldName]; !ok && !isPrimaryKey {
			return fmt.Errorf("form to SQL params: form is missing key \"%s\"", fieldName)
		}
	}

	if err := decoder.Decode(dest, hs.Local.Request.Form); err != nil {
		return err
	}

	return nil
}

type AuthenticatedUsers struct {
	list []jwt.RegisteredClaims
	lock sync.RWMutex
}

// Add claims for a logged in user into runtime cache
func (au *AuthenticatedUsers) Add(claims jwt.RegisteredClaims) {
	au.lock.Lock()
	au.list = append(au.list, claims)
	au.lock.Unlock()
}

// Runtime cache of authenticated users
var authenticatedUsers = AuthenticatedUsers{
	list: []jwt.RegisteredClaims{},
}

func tokenParser(t *jwt.Token) (interface{}, error) {
	return tokenSecret, nil
}

func (hs HandlerState) GetClaims(claims *jwt.RegisteredClaims) error {

	if claims == nil {
		return errors.New("claim pointer is nil")
	}

	// Get signed auth token from cookies.
	c, err := hs.Local.Request.Cookie("authenticationToken")
	if err != nil {
		return err
	}

	// Parse claims from signed token.
	parsedClaims := jwt.RegisteredClaims{}
	if _, err := jwt.ParseWithClaims(c.Value, &parsedClaims, tokenParser); err != nil {
		return err
	}

	// Store into target
	*claims = parsedClaims

	return nil

}

func (hs HandlerState) GetUser(user *User) error {

	if user == nil {
		return errors.New("user pointer is nil")
	}

	// Get claims from handler
	claims := jwt.RegisteredClaims{}
	if err := hs.GetClaims(&claims); err != nil {
		return err
	}

	// Query that gets users that match claim's Subject, which should contain a user ID
	query := ReadUsers().WHERE(t.Baseuser.ID.EQ(pg.String(claims.Subject)))

	// Run query, store list of users.
	users := []User{}
	if err := runQuery(hs, query, &users); err != nil {
		return err
	}

	// A length of 0 means no user with ID exists
	if len(users) == 0 {
		return errors.New("unable to find user with ID")
	}

	// Only get first user; database enforces primary key, so there should only ever be at most one matching user.
	*user = users[0]

	return nil

}

// Check if given authentication token is still valid.
func (hs HandlerState) Authenticated() error {

	// Get signed auth token from cookies.
	c, err := hs.Local.Request.Cookie("authenticationToken")
	if err != nil {
		return err
	}

	// Parse claims from signed token.
	parsedClaims := jwt.RegisteredClaims{}
	if _, err := jwt.ParseWithClaims(c.Value, &parsedClaims, tokenParser); err != nil {
		return err
	}

	authenticatedUsers.lock.RLock()
	defer authenticatedUsers.lock.RUnlock()

	// Attempt to check if token has been cached and still valid.
	for _, cachedClaims := range authenticatedUsers.list {

		sub, err := cachedClaims.GetSubject()
		if err != nil {
			continue
		}

		// If auth token is cached and used between Expiration and NotBefore timeframe, allow.
		now := time.Now()
		if parsedClaims.Subject == sub && now.After(parsedClaims.NotBefore.Time) && parsedClaims.ExpiresAt.After(now) {
			return nil
		}
	}

	return errors.New("provided authentication token is invalid")

}

// Accessible paths with associated methods that don't require authentication.
var noAuth = map[string][]string{
	"/":                     {"get"},
	utils.ApiPath("login"):  {"get", "post"},
	utils.ApiPath("signup"): {"get", "post"},
	utils.ApiPath("init"):   {"post"},
}

// Get value from form; string cannot be empty
func (hs HandlerState) FormGet(key string) (string, error) {

	val, err := hs.FormGetOpt(key)
	if err != nil {
		return "", err
	}

	if val == "" {
		return "", errors.New("form with provided key has no value")
	}

	return val, nil

}

// Get value from form; string can be empty
func (hs HandlerState) FormGetOpt(key string) (string, error) {

	var val string

	// Parse form, exit on error
	if err := hs.ParseForm(); err != nil {
		return val, err
	}

	// Get value, exit if empty
	val = hs.Local.Request.Form.Get(key)

	return val, nil

}

func (hs HandlerState) Authenticate(user User) error {

	// Create claims.
	now := jwt.NumericDate{Time: time.Now()}
	expires := jwt.NumericDate{Time: now.Add(time.Hour * 24 * 3)}

	claims := jwt.RegisteredClaims{
		Issuer:    "college_event_website",
		Subject:   user.Baseuser.ID,
		Audience:  jwt.ClaimStrings{"user"},
		ExpiresAt: &expires,
		NotBefore: &now,
		IssuedAt:  &now,
		ID:        user.Baseuser.ID,
	}

	// Cache claims.
	authenticatedUsers.Add(claims)

	// Create signed token string.
	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(tokenSecret)
	if err != nil {
		return err
	}

	// Create cookie storing signed token string.
	authCookie := http.Cookie{
		Name:     "authenticationToken",
		Value:    ss,
		Path:     "/",
		Expires:  expires.Time,
		HttpOnly: true,
	}

	// Store cookie into Set-Cookie header for future HTTP access
	http.SetCookie(hs.Local.ResponseWriter, &authCookie)

	// Also copy cookie into Request for later usage in the same handler
	hs.Local.Request.Header.Set("Cookie", authCookie.Name+"="+authCookie.Value)

	return nil
}
