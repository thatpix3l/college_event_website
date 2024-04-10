package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/microcosm-cc/bluemonday"
	. "github.com/thatpix3l/cew/src/gen_sql"
	"github.com/thatpix3l/cew/src/utils"
	"golang.org/x/crypto/bcrypt"

	app "github.com/thatpix3l/cew/src/gen_templ"

	m "github.com/thatpix3l/cew/src/gen_sql/college_event_website/cew/model"
	t "github.com/thatpix3l/cew/src/gen_sql/college_event_website/cew/table"
)

// Get homepage.
var ReadHomepageErr = addHandlerFunc("/", "get", func(hs HandlerState) error {

	comp := app.LoginForm()

	// If authenticated and authorized, allow access to default homepage.
	if err := hs.Authenticated(); err == nil {

		// Get events home UI
		eventListComp, err := eventListHome(hs)
		if err != nil {
			return err
		}

		// Set as component to send.
		comp = eventListComp

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
	query := ReadUsers().WHERE(t.Baseuser.Email.EQ(pg.String(email)))
	users := []User{}
	runQuery(hs, query, &users)
	if err != nil {
		return err
	}

	// Error if could not find user with email
	if len(users) == 0 {
		hs.Local.RespondHtml(app.StatusMessage(app.Failure, "unable to find user with email/password combination"), http.StatusInternalServerError)
		return errors.New("could not find user with matching email")
	}

	// Get first match
	user := users[0]

	// Check if provided password matches user.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordPlaintext)); err != nil {
		hs.Local.RespondHtml(app.StatusMessage(app.Failure, "unable to find user with email/password combination"), http.StatusInternalServerError)
		return utils.ErrPrep(err, "password does not match user with provided email")
	}

	// Authenticate
	hs.Authenticate(user)

	// Get list of events; customized based on user authorization
	comp, err := eventListHome(hs)
	if err != nil {
		return err
	}

	// Respond to user list of events
	if err := hs.Local.RespondHtml(comp); err != nil {
		return err
	}

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
	if err := runQuery(hs, ReadUniversities(), &universities); err != nil {
		return err
	}

	// Respond with HTML.
	hs.Local.RespondHtml(app.SignupForm(universities))
	return nil
})

// Build query that returns list of events viewable by user in current handle
func eventListQuery(hs HandlerState) (pg.SelectStatement, error) {

	// Get user from HandlerState
	user := User{}
	if err := hs.GetUser(&user); err != nil {
		return nil, err
	}

	// Make bool expression that checks if an event and student both are from the same university
	sameUniversity := pg.Bool(false)
	if user.Student != nil && user.Student.UniversityID != nil {
		sameUniversity = t.Baseevent.UniversityID.EQ(pg.String(*user.Student.UniversityID))
	}

	// If event is public and approved, allow
	public := t.Baseevent.ID.EQ(t.Publicevent.ID).AND(t.Publicevent.Approved)

	// If event is private and user is part of the same university, allow
	private := t.Baseevent.ID.EQ(t.Privateevent.ID).AND(sameUniversity)

	// If event is rso and user is part of the same university and part of the rso related to event, allow
	sameRso := pg.Bool(false)
	if user.Rsomember != nil {
		sameRso = t.Rsoevent.RsoID.EQ(pg.String(user.Rsomember.RsoID))
	}

	rso := t.Baseevent.ID.EQ(t.Rsoevent.ID).AND(sameRso).AND(sameUniversity)

	query := ReadEvents().WHERE(public.OR(private).OR(rso))

	return query, nil
}

// Get list of events
func eventList(hs HandlerState) ([]Event, error) {

	query, err := eventListQuery(hs)
	if err != nil {
		return nil, err
	}

	events := []Event{}
	if err := runQuery(hs, query, &events); err != nil {
		return nil, err
	}

	return events, nil

}

// Create UI for the homepage of a list of events
func eventListHome(hs HandlerState) (templ.Component, error) {

	events, err := eventList(hs)
	if err != nil {
		return nil, err
	}

	// Create UI that uses current state of list of events
	return app.StackComponents(
		app.NavBar("events"),
		app.EventSearch(),
		app.Interactive(app.EventList(events)),
	), nil
}

// Create new student that's associated with a university.
var CreateStudentErr = addHandlerFunc(utils.ApiPath("signup"), "post", func(hs HandlerState) error {

	// Hash user-provided password input.
	if err := hs.HashPasswordInput(); err != nil {
		hs.Local.RespondHtml(app.StatusMessage(app.Failure, err.Error()), http.StatusInternalServerError)
		return err
	}

	// Copy Form data into base user parameters
	baseUserParams := m.Baseuser{}
	if err := hs.ToParams(&baseUserParams); err != nil {
		return err
	}

	// User we're going to create
	user := User{}

	// Create BaseUser
	if err := runQuery(hs,
		CreateBaseUser().MODEL(baseUserParams).RETURNING(t.Baseuser.AllColumns),
		&user.Baseuser); err != nil {
		return err
	}

	// Create Student based off of base userr
	user.Student = &m.Student{}
	if err := runQuery(hs, CreateStudent().VALUES(user.Baseuser.ID).RETURNING(t.Student.AllColumns), user.Student); err != nil {
		hs.Local.RespondHtml(app.StatusMessage(app.Failure, err.Error()), http.StatusInternalServerError)
		return err
	}

	// Authenticate and cache user
	hs.Authenticate(user)

	// Get list events viewable by user
	comp, err := eventListHome(hs)
	if err != nil {
		return err
	}

	// Respond to user list of events
	if err := hs.Local.RespondHtml(comp); err != nil {
		return err
	}

	return nil

})

// Get list of users.
var ReadUsersErr = addHandlerFunc(utils.ApiPath("users"), "get", func(hs HandlerState) error {

	students := []User{}
	if err := runQuery(hs, ReadUsers(), students); err != nil {
		return err
	}

	return hs.Local.RespondHtml(app.CreatedBaseUsers(students))

})

// Get list of universities.
var ReadUniversitiesErr = addHandlerFunc(utils.ApiPath("university"), "get", func(hs HandlerState) error {

	// Read list of created universities.
	universities := []m.University{}
	if err := runQuery(hs, ReadUniversities(), &universities); err != nil {
		return err
	}

	// Respond with universities.
	if err := hs.Local.RespondHtml(app.CreatedUniversities(universities)); err != nil {
		return err
	}

	return nil

})

// Create a new university record.
var CreateUniversityErr = addHandlerFunc(utils.ApiPath("university"), "post", func(hs HandlerState) error {

	// Create new university.
	universities := []m.University{}
	if err := runQuery(hs, ReadUniversities(), &universities); err != nil {
		hs.Local.RespondHtml(app.StatusMessage(app.Failure, "Unable to create university"), http.StatusInternalServerError)
		return err
	}

	// Respond with status.
	if err := hs.Local.RespondHtml(app.StatusMessage(app.Success, "Created new university!")); err != nil {
		return err
	}

	return nil

})

var CreateEventErr = addHandlerFunc(utils.ApiPath("event"), "post", func(hs HandlerState) error {

	createUniversityParams := m.University{}
	if err := hs.ToParams(&createUniversityParams); err != nil {
		return err
	}

	if err := runQuery(hs, CreateUniversity().MODEL(createUniversityParams), nil); err != nil {
		return err
	}

	return nil
})

var ReadRsosErr = addHandlerFunc(utils.ApiPath("home/rsos"), "get", func(hs HandlerState) error {

	rsos := []m.Rso{}
	if err := runQuery(hs, ReadRsos(), &rsos); err != nil {
		return err
	}

	if err := hs.Local.RespondHtml(app.StackComponents(
		app.NavBar("rsos"),
		app.CreatedRsos(rsos),
	)); err != nil {
		return err
	}

	return nil
})

var ReadEventsErr = addHandlerFunc(utils.ApiPath("event"), "get", func(hs HandlerState) error {

	// Copy of query for getting events
	query, err := eventListQuery(hs)
	if err != nil {
		return err
	}

	// Get url queries
	urlQueries := hs.Local.Request.URL.Query()

	// Get search term for filtering; if empty, simply return list of ALL events
	searchTerm := urlQueries.Get("search")
	if searchTerm == "" {

		events := []Event{}
		if err := runQuery(hs, query, &events); err != nil {
			return err
		}

		hs.Local.RespondHtml(app.EventList(events))
		return nil
	}

	// Modify search term so it can match as a substring of any stored strings
	searchTerm = "%" + searchTerm + "%"

	// Part of event title matches search
	title := t.Baseevent.Title.LIKE(pg.String(searchTerm))

	// Part of event body matches search
	body := t.Baseevent.About.LIKE(pg.String(searchTerm))

	// At least one of the event's tags matches search
	tags := t.Baseevent.ID.EQ(t.Taggedevent.BaseEventID).AND(
		t.Taggedevent.TagID.EQ(t.Tag.ID),
	).AND(
		t.Tag.Title.EQ(pg.String(searchTerm)),
	)

	// Modify query to only return events that have either the title, body, or tags matching search query
	query = query.WHERE(title.OR(body).OR(tags))

	// Get list of events, based off of criteria
	events := []Event{}
	if err := runQuery(hs, query, &events); err != nil {
		return err
	}

	// Respond to client
	hs.Local.RespondHtml(app.EventList(events))

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

type ucfEvent struct {
	Title        string
	Description  string
	Starts       string
	ContactPhone string
	ContactEmail string
	Tags         []string
}

var p = bluemonday.UGCPolicy()

func unescape(input string) string {
	output := p.Sanitize(html.UnescapeString(input))
	return output
}

// Helper route to populate database with default values.
var InitDatabaseErr = addHandlerFunc(utils.ApiPath("init"), "post", func(hs HandlerState) error {

	universitiesParams := []m.University{}
	for _, form := range universityForms {

		// Copy current university Form into HandlerState's Form
		hs.Local.Request.Form = url.Values(form)

		// Marhsal current Form into params struct
		params := m.University{}
		if err := hs.ToParams(&params); err != nil {
			return err
		}

		// Store params for later insertion
		universitiesParams = append(universitiesParams, params)

	}

	// Insert and return copy of universities
	fmt.Print("inserting default universities...")
	universities := []m.University{}
	if err := runQuery(hs,
		CreateUniversity().MODELS(universitiesParams).RETURNING(t.University.AllColumns),
		&universities); err != nil {
		return err
	}
	fmt.Println("done")

	// ID of the UCF university
	ucfId := func() string {
		for _, university := range universities {
			if university.Title == "University of Central Florida" {
				return university.ID
			}
		}
		return ""
	}()

	// Get HTTP Response from UCF events feed
	resp, err := http.Get("https://events.ucf.edu/feed.json")
	if err != nil {
		return err
	}

	// Read body of response into buffer
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unmarshal buffer into structured data
	ucfEvents := []ucfEvent{}
	if err := json.Unmarshal(buf, &ucfEvents); err != nil {
		return err
	}

	// Convert UCF-created event into version the database can work with
	eventsParams := []m.Baseevent{}
	for _, ue := range ucfEvents {
		newParams := m.Baseevent{
			Title:        unescape(ue.Title),
			About:        unescape(ue.Description),
			ContactPhone: unescape(ue.ContactPhone),
			ContactEmail: unescape(ue.ContactEmail),
			UniversityID: ucfId,
		}

		eventsParams = append(eventsParams, newParams)

	}

	events := []m.Baseevent{} // inserted events

	// Store database-compatible events
	fmt.Print("inserting events...")
	if err := runQuery(hs, CreateEvent().MODELS(eventsParams).RETURNING(t.Baseevent.AllColumns), &events); err != nil {
		return err
	}
	fmt.Println("done")

	// Make all events pulled from UCF public
	publicEvents := []m.Publicevent{}
	for _, event := range events {
		publicEvents = append(publicEvents, m.Publicevent{
			ID:       event.ID,
			Approved: true,
		})
	}

	// Insert public events
	if err := runQuery(hs, CreatePublicEvent().MODELS(publicEvents), nil); err != nil {
		return err
	}

	taggedEventsParams := []m.Taggedevent{} // list of event-tag tuples

	fmt.Print("inserting tags...")
	for i := range ucfEvents {
		eventUcf := ucfEvents[i] // current UCF-created event
		tagsParams := []m.Tag{}  // database-compatible list of tags for current event

		// Make current event's tags database compatible
		for _, tagName := range eventUcf.Tags {
			tagsParams = append(tagsParams, m.Tag{Title: tagName})
		}

		// Store current event's tags into database
		tags := []m.Tag{}
		if err := runQuery(hs, CreateTag().MODELS(tagsParams).RETURNING(t.Tag.AllColumns), &tags); err != nil {
			return err
		}

		eventDb := events[i] // current database-compatible event

		// For each tag, make tuple with it and current event; add to list of tagged events
		for _, tag := range tags {
			taggedEventsParams = append(taggedEventsParams, m.Taggedevent{TagID: tag.ID, BaseEventID: eventDb.ID})
		}

	}
	fmt.Println("done")

	// Store list of event-tag tuples into database
	fmt.Print("storing all event-tag tuples...")
	if err := runQuery(hs, CreateTaggedEvent().MODELS(taggedEventsParams), nil); err != nil {
		return err
	}
	fmt.Println("done")

	hs.Local.RespondHtml(app.StatusMessage(app.Success, "Initialized database with default values"))

	return nil
})

// Respond to client the UI for viewing an event
func eventInfo(hs HandlerState) error {

	// Get event ID
	eventId := chi.URLParam(hs.Local.Request, "event_id")
	if eventId == "" {
		return errors.New("get event info: did not provide event ID")
	}

	// Build query for getting list of events user can view
	query, err := eventListQuery(hs)
	if err != nil {
		return err
	}

	query = query.WHERE(t.Baseevent.ID.EQ(pg.String(eventId))).ORDER_BY(t.Comment.PostTimestamp.ASC())

	// Run query and story output
	events := []Event{}
	if err := runQuery(hs, query, &events); err != nil {
		return err
	}

	// Error if event with ID does not exist
	if len(events) == 0 {
		return err
	}

	// Get first match; DB constraints ensure there is at most 1
	event := events[0]

	// Get user that's accessing page based upon their JWT credentials
	user := User{}
	if err := hs.GetUser(&user); err != nil {
		return err
	}

	// Respond to user with UI for viewing event
	hs.Local.RespondHtml(app.StackComponents(
		app.Event(event),
		app.CreateComment(event),
		app.CommentList(user, event.Comments),
	))

	return nil

}

var ReadEventInfo = addHandlerFunc(utils.ApiPath("event/{event_id}"), "get", eventInfo)

var ReadEventsCommentCreated = addHandlerFunc(utils.ApiPath("event/{event_id}/comment"), "post", func(hs HandlerState) error {

	// User that initiated request
	user := User{}
	if err := hs.GetUser(&user); err != nil {
		return err
	}

	// Body of comment to insert
	commentBody, err := hs.FormGet("CommentBody")
	if err != nil {
		return err
	}

	// Event ID the comment is supposed to be for
	eventId := chi.URLParam(hs.Local.Request, "event_id")
	if eventId == "" {
		return errors.New("create comment for event: did not provide event ID")
	}

	// Build comment params
	comment := m.Comment{
		Body:          commentBody,
		StudentID:     &user.Baseuser.ID,
		BaseEventID:   eventId,
		PostTimestamp: time.Now(),
	}

	// Store comment
	if err := runQuery(hs, CreateComment().MODEL(comment), nil); err != nil {
		return err
	}

	// Respond to user with specific event they just commented on
	if err := eventInfo(hs); err != nil {
		return err
	}

	return nil
})

var ReadEventsCommentRemoved = addHandlerFunc(utils.ApiPath("event/{event_id}/comment/{comment_id}"), "delete", func(hs HandlerState) error {

	// Get comment ID
	commentId := chi.URLParam(hs.Local.Request, "comment_id")
	if commentId == "" {
		return errors.New("remove comment from event: did not provide comment ID")
	}

	// Get user
	user := User{}
	if err := hs.GetUser(&user); err != nil {
		return err
	}

	// Error if user who initiated request is not a student
	if user.Student == nil {
		return errors.New("remove comment from event: user is not a student")
	}

	// Query to remove chosen comment that was posted by student initiating request
	query := t.Comment.DELETE().WHERE(t.Comment.ID.EQ(pg.String(commentId)).AND(t.Comment.StudentID.EQ(pg.String(user.Student.ID))))

	if err := runQuery(hs, query, nil); err != nil {
		return err
	}

	if err := eventInfo(hs); err != nil {
		return err
	}

	return nil
})

var ReadEventListHome = addHandlerFunc(utils.ApiPath("home/events"), "get", func(hs HandlerState) error {

	comp, err := eventListHome(hs)
	if err != nil {
		return err
	}

	hs.Local.RespondHtml(comp)

	return nil
})
