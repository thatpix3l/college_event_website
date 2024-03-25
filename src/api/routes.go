package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/thatpix3l/collge_event_website/src/gen_sql"
	app "github.com/thatpix3l/collge_event_website/src/gen_templ"
	"github.com/thatpix3l/collge_event_website/src/utils"
	"golang.org/x/crypto/bcrypt"
)

// Get homepage.
var ReadHomepage = addHandlerFunc("/", "get", func(hs HandlerState) error {

	comp := app.LoginForm()

	// If authenticated and authorized, allow access the list of events
	if authenticated(hs.Local.Request) {

		// Create queries connection
		if err := hs.Queries(); err != nil {
			return err
		}

		// Get list of universities
		universities, err := runTx(hs.Local, noParamTx(hs.Local.Queries.ReadUniversities))
		if err != nil {
			return err
		}

		// Set as component to send
		comp = app.CreatedUniversities(universities)

	}

	// Respond to client with fully rendered home
	if err := hs.Local.RespondHtml(app.Home(comp)); err != nil {
		return err
	}

	return nil

})

// Create login session based on provided form credentials.
var CreateLogin = addHandlerFunc(utils.ApiPath("login"), "post", func(hs HandlerState) error {

	// Acquire queries connection
	if err := hs.Queries(); err != nil {
		return err
	}

	// Retrieve email from user
	email, err := hs.Local.FormGet("Email")
	if err != nil {
		return utils.ErrPrep(err, "unable to get Email")
	}

	// Retrieve plaintext password from user
	passwordPlaintext, err := hs.Local.FormGet("PasswordPlaintext")
	if err != nil {
		return utils.ErrPrep(err, "unable to get PasswordPlaintext")
	}

	// Get list of users
	users, err := runTx(hs.Local, noParamTx(hs.Local.Queries.ReadStudents), "unable to get list of students")
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

	// Check if user with provided email exists
	if !userExists {
		hs.Local.RespondHtml(app.StatusMessage("danger", "unable to find user with email/password combination"), http.StatusInternalServerError)
		return errors.New("user with provided email does not exist")
	}

	// Check if provided password matches email
	if err := bcrypt.CompareHashAndPassword([]byte(baseUser.PasswordHash), []byte(passwordPlaintext)); err != nil {
		hs.Local.RespondHtml(app.StatusMessage("danger", "unable to find user with email/password combination"), http.StatusInternalServerError)
		return utils.ErrPrep(err, "password does not match user with provided email")
	}

	// Create new authentication token
	token, err := newToken(int(baseUser.ID))
	if err != nil {
		hs.Local.RespondHtml(app.StatusMessage("warning", "unable to create JWT token"), http.StatusInternalServerError)
		return err
	}
	tokenCookie := newTokenCookie(token)

	// Get list of universities
	universities, err := runTx(hs.Local, noParamTx(hs.Local.Queries.ReadUniversities))
	if err != nil {
		hs.Local.RespondHtml(app.StatusMessage("warning", err.Error()), http.StatusInternalServerError)
		return err
	}

	// Store cookie into Set-Cookie header; client should respond back with the given token
	hs.Local.ResponseWriter.Header().Set("Set-Cookie", tokenCookie.String())
	hs.Local.RespondHtml(app.CreatedUniversities(universities))

	return nil

})

// Get login form used to create a login session.
var ReadLogin = addHandlerFunc(utils.ApiPath("login"), "get", func(hs HandlerState) error {
	hs.Local.RespondHtml(app.LoginForm())
	return nil
})

// Get signup form used to create a student account.
var ReadSignup = addHandlerFunc(utils.ApiPath("signup"), "get", func(hs HandlerState) error {

	// Acquire queries connection
	if err := hs.Queries(); err != nil {
		return err
	}

	// Get list of created universities
	universities, err := runTx(hs.Local, noParamTx(hs.Local.Queries.ReadUniversities))
	if err != nil {
		return err
	}

	// Respond with HTML
	hs.Local.RespondHtml(app.SignupForm(universities))
	return nil
})

// Create new student that's associated with a university.
var CreateStudent = addHandlerFunc(utils.ApiPath("signup"), "post", func(hs HandlerState) error {

	// Parse form of request
	if err := hs.Local.ParseForm(); err != nil {
		hs.Local.RespondHtml(app.StatusMessage("danger", err.Error()), http.StatusInternalServerError)
		return err
	}

	// Hash user-provided password input
	if err := hs.Local.HashPasswordInput(); err != nil {
		hs.Local.RespondHtml(app.StatusMessage("danger", err.Error()), http.StatusInternalServerError)
		return err
	}

	// Acquire queries connection
	if err := hs.Queries(); err != nil {
		hs.Local.RespondHtml(app.StatusMessage("danger", err.Error()), http.StatusInternalServerError)
		return err
	}

	// Create student
	if _, err := runTx(hs.Local, hs.Local.Queries.CreateStudent); err != nil {
		hs.Local.RespondHtml(app.StatusMessage("danger", err.Error()), http.StatusBadRequest)
		return err
	}

	return nil

})

// Get list of students
var ReadStudents = addHandlerFunc(utils.ApiPath("student"), "get", func(hs HandlerState) error {

	if err := hs.Queries(); err != nil {
		return err
	}

	students, err := runTx(hs.Local, noParamTx(hs.Local.Queries.ReadStudents))
	if err != nil {
		return err
	}

	hs.Local.RespondHtml(app.CreatedStudents(students))
	log.Println(students)

	return nil

})

// Get list of universities.
var ReadUniversities = addHandlerFunc(utils.ApiPath("university"), "get", func(hs HandlerState) error {

	if err := hs.Queries(); err != nil {
		return err
	}

	// Read list of created universities
	universities, err := hs.Local.Queries.ReadUniversities(hs.Local.Request.Context())
	if err != nil {
		return err
	}

	// Respond with universities
	if err := hs.Local.RespondHtml(app.CreatedUniversities(universities)); err != nil {
		return err
	}

	return nil

})

// Create a new university record.
var CreateUniversity = addHandlerFunc(utils.ApiPath("university"), "post", func(hs HandlerState) error {

	// Get queries connection to database
	if err := hs.Queries(); err != nil {
		return err
	}

	// Create new university
	if _, err := runTx(hs.Local, hs.Local.Queries.CreateUniversity); err != nil {

		hs.Local.RespondHtml(app.StatusMessage("danger", "Unable to create university"), http.StatusInternalServerError)
		return err
	}

	// Respond with status
	if err := hs.Local.RespondHtml(app.StatusMessage("success", "Created new university!")); err != nil {
		return err
	}

	return nil

})
