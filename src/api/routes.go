package api

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	. "github.com/thatpix3l/collge_event_website/src/gen_sql"
	"github.com/thatpix3l/collge_event_website/src/utils"
	"golang.org/x/crypto/bcrypt"

	app "github.com/thatpix3l/collge_event_website/src/gen_templ"

	m "github.com/thatpix3l/collge_event_website/src/gen_sql/college_event_website/cew/model"
	t "github.com/thatpix3l/collge_event_website/src/gen_sql/college_event_website/cew/table"
)

// Get homepage.
var ReadHomepageErr = addHandlerFunc("/", "get", func(hs HandlerState) error {

	comp := app.LoginForm()

	// If authenticated and authorized, allow access to default homepage.
	if err := hs.Authenticated(); err == nil {

		// Get list of events.
		events := []Event{}
		if err := runQuery(hs, ReadEvents, &events); err != nil {
			return err
		}

		// Set as component to send.
		comp = app.EventsHome(events)

	} else {
		log.Println(err)
	}

	// Respond to client with fully rendered home.
	if err := hs.Local.RespondHtml(app.Home(comp)); err != nil {
		return err
	}

	return nil

})

// Create login session based on provided form credentials.
var CreateLoginErr = addHandlerFunc(utils.ApiPath("login"), "post", func(hs HandlerState) error {

	// Retrieve email from user.
	email, err := hs.FormGet("Email")
	if err != nil {
		return utils.ErrPrep(err, "unable to get Email")
	}

	// Retrieve plaintext password from user.
	passwordPlaintext, err := hs.FormGet("PasswordPlaintext")
	if err != nil {
		return utils.ErrPrep(err, "unable to get PasswordPlaintext")
	}

	// Get list of existing users that have matching email
	readUsersWithEmail := ReadUsers.WHERE(t.Baseuser.Email.EQ(postgres.String(email)))
	usersWithEmail := []User{}
	runQuery(hs, readUsersWithEmail, &usersWithEmail)
	if err != nil {
		return err
	}

	users := []User{}
	if err := runQuery(hs, ReadUsers, &users); err != nil {
		return err
	}

	// Error if could not find user with email
	if len(usersWithEmail) == 0 {
		hs.Local.RespondHtml(app.StatusMessage(app.Failure(), "unable to find user with email/password combination"), http.StatusInternalServerError)
		return errors.New("could not find user with matching email")
	}

	// Get first match
	user := &usersWithEmail[0]

	// Check if provided password matches user.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordPlaintext)); err != nil {
		hs.Local.RespondHtml(app.StatusMessage("danger", "unable to find user with email/password combination"), http.StatusInternalServerError)
		return utils.ErrPrep(err, "password does not match user with provided email")
	}

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
		ID:        uuid.NewString(),
	}

	// Cache claims.
	authenticatedUsers.Add(claims)

	// Create signed token string.
	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(tokenSecret)
	if err != nil {
		hs.Local.RespondHtml(app.StatusMessage("danger", "unable to sign JWT token"), http.StatusInternalServerError)
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

	// Store cookie into Set-Cookie header for future usage.
	http.SetCookie(hs.Local.ResponseWriter, &authCookie)

	// Get list of events.
	events := []Event{}
	if err := runQuery(hs, ReadEvents, &events); err != nil {
		return err
	}

	// Respond to request with list of events
	hs.Local.RespondHtml(app.EventsHome(events))

	return nil

})

// Get login form used to create a login session.
var ReadLoginErr = addHandlerFunc(utils.ApiPath("login"), "get", func(hs HandlerState) error {
	hs.Local.RespondHtml(app.LoginForm())
	return nil
})

// Get signup form used to create a student account.
var ReadSignupErr = addHandlerFunc(utils.ApiPath("signup"), "get", func(hs HandlerState) error {

	universities := []m.University{}
	if err := runQuery(hs, ReadUniversities, &universities); err != nil {
		return err
	}

	// Respond with HTML.
	hs.Local.RespondHtml(app.SignupForm(universities))
	return nil
})

var createStudentParamKeys = []string{
	"Email",
	"PasswordRaw",
	"NameFirst",
	"NameLast",
}

// Create new student that's associated with a university.
var CreateStudentErr = addHandlerFunc(utils.ApiPath("signup"), "post", func(hs HandlerState) error {

	// Hash user-provided password input.
	if err := hs.HashPasswordInput(); err != nil {
		hs.Local.RespondHtml(app.StatusMessage("danger", err.Error()), http.StatusInternalServerError)
		return err
	}

	// Copy Form data into base user parameters
	baseUserParams := m.Baseuser{}
	if err := hs.ToParams(&baseUserParams); err != nil {
		return err
	}

	// Create BaseUser
	createBaseUser := CreateBaseUser.MODEL(baseUserParams).RETURNING(t.Baseuser.ID)
	newBaseUsers := []m.Baseuser{}
	if err := runQuery(hs, createBaseUser, &newBaseUsers); err != nil {
		return err
	}

	// Prepare SQL statement for promoting the BaseUser into a Student.
	createStudent := CreateStudent.VALUES(newBaseUsers[0].ID)

	// Create Student
	if err := runQuery(hs, createStudent, nil); err != nil {
		hs.Local.RespondHtml(app.StatusMessage("danger", err.Error()), http.StatusInternalServerError)
		return err
	}

	return nil

})

// Get list of users.
var ReadUsersErr = addHandlerFunc(utils.ApiPath("users"), "get", func(hs HandlerState) error {

	students := []User{}
	if err := runQuery(hs, ReadUsers, students); err != nil {
		return err
	}

	return hs.Local.RespondHtml(app.CreatedBaseUsers(students))

})

// Get list of universities.
var ReadUniversitiesErr = addHandlerFunc(utils.ApiPath("university"), "get", func(hs HandlerState) error {

	// Read list of created universities.
	universities := []m.University{}
	if err := runQuery(hs, ReadUniversities, &universities); err != nil {
		return err
	}

	// Respond with universities.
	if err := hs.Local.RespondHtml(app.CreatedUniversities(universities)); err != nil {
		return err
	}

	return nil

})

var createUniversityParamKeys = []string{
	"Title",
	"About",
	"Latitude",
	"Longitude",
}

// Create a new university record.
var CreateUniversityErr = addHandlerFunc(utils.ApiPath("university"), "post", func(hs HandlerState) error {

	// Create new university.
	universities := []m.University{}
	if err := runQuery(hs, ReadUniversities, &universities); err != nil {
		hs.Local.RespondHtml(app.StatusMessage("danger", "Unable to create university"), http.StatusInternalServerError)
		return err
	}

	// Respond with status.
	if err := hs.Local.RespondHtml(app.StatusMessage("success", "Created new university!")); err != nil {
		return err
	}

	return nil

})

type eventParams struct {
	Event
	EventType string
}

var CreateEventErr = addHandlerFunc(utils.ApiPath("event"), "post", func(hs HandlerState) error {

	createUniversityParams := m.University{}
	if err := hs.ToParams(&createUniversityParams); err != nil {
		return err
	}

	if err := runQuery(hs, CreateUniversity.MODEL(createUniversityParams), nil); err != nil {
		return err
	}

	return nil
})

var universityForms = []map[string][]string{
	{
		"Title":     {"University of Central Florida"},
		"About":     {"A public research university with its main campus in unincorporated Orange County, Florida"},
		"Latitude":  {"28.602540027323045"},
		"Longitude": {"-81.20002623717798"},
	},
	{
		"Title":     {"Massachusetts Institute of Technology"},
		"About":     {"A private land-grant research university in Cambridge, Massachusetts"},
		"Latitude":  {"42.360134050711146"},
		"Longitude": {"-71.09410939970417"},
	},
}

// Helper route to populate database with default values.
var InitDatabaseErr = addHandlerFunc(utils.ApiPath("init"), "post", func(hs HandlerState) error {

	// Get list of existing universities
	universities := []m.University{}
	if err := runQuery(hs, ReadUniversities, &universities); err != nil {
		return err
	}

	// Check if already created universities
	if len(universities) != 0 {
		return errors.New("database already has universities, skipping")
	}

	for _, form := range universityForms {

		// Copy into HandlerState's form data the university Form
		hs.Local.Request.Form = url.Values(form)

		// Copy Form data into params struct for University table
		params := m.University{}
		if err := hs.ToParams(&params); err != nil {
			return err
		}

		// Insert university
		if err := runQuery(hs, CreateUniversity.MODEL(params), nil); err != nil {
			return err
		}

	}

	hs.Local.RespondHtml(app.StatusMessage(app.Success(), "initialized database with default values"))

	return nil
})
