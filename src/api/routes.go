package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
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

	isRsoSameSource := pg.Bool(false)
	if user.Rsomember != nil {
		isRsoSameSource = IsRsoSameSource(user.Rsomember.RsoID)
	}

	query := ReadEvents().WHERE(
		IsPublicApproved().OR(
			IsPrivateSameUniversity(user.UniversityID),
		).OR(
			isRsoSameSource,
		),
	)

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
	return app.EventListHome(events), nil
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
	baseUser := m.Baseuser{}

	// Create BaseUser
	if err := runQuery(hs,
		CreateBaseUser().MODEL(baseUserParams).RETURNING(t.Baseuser.AllColumns),
		&baseUser); err != nil {
		return err
	}

	// Pull university ID user wants to be a part of
	universityId, err := hs.FormGet("UniversityID")
	if err != nil {
		return err
	}

	// Create student based off of base user
	student := m.Student{
		ID:           baseUser.ID,
		UniversityID: universityId,
	}

	// Insert student into database
	if err := runQuery(hs, CreateStudent().MODEL(student), nil); err != nil {
		hs.Local.RespondHtml(app.StatusMessage(app.Failure, err.Error()), http.StatusInternalServerError)
		return err
	}

	// Retrieve fully configured student from database
	user := User{}
	readNewStudentQuery := ReadUsers().WHERE(t.Baseuser.ID.EQ(pg.String(baseUser.ID)))

	if err := runQuery(hs, readNewStudentQuery, &user); err != nil {
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

var ReadEventCreatorErr = addHandlerFunc(utils.ApiPath("event/creator"), "get", func(hs HandlerState) error {

	// Get universities
	universities := []m.University{}
	if err := runQuery(hs, ReadUniversities(), &universities); err != nil {
		return err
	}

	// Get RSOs
	rsos := []Rso{}
	if err := runQuery(hs, ReadRsosMinMembers(), &rsos); err != nil {
		return err
	}

	// Respond to client with UI for creating a new event
	hs.Local.RespondHtml(app.EventCreator(universities, rsos))

	return nil
})

func nothing(hs HandlerState) error {
	return nil
}

var ReadEventCreatorPublicOptionsErr = addHandlerFunc(utils.ApiPath("event/creator/public"), "get", nothing)

var ReadEventCreatorPrivateOptionsErr = addHandlerFunc(utils.ApiPath("event/creator/private"), "get", nothing)

var ReadEventCreatorRsoOptionsErr = addHandlerFunc(utils.ApiPath("event/creator/rso"), "get", func(hs HandlerState) error {

	// Get RSOs with minimum amount of members required
	rsos := []Rso{}
	if err := runQuery(hs, ReadRsosMinMembers(), &rsos); err != nil {
		return err
	}

	hs.Local.RespondHtml(app.EventCreatorRsoOptions(rsos))

	return nil

})

func fixTime(hs HandlerState, key string) error {

	// Get time
	startTime, err := hs.FormGet("StartTime")
	if err != nil {
		return err
	}

	// Append thing so Go can parse/marshal/whatever in background
	hs.Local.Request.Form[key] = []string{startTime + ":00Z"}

	log.Println(startTime)

	return nil

}

var CreateEventErr = addHandlerFunc(utils.ApiPath("event"), "post", func(hs HandlerState) error {

	// Fix start time
	if err := fixTime(hs, "StartTime"); err != nil {
		return err
	}

	// Get event type
	eventType, err := hs.FormGet("EventType")
	if err != nil {
		return err
	}

	// Decode base event parameters
	event := Event{}
	if err := hs.ToParams(&event.Baseevent, "PostTime"); err != nil {
		return err
	}

	// Store base event
	if err := runQuery(hs, CreateEvent().MODEL(event.Baseevent), nil); err != nil {
		return err
	}

	// Add appropriate event type to database
	switch eventType {

	// Make event private and viewable to only those in the same school
	case "private":
		event.Privateevent = &m.Privateevent{ID: event.Baseevent.ID}
		query := t.Privateevent.INSERT(t.Privateevent.ID).MODEL(event.Privateevent)
		if err := runQuery(hs, query, nil); err != nil {
			return err
		}

		// Make event public and viewable to everyone
	case "public":
		event.Publicevent = &m.Publicevent{ID: event.Baseevent.ID}
		query := t.Publicevent.INSERT(t.Publicevent.ID).MODEL(event.Publicevent)
		if err := runQuery(hs, query, nil); err != nil {
			return err
		}

		// Make event Rso-specific and viewable only to its associated Rso members
	case "rso":
		rsoId, err := hs.FormGet("RsoId")
		if err != nil {
			return err
		}

		event.Rsoevent = &m.Rsoevent{ID: event.Baseevent.ID, RsoID: rsoId}
		query := t.Rsoevent.INSERT(t.Rsoevent.ID, t.Rsoevent.RsoID).MODEL(event.Rsoevent)
		if err := runQuery(hs, query, nil); err != nil {
			return err
		}
	}

	// Retrieved updated list of events after inserting event
	events, err := eventList(hs)
	if err != nil {
		return err
	}

	hs.Local.RespondHtml(app.EventListInteractive(events))

	return nil
})

var ReadRsosHomeErr = addHandlerFunc(utils.ApiPath("home/rsos"), "get", func(hs HandlerState) error {

	rsos := []Rso{}
	if err := runQuery(hs, ReadRsos(), &rsos); err != nil {
		return err
	}

	if err := hs.Local.RespondHtml(app.RsoListHome(rsos)); err != nil {
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
	if len(searchTerm) == 0 {

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
		"Title": {"University of Central Florida"},
		"About": {"A public research university with its main campus in unincorporated Orange County, Florida"},
	},
	{
		"Title": {"Massachusetts Institute of Technology"},
		"About": {"A private land-grant research university in Cambridge, Massachusetts"},
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

	// Make all events pulled from UCF public and approved
	publicEvents := []m.Publicevent{}
	for _, event := range events {
		publicEvents = append(publicEvents, m.Publicevent{
			ID:       event.ID,
			Approved: true,
		})
	}

	// Insert public events
	query := t.Publicevent.INSERT(t.Publicevent.ID, t.Publicevent.Approved).MODELS(publicEvents)
	if err := runQuery(hs, query, nil); err != nil {
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
	fmt.Print("inserting tags for each event...")
	if err := runQuery(hs, CreateTaggedEvent().MODELS(taggedEventsParams), nil); err != nil {
		return err
	}
	fmt.Println("done")

	hs.Local.RespondHtml(app.StatusMessage(app.Success, "Initialized database with default values"))

	return nil
})

// Respond to client the UI for viewing an event
func eventInfo(hs HandlerState, eventId string) error {

	// Get event ID
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
		app.CommentCreator(event),
		app.CommentList(user, event.Comments),
	))

	return nil

}

var ReadEventInfo = addHandlerFunc(utils.ApiPath("event/{event_id}"), "get", func(hs HandlerState) error {

	eventId := chi.URLParam(hs.Local.Request, "event_id")
	if err := eventInfo(hs, eventId); err != nil {
		return err
	}

	return nil
})

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
	if err := eventInfo(hs, eventId); err != nil {
		return err
	}

	return nil
})

var DeleteComment = addHandlerFunc(utils.ApiPath("comment/{comment_id}"), "delete", func(hs HandlerState) error {

	// Get user
	user := User{}
	if err := hs.GetUser(&user); err != nil {
		return err
	}

	// Error if user who initiated request is not a student
	if user.StudentFull == nil {
		return errors.New("remove comment from event: user is not a student")
	}

	// Get comment ID
	commentId := chi.URLParam(hs.Local.Request, "comment_id")
	if commentId == "" {
		return errors.New("delete comment: did not provide comment ID")
	}

	// Cache events associated with comment ID; obviously, should only be one
	events := []Event{}
	if err := runQuery(hs, ReadEvents().WHERE(t.Baseevent.ID.EQ(t.Comment.BaseEventID)), &events); err != nil {
		return err
	}

	if len(events) == 0 {
		return errors.New("delete comment: could not find event comment is associated with")
	}

	eventId := events[0].Baseevent.ID // Specifically, cache the event's ID

	// Query to remove chosen comment that was posted by student initiating request
	query := t.Comment.DELETE().WHERE(t.Comment.ID.EQ(pg.String(commentId)).AND(t.Comment.StudentID.EQ(pg.String(user.Student.ID))))

	// Attempt to remove comment
	if err := runQuery(hs, query, nil); err != nil {
		return err
	}

	// Respond to client with UI for viewing event that the comment was deleted from
	if err := eventInfo(hs, eventId); err != nil {
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

var ReadCommentErr = addHandlerFunc(utils.ApiPath("comment/{comment_id}"), "get", func(hs HandlerState) error {

	// Get comment ID
	commentId := chi.URLParam(hs.Local.Request, "comment_id")
	if commentId == "" {
		return errors.New("get comment: comment ID is empty")
	}

	// Build query to select all comments with matching ID; obviously, at most one
	query := t.Comment.SELECT(t.Comment.AllColumns).WHERE(t.Comment.ID.EQ(pg.String(commentId)))

	// Run query and store comments
	comments := []m.Comment{}
	if err := runQuery(hs, query, &comments); err != nil {
		return err
	}

	// Error if could not find comment
	if len(comments) == 0 {
		return errors.New("get comment: comment with ID does not exist")
	}

	// Get user
	user := User{}
	if err := hs.GetUser(&user); err != nil {
		return err
	}

	// Respond to given user the comment
	hs.Local.RespondHtml(app.Comment(user, comments[0]))

	return nil

})

var ReadUpdateCommentErr = addHandlerFunc(utils.ApiPath("comment/{comment_id}/update"), "get", func(hs HandlerState) error {

	commentId := chi.URLParam(hs.Local.Request, "comment_id")
	if commentId == "" {
		return errors.New("get comment updater: did not provide comment ID")
	}

	hs.Local.RespondHtml(app.CommentUpdater(commentId))

	return nil

})

var UpdateComment = addHandlerFunc(utils.ApiPath("comment/{comment_id}"), "patch", func(hs HandlerState) error {

	commentId := chi.URLParam(hs.Local.Request, "comment_id")
	if commentId == "" {
		return errors.New("update comment: did not provide comment ID")
	}

	user := User{}
	if err := hs.GetUser(&user); err != nil {
		return err
	}

	if user.StudentFull == nil {
		return errors.New("update comment: user is not a student")
	}

	commentBody, err := hs.FormGet("CommentBody")
	if err != nil {
		return err
	}

	query := t.Comment.UPDATE(t.Comment.Body).SET(commentBody).WHERE(
		t.Comment.ID.EQ(pg.String(commentId)).AND(
			t.Comment.StudentID.EQ(pg.String(user.Student.ID)),
		),
	).RETURNING(t.Comment.AllColumns)

	comments := []m.Comment{}
	if err := runQuery(hs, query, &comments); err != nil {
		return err
	}

	if len(comments) == 0 {
		return errors.New("update comment: comment with ID does not exist")
	}

	hs.Local.RespondHtml(app.Comment(user, comments[0]))

	return nil
})

var ReadRsoCreate = addHandlerFunc(utils.ApiPath("rso/create"), "get", func(hs HandlerState) error {

	universities := []m.University{}
	if err := runQuery(hs, ReadUniversities(), &universities); err != nil {
		return err
	}

	hs.Local.RespondHtml(app.CreateRso(universities))

	return nil
})

type RsoParams struct {
	Title        string
	About        string
	UniversityID string
}

var RsoCreate = addHandlerFunc(utils.ApiPath("rso"), "post", func(hs HandlerState) error {

	params := RsoParams{}
	if err := hs.ToParams(&params); err != nil {
		return err
	}

	query := t.Rso.INSERT(t.Rso.MutableColumns).MODEL(params)

	if err := runQuery(hs, query, nil); err != nil {
		return err
	}

	rsos := []Rso{}
	if err := runQuery(hs, ReadRsos(), &rsos); err != nil {
		return err
	}

	hs.Local.RespondHtml(app.RsoListHome(rsos))

	return nil
})

var ReadRsoList = addHandlerFunc(utils.ApiPath("rso/{rso_id}"), "get", func(hs HandlerState) error {

	user := User{}
	if err := hs.GetUser(&user); err != nil {
		return err
	}

	rsoId := chi.URLParam(hs.Local.Request, "rso_id")
	if rsoId == "" {
		return errors.New("get rso: did not provide rso ID")
	}

	rsos := []Rso{}
	if err := runQuery(hs, ReadRsos().WHERE(t.Rso.ID.EQ(pg.String(rsoId))), &rsos); err != nil {
		return err
	}

	if len(rsos) == 0 {
		return errors.New("get rso: rso with provided ID does not exist")
	}

	hs.Local.RespondHtml(app.RsoInfo(rsos[0], user))

	return nil
})

var CreateRsoMemberErr = addHandlerFunc(utils.ApiPath("rso/{rso_id}/member"), "post", func(hs HandlerState) error {

	// Get user
	user := User{}
	if err := hs.GetUser(&user); err != nil {
		return err
	}

	// Error if user is not a student
	if user.StudentFull == nil {
		return errors.New("create rso member: user is not a student")
	}

	// Error if did not provide ID for an rso
	rsoId := chi.URLParam(hs.Local.Request, "rso_id")
	if rsoId == "" {
		return errors.New("create rso member: did not provide rso ID")
	}

	rsoMember := m.Rsomember{
		ID:    user.Student.ID,
		RsoID: rsoId,
	}

	// Create RSO member, based off of user initiating request
	createRsoMemberQuery := CreateRsoMember().MODEL(rsoMember)
	if err := runQuery(hs, createRsoMemberQuery, nil); err != nil {
		return err
	}

	// Retrievespecific RSO for whom we just created a member for
	rso := []Rso{}
	readRsosQuery := ReadRsos().WHERE(t.Rso.ID.EQ(pg.String(rsoId)))
	if err := runQuery(hs, readRsosQuery, &rso); err != nil {
		return err
	}

	// Error if could not retrieve specific Rso
	if len(rso) == 0 {
		return errors.New("create rso member: cannot retrieve specified rso for viewing")
	}

	// Get user from database after addition of their Rso membership
	userAfterDeletion := User{}
	if err := hs.GetUser(&userAfterDeletion); err != nil {
		return err
	}

	// Respond to client with UI for viewing the info of the RSO
	hs.Local.RespondHtml(app.RsoInfo(rso[0], userAfterDeletion))

	return nil
})

var DeleteRsoMemberErr = addHandlerFunc(utils.ApiPath("rso/{rso_id}/member"), "delete", func(hs HandlerState) error {

	// Get Rso id
	rsoId := chi.URLParam(hs.Local.Request, "rso_id")
	if rsoId == "" {
		return errors.New("delete rso member: did not provide rso ID")
	}

	// Get user initiating request
	user := User{}
	if err := hs.GetUser(&user); err != nil {
		return err
	}

	// Error if request initiator is not a member of any Rso
	if user.Rsomember == nil {
		return errors.New("delete rso member: user initiating request is not a member of any RSO")
	}

	// Error if request initiator is not a member of Rso that has provided ID
	if user.Rsomember.RsoID != rsoId {
		return errors.New("delete rso member: user initiating request is not a member of provided RSO")
	}

	deleteRsoMemberQuery := t.Rsomember.DELETE().WHERE(
		t.Rsomember.RsoID.EQ(pg.String(user.Rsomember.RsoID)).AND(
			t.Rsomember.ID.EQ(pg.String(user.Rsomember.ID)),
		),
	)

	// Delete request initiator from list of Rso members for provided Rso
	if err := runQuery(hs, deleteRsoMemberQuery, nil); err != nil {
		return err
	}

	// Retrieve details for Rso after user removal
	rsos := []Rso{}
	readRsoQuery := ReadRsos().WHERE(t.Rso.ID.EQ(pg.String(user.Rsomember.RsoID)))
	if err := runQuery(hs, readRsoQuery, &rsos); err != nil {
		return err
	}

	// Error if could not retrieve specific Rso
	if len(rsos) == 0 {
		return errors.New("delete rso member: cannot retrieve specified rso for viewing")
	}

	// Get user from database after removal of their rso membership
	userAfterDeletion := User{}
	if err := hs.GetUser(&userAfterDeletion); err != nil {
		return err
	}

	// Respond to clinet the UI for viewing details about an Rso
	hs.Local.RespondHtml(app.RsoInfo(rsos[0], userAfterDeletion))

	return nil
})
