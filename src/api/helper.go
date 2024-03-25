package api

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/cristalhq/jwt/v5"
	"github.com/gorilla/schema"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thatpix3l/collge_event_website/src/gen_sql"
	"github.com/thatpix3l/collge_event_website/src/utils"
	"golang.org/x/crypto/bcrypt"
)

var tokenSecret = func() []byte {
	b := make([]byte, 256)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}()

var tokenSigner = func() *jwt.HSAlg {
	signer, err := jwt.NewSignerHS(jwt.HS256, tokenSecret)
	if err != nil {
		panic(err)
	}

	return signer
}()

var tokenBuilder = jwt.NewBuilder(tokenSigner)

type JwtTime struct {
	time.Time
}

func (t *JwtTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02 15:04:05.999999999 -0700 MST"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func (t JwtTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.String() + "\""), nil
}

type JwtClaim struct {
	UserId int     `json:"sub"` // ID of user the token was issued for
	JwtId  string  `json:"jti"` // ID of token
	Issued JwtTime `json:"iat"` // when token was issued
}

// State accessible by all handlers
type globalState struct {
	Pool *pgxpool.Pool
}

// API package's access to the global state
var GlobalState = globalState{}

// State accessible by a single handler
type LocalState struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	ParsedForm     bool             // already parsed request's form data?
	Conn           *pgxpool.Conn    // database connection
	Queries        *gen_sql.Queries // queries connection
}

// Parse form if never attempted
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

	// HTML strings representation of component
	htmlStr, err := templ.ToGoHTML(ls.Request.Context(), component)
	if err != nil {
		return err
	}

	// HTML content type
	ls.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Use optional status code if not empty
	if len(status) > 0 {
		ls.ResponseWriter.WriteHeader(status[0])
	}

	// Write HTML string to client (by default, use status code 200)
	if _, err := ls.ResponseWriter.Write([]byte(htmlStr)); err != nil {
		return err
	}
	return nil
}

// Hash plaintext password and store back into form for later usage
func (ls *LocalState) HashPasswordInput() error {

	// Parse form data if needed
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

// Current state for each handler
type HandlerState struct {
	Global *globalState
	Local  *LocalState
}

// Create database connection if not already exist
func (hs *HandlerState) Conn() error {
	if hs.Local.Conn != nil {
		return nil
	}

	// Acquire connection
	conn, err := hs.Global.Pool.Acquire(hs.Local.Request.Context())
	if err != nil {
		return errors.New("unable to acquire database connection")
	}

	// Store connection
	hs.Local.Conn = conn

	return nil
}

// Create queries connection if not already exist
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

// Convert this package's custom handler function signature into Go's stdlib version
func StdHttpFunc(path string, method string, callback HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {

		// Create handler state
		hs := HandlerState{
			Global: &GlobalState,
			Local: &LocalState{
				Request:        req,
				ResponseWriter: rw,
			},
		}

		// Run handler
		if err := callback(hs); err != nil {
			log.Println(utils.ErrPrep(err, path, method))
		}

		// Perform cleanup
		if hs.Local.Conn != nil {
			hs.Local.Conn.Release()
		}

	}
}

// List of configured generic handle funcs to be used in any generic router
var HandleFuncs = make(map[string]map[string]func(http.ResponseWriter, *http.Request))

// Add given handler func to list of configured handler funcs to be used in any generic router
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

// Run transaction; parse request for user input if necessary.
func runTx[Params any, Output any](ls *LocalState, tx func(context.Context, Params) (Output, error), errorInfo ...string) (Output, error) {

	// Empty inserted record
	var output Output

	// Parse form data
	if err := ls.ParseForm(); err != nil {
		return output, utils.ErrPrep(err, errorInfo...)
	}

	// Deserialize params used to create record
	var params Params
	if err := decoder.Decode(&params, ls.Request.Form); err != nil {
		return output, utils.ErrPrep(err, errorInfo...)
	}

	// Run transaction
	if tempOut, err := tx(ls.Request.Context(), params); err != nil {

		return output, err
	} else {
		output = tempOut
	}

	// Return output from transaction
	return output, nil
}

// Information about an authenticated user
type AuthenticatedUser struct {
	Id int // database ID of user
	http.Cookie
}

var authenticatedUsers map[string]AuthenticatedUser = make(map[string]AuthenticatedUser) // Authenticated users cache
var authenticatedUsersLock sync.RWMutex                                                  // Authenticated users cache lock

// Check if given authentication token is still valid
func validToken(authToken string) bool {
	authenticatedUsersLock.RLock()
	defer authenticatedUsersLock.RUnlock()
	for authenticated := range authenticatedUsers {
		if authToken == authenticated {
			return true
		}
	}

	return false
}

var noAuthPaths = []string{"/", utils.ApiPath("login"), utils.ApiPath("signup")}

func authenticated(req *http.Request) bool {

	c, err := req.Cookie("authentication_token")
	if err != nil {
		return false
	}

	for _, authUser := range authenticatedUsers {
		if c.Value == authUser.Cookie.Value {
			return true
		}
	}

	return false

}

// Wrapper for generic transaction that does NOT accept any parameters.
func noParamTx[Output any, Param struct{}](tx func(context.Context) (Output, error)) func(context.Context, Param) (Output, error) {
	return func(ctx context.Context, p Param) (Output, error) {
		out, err := tx(ctx)
		return out, err
	}
}

// Create new JWT token
func newToken(user int) (*jwt.Token, error) {

	// Create unique JWT token id
	jwtId := make([]byte, 256)
	if _, err := rand.Read(jwtId); err != nil {
		return &jwt.Token{}, err
	}

	// Build token that user can reuse for continued access and refreshes
	token, err := tokenBuilder.Build(JwtClaim{
		UserId: int(user),
		JwtId:  string(jwtId),
		Issued: JwtTime{time.Now()},
	})

	// Return token and possible error
	return token, err

}

// Create new JWT token cookie
func newTokenCookie(token *jwt.Token) http.Cookie {
	tokenCookie := http.Cookie{
		Name:     "authenticationToken",
		Value:    "Bearer " + token.String(),
		HttpOnly: true,
	}
	return tokenCookie
}

func (ls *LocalState) FormGet(key string) (string, error) {

	var val string

	// Parse form, exit on error
	if err := ls.ParseForm(); err != nil {
		return "", err
	}

	// get value, exit if empty.
	val = ls.Request.Form.Get(key)
	if val == "" {
		return val, errors.New("form with provided key has no value")
	}

	return val, nil
}
