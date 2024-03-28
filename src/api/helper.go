package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thatpix3l/collge_event_website/src/gen_sql"
	"github.com/thatpix3l/collge_event_website/src/utils"
	"golang.org/x/crypto/bcrypt"
)

var tokenSecret = []byte(uuid.NewString())

// State accessible by all handlers.
type globalState struct {
	Pool *pgxpool.Pool
}

// API package's access to the global state.
var GlobalState = globalState{}

// State accessible by a single handler.
type LocalState struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	ParsedForm     bool             // already parsed request's form data?
	Conn           *pgxpool.Conn    // database connection.
	Queries        *gen_sql.Queries // queries connection.
}

// Parse form if never attempted.
func (ls *LocalState) ParseForm() error {
	if ls.ParsedForm {
		return nil
	}

	if err := ls.Request.ParseForm(); err != nil {
		return utils.ErrPrep(err, "unable to parse form")
	}

	ls.ParsedForm = true

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
func (ls *LocalState) HashPasswordInput() error {

	// Parse form data if needed.
	ls.ParseForm()

	passwordRaw, err := ls.FormGet("PasswordPlaintext")
	if err != nil {
		return utils.ErrPrep(err, "unable to get PasswordPlaintext")
	}

	passwordHashBuf, err := bcrypt.GenerateFromPassword([]byte(passwordRaw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	ls.Request.Form["PasswordHash"] = []string{string(passwordHashBuf)}

	return nil

}

// Current state for each handler.
type HandlerState struct {
	Global *globalState
	Local  *LocalState
}

// Create database connection if not already exist.
func (hs *HandlerState) Conn() error {
	if hs.Local.Conn != nil {
		return nil
	}

	// Acquire connection.
	conn, err := hs.Global.Pool.Acquire(hs.Local.Request.Context())
	if err != nil {
		return errors.New("unable to acquire database connection")
	}

	// Store connection.
	hs.Local.Conn = conn

	return nil
}

// Create queries connection if not already exist.
func (hs *HandlerState) Queries() error {

	if hs.Local.Queries != nil {
		return nil
	}

	if err := hs.Conn(); err != nil {
		return err
	}

	hs.Local.Queries = gen_sql.New(hs.Local.Conn)

	return nil
}

type HandlerFunc func(hs HandlerState) error

// Convert this package's custom handler function signature into Go's stdlib version.
func StdHttpFunc(path string, method string, callback HandlerFunc) func(http.ResponseWriter, *http.Request) {
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
		if err := callback(hs); err != nil {
			log.Println(utils.ErrPrep(err, path, method))
		}

		// Perform cleanup.
		if hs.Local.Conn != nil {
			hs.Local.Conn.Release()
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

// Run query; parse request for user input if necessary.
func runQuery[Params any, Output any](hs HandlerState, query func(context.Context, Params) (Output, error), errorInfo ...string) (Output, error) {

	// Output from running query.
	var output Output

	// Parse form data.
	if err := hs.Local.ParseForm(); err != nil {
		return output, utils.ErrPrep(err, errorInfo...)
	}

	// Deserialize params needed by query.
	var params Params
	if err := decoder.Decode(&params, hs.Local.Request.Form); err != nil {
		return output, utils.ErrPrep(err, errorInfo...)
	}

	// Run query.
	if tempOut, err := query(hs.Local.Request.Context(), params); err != nil {

		return output, utils.ErrPrep(err, errorInfo...)
	} else {
		output = tempOut
	}

	// Return output from transaction.
	return output, nil
}

type AuthenticatedUsers struct {
	list []jwt.RegisteredClaims
	lock sync.RWMutex
}

func (au *AuthenticatedUsers) Add(claims jwt.RegisteredClaims) {
	au.lock.Lock()
	au.list = append(au.list, claims)
	au.lock.Unlock()
}

var authenticatedUsers = AuthenticatedUsers{
	list: []jwt.RegisteredClaims{},
}

func tokenParser(t *jwt.Token) (interface{}, error) {
	return tokenSecret, nil
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
}

// Wrapper that converts a query that doesn't accept concrete parameters into one that accepts an empty struct.
func noParam[Output any](query func(context.Context) (Output, error)) func(context.Context, struct{}) (Output, error) {
	return func(ctx context.Context, p struct{}) (Output, error) {
		return query(ctx)
	}
}

func (ls *LocalState) FormGet(key string) (string, error) {

	var val string

	// Parse form, exit on error.
	if err := ls.ParseForm(); err != nil {
		return "", err
	}

	// Get value, exit if empty.
	val = ls.Request.Form.Get(key)
	if val == "" {
		return val, errors.New("form with provided key has no value")
	}

	return val, nil
}
