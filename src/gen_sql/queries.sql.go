// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package gen_sql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createBaseEvent = `-- name: CreateBaseEvent :one
INSERT INTO BaseEvent (
        title,
        body,
        university,
        occurrence_time,
        occurrence_location,
        contact_phone,
        contact_email
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, title, body, university, occurrence_time, occurrence_location, contact_phone, contact_email
`

type CreateBaseEventParams struct {
	Title              string           `schema:",required"`
	Body               string           `schema:",required"`
	University         int32            `schema:",required"`
	OccurrenceTime     pgtype.Timestamp `schema:",required"`
	OccurrenceLocation int32            `schema:",required"`
	ContactPhone       string           `schema:",required"`
	ContactEmail       string           `schema:",required"`
}

func (q *Queries) CreateBaseEvent(ctx context.Context, arg CreateBaseEventParams) (Baseevent, error) {
	row := q.db.QueryRow(ctx, createBaseEvent,
		arg.Title,
		arg.Body,
		arg.University,
		arg.OccurrenceTime,
		arg.OccurrenceLocation,
		arg.ContactPhone,
		arg.ContactEmail,
	)
	var i Baseevent
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Body,
		&i.University,
		&i.OccurrenceTime,
		&i.OccurrenceLocation,
		&i.ContactPhone,
		&i.ContactEmail,
	)
	return i, err
}

const createComment = `-- name: CreateComment :one
INSERT INTO Comment (body, posted_by, base_event)
VALUES ($1, $2, $3)
RETURNING id, body, posted_by, base_event
`

type CreateCommentParams struct {
	Body      string      `schema:",required"`
	PostedBy  pgtype.Int4 `schema:",required"`
	BaseEvent int32       `schema:",required"`
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRow(ctx, createComment, arg.Body, arg.PostedBy, arg.BaseEvent)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Body,
		&i.PostedBy,
		&i.BaseEvent,
	)
	return i, err
}

const createCoordinate = `-- name: CreateCoordinate :one
INSERT INTO Coordinate (title, latitude, longitude)
VALUES ($1, $2, $3)
RETURNING id, title, latitude, longitude
`

type CreateCoordinateParams struct {
	Title     string  `schema:",required"`
	Latitude  float64 `schema:",required"`
	Longitude float64 `schema:",required"`
}

func (q *Queries) CreateCoordinate(ctx context.Context, arg CreateCoordinateParams) (Coordinate, error) {
	row := q.db.QueryRow(ctx, createCoordinate, arg.Title, arg.Latitude, arg.Longitude)
	var i Coordinate
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Latitude,
		&i.Longitude,
	)
	return i, err
}

const createPrivateEvent = `-- name: CreatePrivateEvent :one
INSERT INTO PrivateEvent (id)
VALUES ($1)
RETURNING id
`

func (q *Queries) CreatePrivateEvent(ctx context.Context, id int32) (int32, error) {
	row := q.db.QueryRow(ctx, createPrivateEvent, id)
	err := row.Scan(&id)
	return id, err
}

const createPublicEvent = `-- name: CreatePublicEvent :one
INSERT INTO PublicEvent (id)
VALUES ($1)
RETURNING id, approved
`

func (q *Queries) CreatePublicEvent(ctx context.Context, id int32) (Publicevent, error) {
	row := q.db.QueryRow(ctx, createPublicEvent, id)
	var i Publicevent
	err := row.Scan(&i.ID, &i.Approved)
	return i, err
}

const createRating = `-- name: CreateRating :one
INSERT INTO Rating (stars, posted_by, base_event)
VALUES ($1, $2, $3)
RETURNING id, stars, posted_by, base_event
`

type CreateRatingParams struct {
	Stars     int32       `schema:",required"`
	PostedBy  pgtype.Int4 `schema:",required"`
	BaseEvent int32       `schema:",required"`
}

func (q *Queries) CreateRating(ctx context.Context, arg CreateRatingParams) (Rating, error) {
	row := q.db.QueryRow(ctx, createRating, arg.Stars, arg.PostedBy, arg.BaseEvent)
	var i Rating
	err := row.Scan(
		&i.ID,
		&i.Stars,
		&i.PostedBy,
		&i.BaseEvent,
	)
	return i, err
}

const createRso = `-- name: CreateRso :one
INSERT INTO Rso (title, university)
VALUES ($1, $2)
RETURNING id, title, university
`

type CreateRsoParams struct {
	Title      string `schema:",required"`
	University int32  `schema:",required"`
}

func (q *Queries) CreateRso(ctx context.Context, arg CreateRsoParams) (Rso, error) {
	row := q.db.QueryRow(ctx, createRso, arg.Title, arg.University)
	var i Rso
	err := row.Scan(&i.ID, &i.Title, &i.University)
	return i, err
}

const createRsoEvent = `-- name: CreateRsoEvent :one
INSERT INTO RsoEvent (id, rso)
VALUES ($1, $2)
RETURNING id, rso
`

type CreateRsoEventParams struct {
	ID  int32 `schema:",required"`
	Rso int32 `schema:",required"`
}

func (q *Queries) CreateRsoEvent(ctx context.Context, arg CreateRsoEventParams) (Rsoevent, error) {
	row := q.db.QueryRow(ctx, createRsoEvent, arg.ID, arg.Rso)
	var i Rsoevent
	err := row.Scan(&i.ID, &i.Rso)
	return i, err
}

const createRsoMember = `-- name: CreateRsoMember :one
INSERT INTO RsoMember (rso)
VALUES ($1)
RETURNING id, rso, is_admin
`

func (q *Queries) CreateRsoMember(ctx context.Context, rso int32) (Rsomember, error) {
	row := q.db.QueryRow(ctx, createRsoMember, rso)
	var i Rsomember
	err := row.Scan(&i.ID, &i.Rso, &i.IsAdmin)
	return i, err
}

const createStudent = `-- name: CreateStudent :one
WITH base_user AS (
    INSERT INTO BaseUser (
            name_first,
            name_last,
            email,
            password_hash
        )
    VALUES (
            $1,
            $2,
            $3,
            $4
        )
    RETURNING id, name_first, name_last, email, password_hash
)
INSERT INTO Student(id, university)
VALUES (
        (
            SELECT id
            FROM base_user
        ),
        $5
    )
RETURNING id, university
`

type CreateStudentParams struct {
	NameFirst    string `schema:",required"`
	NameLast     string `schema:",required"`
	Email        string `schema:",required"`
	PasswordHash string `schema:",required"`
	University   int32  `schema:",required"`
}

func (q *Queries) CreateStudent(ctx context.Context, arg CreateStudentParams) (Student, error) {
	row := q.db.QueryRow(ctx, createStudent,
		arg.NameFirst,
		arg.NameLast,
		arg.Email,
		arg.PasswordHash,
		arg.University,
	)
	var i Student
	err := row.Scan(&i.ID, &i.University)
	return i, err
}

const createSuperAdmin = `-- name: CreateSuperAdmin :one
WITH base_user AS (
    INSERT INTO BaseUser (
            name_first,
            name_last,
            email,
            password_hash
        )
    VALUES (
            $1,
            $2,
            $3,
            $4
        )
    RETURNING id, name_first, name_last, email, password_hash
)
INSERT INTO SuperAdmin(id)
VALUES (
        (
            SELECT id
            FROM base_user
        )
    )
RETURNING id
`

type CreateSuperAdminParams struct {
	NameFirst    string `schema:",required"`
	NameLast     string `schema:",required"`
	Email        string `schema:",required"`
	PasswordHash string `schema:",required"`
}

func (q *Queries) CreateSuperAdmin(ctx context.Context, arg CreateSuperAdminParams) (int32, error) {
	row := q.db.QueryRow(ctx, createSuperAdmin,
		arg.NameFirst,
		arg.NameLast,
		arg.Email,
		arg.PasswordHash,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createTag = `-- name: CreateTag :one
INSERT INTO Tag (title)
VALUES ($1)
RETURNING id, title
`

func (q *Queries) CreateTag(ctx context.Context, title string) (Tag, error) {
	row := q.db.QueryRow(ctx, createTag, title)
	var i Tag
	err := row.Scan(&i.ID, &i.Title)
	return i, err
}

const createTaggedEvent = `-- name: CreateTaggedEvent :one
INSERT INTO TaggedEvent (tag, base_event)
VALUES ($1, $2)
RETURNING tag, base_event
`

type CreateTaggedEventParams struct {
	Tag       int32 `schema:",required"`
	BaseEvent int32 `schema:",required"`
}

func (q *Queries) CreateTaggedEvent(ctx context.Context, arg CreateTaggedEventParams) (Taggedevent, error) {
	row := q.db.QueryRow(ctx, createTaggedEvent, arg.Tag, arg.BaseEvent)
	var i Taggedevent
	err := row.Scan(&i.Tag, &i.BaseEvent)
	return i, err
}

const createTaggedRso = `-- name: CreateTaggedRso :one
INSERT INTO TaggedRso (tag, rso)
VALUES ($1, $2)
RETURNING tag, rso
`

type CreateTaggedRsoParams struct {
	Tag int32 `schema:",required"`
	Rso int32 `schema:",required"`
}

func (q *Queries) CreateTaggedRso(ctx context.Context, arg CreateTaggedRsoParams) (Taggedrso, error) {
	row := q.db.QueryRow(ctx, createTaggedRso, arg.Tag, arg.Rso)
	var i Taggedrso
	err := row.Scan(&i.Tag, &i.Rso)
	return i, err
}

const createUniversity = `-- name: CreateUniversity :one
WITH new_coord AS (
    INSERT INTO Coordinate (title, latitude, longitude)
    VALUES ($3, $4, $5)
    RETURNING id, title, latitude, longitude
)
INSERT INTO University (title, coordinate, about)
VALUES (
        $1,
        (
            SELECT id
            FROM new_coord
        ),
        $2
    )
RETURNING id, title, coordinate, about
`

type CreateUniversityParams struct {
	UniversityTitle string  `schema:",required"`
	UniversityAbout string  `schema:",required"`
	CoordTitle      string  `schema:",required"`
	CoordLatitude   float64 `schema:",required"`
	CoordLongitude  float64 `schema:",required"`
}

func (q *Queries) CreateUniversity(ctx context.Context, arg CreateUniversityParams) (University, error) {
	row := q.db.QueryRow(ctx, createUniversity,
		arg.UniversityTitle,
		arg.UniversityAbout,
		arg.CoordTitle,
		arg.CoordLatitude,
		arg.CoordLongitude,
	)
	var i University
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Coordinate,
		&i.About,
	)
	return i, err
}

const readEvents = `-- name: ReadEvents :many
SELECT BE.id,
    BE.title,
    BE.body,
    BE.university,
    BE.occurrence_time,
    BE.contact_phone,
    BE.contact_email,
    C.title AS coord_title,
    C.latitude,
    C.longitude,
    PUE.id AS public_event,
    PUE.approved,
    PRE.id AS private_event,
    RE.id AS rso_event,
    RE.rso AS rso
FROM BaseEvent BE
    FULL OUTER JOIN Coordinate C ON BE.occurrence_location = C.id
    FULL OUTER JOIN PublicEvent PUE ON BE.id = PUE.id
    FULL OUTER JOIN PrivateEvent PRE ON BE.id = PRE.id
    FULL OUTER JOIN RsoEvent RE ON BE.id = RE.id
`

type ReadEventsRow struct {
	ID             pgtype.Int4      `schema:",required"`
	Title          pgtype.Text      `schema:",required"`
	Body           pgtype.Text      `schema:",required"`
	University     pgtype.Int4      `schema:",required"`
	OccurrenceTime pgtype.Timestamp `schema:",required"`
	ContactPhone   pgtype.Text      `schema:",required"`
	ContactEmail   pgtype.Text      `schema:",required"`
	CoordTitle     pgtype.Text      `schema:",required"`
	Latitude       pgtype.Float8    `schema:",required"`
	Longitude      pgtype.Float8    `schema:",required"`
	PublicEvent    pgtype.Int4      `schema:",required"`
	Approved       pgtype.Bool      `schema:",required"`
	PrivateEvent   pgtype.Int4      `schema:",required"`
	RsoEvent       pgtype.Int4      `schema:",required"`
	Rso            pgtype.Int4      `schema:",required"`
}

func (q *Queries) ReadEvents(ctx context.Context) ([]ReadEventsRow, error) {
	rows, err := q.db.Query(ctx, readEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadEventsRow
	for rows.Next() {
		var i ReadEventsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Body,
			&i.University,
			&i.OccurrenceTime,
			&i.ContactPhone,
			&i.ContactEmail,
			&i.CoordTitle,
			&i.Latitude,
			&i.Longitude,
			&i.PublicEvent,
			&i.Approved,
			&i.PrivateEvent,
			&i.RsoEvent,
			&i.Rso,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readStudents = `-- name: ReadStudents :many
SELECT S.id,
    BU.name_first,
    BU.name_last,
    BU.email,
    BU.password_hash,
    U.title AS university_name
FROM Student S,
    BaseUser BU,
    University U
WHERE S.id = BU.id
    AND S.university = U.id
`

type ReadStudentsRow struct {
	ID             int32  `schema:",required"`
	NameFirst      string `schema:",required"`
	NameLast       string `schema:",required"`
	Email          string `schema:",required"`
	PasswordHash   string `schema:",required"`
	UniversityName string `schema:",required"`
}

func (q *Queries) ReadStudents(ctx context.Context) ([]ReadStudentsRow, error) {
	rows, err := q.db.Query(ctx, readStudents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadStudentsRow
	for rows.Next() {
		var i ReadStudentsRow
		if err := rows.Scan(
			&i.ID,
			&i.NameFirst,
			&i.NameLast,
			&i.Email,
			&i.PasswordHash,
			&i.UniversityName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readUniversities = `-- name: ReadUniversities :many
SELECT id, title, coordinate, about
FROM University
`

func (q *Queries) ReadUniversities(ctx context.Context) ([]University, error) {
	rows, err := q.db.Query(ctx, readUniversities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []University
	for rows.Next() {
		var i University
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Coordinate,
			&i.About,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readUserEmail = `-- name: ReadUserEmail :one
SELECT id, name_first, name_last, email, password_hash
FROM BaseUser
WHERE email = $1
`

func (q *Queries) ReadUserEmail(ctx context.Context, email string) (Baseuser, error) {
	row := q.db.QueryRow(ctx, readUserEmail, email)
	var i Baseuser
	err := row.Scan(
		&i.ID,
		&i.NameFirst,
		&i.NameLast,
		&i.Email,
		&i.PasswordHash,
	)
	return i, err
}

const readUsers = `-- name: ReadUsers :many
SELECT id, name_first, name_last, email, password_hash
FROM BaseUser
`

func (q *Queries) ReadUsers(ctx context.Context) ([]Baseuser, error) {
	rows, err := q.db.Query(ctx, readUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Baseuser
	for rows.Next() {
		var i Baseuser
		if err := rows.Scan(
			&i.ID,
			&i.NameFirst,
			&i.NameLast,
			&i.Email,
			&i.PasswordHash,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
