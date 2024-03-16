package api

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/cristalhq/jwt/v5"
	"github.com/gorilla/schema"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thatpix3l/collge_event_website/src/gen_sql"
	app "github.com/thatpix3l/collge_event_website/src/gen_templ"
	"github.com/thatpix3l/collge_event_website/src/utils"
	"golang.org/x/crypto/bcrypt"
)

const errUniList = "unable to get list of created universities"

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

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02 15:04:05.999999999 -0700 MST"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func (t CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

type JwtClaim struct {
	UserId int        `json:"sub"` // ID of user the token was issued for
	JwtId  string     `json:"jti"` // ID of token
	Issued CustomTime `json:"iat"` // when token was issued
}

// Information about an authenticated user
type AuthenticatedUser struct {
	Id int // database ID of user
	http.Cookie
}

// Local state scoped to a single handler
type LocalState struct {
	ResponseWriter    http.ResponseWriter
	Request           *http.Request
	AlreadyParsedForm bool
}

// Parse form if never attempted
func (l *LocalState) ParseForm() error {
	if l.AlreadyParsedForm {
		return nil
	}

	if err := l.Request.ParseForm(); err != nil {
		return utils.ErrInfo(err, "unable to parse form")
	}

	l.AlreadyParsedForm = true

	return nil
}

func (l LocalState) RespondHtml(component templ.Component, status ...int) error {

	// HTML string representation of component
	htmlStr, err := templ.ToGoHTML(l.Request.Context(), component)
	if err != nil {
		return err
	}

	// HTML content type
	l.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Use optional status code
	if len(status) > 0 {
		l.ResponseWriter.WriteHeader(status[0])
	}

	// Write HTML string to client (by default, use status code 200)
	if _, err := l.ResponseWriter.Write([]byte(htmlStr)); err != nil {
		return err
	}
	return nil
}

type RouteCallback func(LocalState) error

// Convert this package's custom handler callback signature into Go's stdlib version
func StdHttpFunc(callback RouteCallback) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {

		// All route callbacks by default will print error messages, if any
		if err := callback(LocalState{Request: req, ResponseWriter: rw}); err != nil {
			log.Println(err)
		}
	}
}

// Decoder for HTML form data that ignores unknown keys.
var decoder = func() *schema.Decoder {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	return d
}()

func optErr(err error, msgs ...string) error {
	for _, msg := range msgs {
		err = utils.ErrInfo(err, msg)
	}
	return err
}

// Run transaction; parse request for user input if necessary.
func runTx[Params any, Output any](l LocalState, tx func(context.Context, Params) (Output, error), errorInfo ...string) (Output, error) {

	// Empty inserted record
	var output Output

	// Parse form data
	if err := l.ParseForm(); err != nil {
		return output, optErr(err, errorInfo...)
	}

	// Deserialize params used to create record
	var params Params
	if err := decoder.Decode(&params, l.Request.Form); err != nil {
		return output, optErr(err, errorInfo...)
	}

	// Run transaction
	if tempOut, err := tx(l.Request.Context(), params); err != nil {

		return output, err
	} else {
		output = tempOut
	}

	// Return output from transaction
	return output, nil
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

// *Handlers and global state shared between all handlers
type Handlers struct {
	Pool *pgxpool.Pool
}

// Using given HTTP method, attach callback to given path
func (h *Handlers) Add(method func(string, http.HandlerFunc), path string, callback RouteCallback) {
	method(path, StdHttpFunc(callback))
}

// Authentication middleware
func (h *Handlers) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(StdHttpFunc(func(l LocalState) error {

		next.ServeHTTP(l.ResponseWriter, l.Request)
		return nil

		// allow if accessing any resources that don't require authentication or authorization
		for _, path := range noAuthPaths {
			if l.Request.URL.Path == path {
				next.ServeHTTP(l.ResponseWriter, l.Request)
				return nil
			}
		}

		// Should only be here if requested resource requires authentication and authorization

		// Exit early if no authentication token provided
		givenToken, err := l.Request.Cookie("authentication_token")
		if err == http.ErrNoCookie {

			// alert user that they are not authenticated yet
			l.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			if _, err := l.ResponseWriter.Write([]byte("401 - Any access besides homepage requires authentication")); err != nil {
				return err
			}
			return err

		}

		// Exit Early if given token is invalid
		if !validToken(givenToken.Value) {
			if err := l.RespondHtml(app.StatusMessage("danger", "403 - Not authorized to access resource"), http.StatusForbidden); err != nil {
				return err
			}
		}

		// Token exists and is valid, continue
		next.ServeHTTP(l.ResponseWriter, l.Request)

		return nil

	}))
}

// Get homepage; depends on provided token credentials.
func (h *Handlers) ReadHome(l LocalState) error {

	comp := app.LoginForm()

	// If authenticated and authorized, allow access the list of events
	if authenticated(l.Request) {

		// Acquire database connection
		conn, err := h.Pool.Acquire(l.Request.Context())
		if err != nil {
			return utils.ErrDb(err)
		}
		defer conn.Release()
		queries := gen_sql.New(conn)

		// Get list of universities
		universities, err := runTx(l, noParamTx(queries.ReadUniversities), errUniList)
		if err != nil {
			return err
		}

		// Set as component to send
		comp = app.CreatedUniversities(universities)

	}

	// Respond to client with fully rendered home
	if err := l.RespondHtml(app.Home(comp)); err != nil {
		return err
	}

	return nil

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
		Issued: CustomTime{time.Now()},
	})

	// Return token and possible error
	return token, err

}

// Createa new JWT token cookie
func newTokenCookie(token *jwt.Token) http.Cookie {
	tokenCookie := http.Cookie{
		Name:     "authenticationToken",
		Value:    "Bearer " + token.String(),
		HttpOnly: true,
	}
	return tokenCookie
}

func (l *LocalState) FormGet(key string) (string, error) {

	var val string

	// Parse form, exit on error
	if err := l.ParseForm(); err != nil {
		return "", err
	}

	// Get value from form with key, exit if not provided
	val = l.Request.Form.Get(key)
	if val == "" {
		return val, errors.New("form with provided key has no value")
	}

	return val, nil
}

// Create login session based on provided form credentials.
func (h *Handlers) CreateLogin(l LocalState) error {

	// Acquire database connection
	conn, err := h.Pool.Acquire(l.Request.Context())
	if err != nil {
		l.RespondHtml(app.StatusMessage("warning", err.Error()), http.StatusInternalServerError)
		return utils.ErrDb(err)
	}
	defer conn.Release()
	queries := gen_sql.New(conn)

	// Retrieve email from user
	email, err := l.FormGet("Email")
	if err != nil {
		return utils.ErrInfo(err, "unable to get Email")
	}

	// Retrieve plaintext password from user
	passwordPlaintext, err := l.FormGet("PasswordPlaintext")
	if err != nil {
		return utils.ErrInfo(err, "unable to get PasswordPlaintext")
	}

	// Get list of users
	users, err := runTx(l, noParamTx(queries.ReadStudents), "unable to get list of students")
	if err != nil {
		return err
	}

	// Check if user with email exists in database
	userExists := false
	var baseUser gen_sql.ReadStudentsRow
	for _, user := range users {
		if user.Email == email {
			userExists = true
			baseUser = user
			break
		}
	}

	if !userExists {
		l.RespondHtml(app.StatusMessage("danger", "unable to find user with email/password combination"), http.StatusInternalServerError)
	}

	// Check if provided password matches email
	if bcrypt.CompareHashAndPassword([]byte(baseUser.PasswordHash), []byte(passwordPlaintext)) == nil {
		return errors.New("provided password does not match stored password")
	}

	// Create new authentication token
	token, err := newToken(int(baseUser.ID))
	if err != nil {
		l.RespondHtml(app.StatusMessage("warning", err.Error()), http.StatusInternalServerError)
	}
	tokenCookie := newTokenCookie(token)

	// Get list of universities
	universities, err := runTx(l, noParamTx(queries.ReadUniversities), errUniList)
	if err != nil {
		l.RespondHtml(app.StatusMessage("warning", err.Error()), http.StatusInternalServerError)
		return err
	}

	// Store cookie into Set-Cookie header; client should respond back with the given token
	l.ResponseWriter.Header().Set("Set-Cookie", tokenCookie.String())
	l.RespondHtml(app.CreatedUniversities(universities))

	return nil

}

// Get login form used to create a login session.
func (h *Handlers) ReadLogin(l LocalState) error {
	l.RespondHtml(app.LoginForm())
	return nil
}

// Get signup form used to create a student account.
func (h *Handlers) ReadSignup(l LocalState) error {

	// Acquire database connection
	conn, err := h.Pool.Acquire(l.Request.Context())
	if err != nil {
		return utils.ErrDb(err)
	}
	defer conn.Release()
	queries := gen_sql.New(conn)

	// Get list of created universities
	universities, err := runTx(l, noParamTx(queries.ReadUniversities), errUniList)
	if err != nil {
		return err
	}

	// Respond with HTML
	l.RespondHtml(app.SignupForm(universities))
	return nil
}

// Hash plaintext password and store back into form for storing into database
func (l LocalState) HashPasswordInput() error {

	// Parse form data if needed
	l.ParseForm()

	passwordRaw, err := l.FormGet("PasswordPlaintext")
	if err != nil {
		return utils.ErrInfo(err, "unable to get PasswordPlaintext")
	}

	passwordHashBuf, err := bcrypt.GenerateFromPassword([]byte(passwordRaw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	l.Request.Form["PasswordHash"] = []string{string(passwordHashBuf)}

	return nil

}

// Create new student that's associated with a university.
func (h *Handlers) CreateStudent(l LocalState) error {

	// Parse form of request
	if err := l.ParseForm(); err != nil {
		return err
	}

	if err := l.HashPasswordInput(); err != nil {
		return err
	}

	// Acquire database connection
	conn, err := h.Pool.Acquire(l.Request.Context())
	if err != nil {
		return utils.ErrDb(err)
	}
	defer conn.Release()
	queries := gen_sql.New(conn)

	// Create student
	if _, err := runTx(l, queries.CreateStudent); err != nil {
		return err
	}

	return nil

}

// Get list of students
func (h *Handlers) ReadStudents(l LocalState) error {

	// Acquire database connection
	conn, err := h.Pool.Acquire(l.Request.Context())
	if err != nil {
		return utils.ErrDb(err)
	}
	defer conn.Release()
	queries := gen_sql.New(conn)

	students, err := runTx(l, noParamTx(queries.ReadStudents))
	if err != nil {
		return err
	}

	l.RespondHtml(app.CreatedStudent(students[0]))

	return nil

}

// Get list of universities.
func (h *Handlers) ReadUniversities(l LocalState) error {

	// Acquire database connection
	conn, err := h.Pool.Acquire(l.Request.Context())
	if err != nil {
		return utils.ErrDb(err)
	}
	defer conn.Release()
	queries := gen_sql.New(conn)

	// Read list of created universities
	universities, err := queries.ReadUniversities(l.Request.Context())
	if err != nil {
		return err
	}

	// Respond with universities
	if err := l.RespondHtml(app.CreatedUniversities(universities)); err != nil {
		return err
	}

	return nil

}

// Create a new university record.
func (h *Handlers) CreateUniversity(l LocalState) error {

	// Acquire database connection
	conn, err := h.Pool.Acquire(l.Request.Context())
	if err != nil {
		return utils.ErrDb(err)
	}
	defer conn.Release()
	queries := gen_sql.New(conn)

	// Create new university
	if _, err := runTx(l, queries.CreateUniversity); err != nil {

		l.RespondHtml(app.StatusMessage("danger", "Unable to create university"), http.StatusInternalServerError)
		return err
	}

	// Respond with status
	if err := l.RespondHtml(app.StatusMessage("success", "Created new university!")); err != nil {
		return err
	}

	return nil

}
