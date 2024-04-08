-- name: CreateUniversity :one
INSERT INTO University (title, latitude, longitude, about)
VALUES (
        @title,
        @latitude,
        @longitude,
        @about
    )
RETURNING *;
-- name: CreateBaseUser :one
INSERT INTO BaseUser (name_first, name_last, email, password_hash)
VALUES (
        $1,
        $2,
        $3,
        $4
    )
RETURNING *;
-- name: CreateStudent :one
WITH base_user AS (
    INSERT INTO BaseUser (name_first, name_last, email, password_hash)
    VALUES ($1, $2, $3, $4)
    RETURNING *
)
INSERT INTO Student (id, university_id)
VALUES (
        (
            SELECT id
            FROM base_user
        ),
        @university_id
    )
RETURNING *;
-- name: CreateRso :one
INSERT INTO Rso
VALUES (DEFAULT, @title, @about, @university_id)
RETURNING *;
-- name: CreateRsoMember :one
INSERT INTO RsoMember
VALUES (DEFAULT, @rso_id, @is_admin)
RETURNING *;
-- name: CreateBaseEvent :one
INSERT INTO BaseEvent
VALUES (
        DEFAULT,
        @title,
        @about,
        @university_id,
        @start_time,
        @contact_phone,
        @contact_email,
        @latitude,
        @longitude
    )
RETURNING *;
-- name: CreateTag :one
INSERT INTO Tag (title)
VALUES ($1)
RETURNING *;
-- name: CreateTaggedEvent :one
INSERT INTO TaggedEvent (tag_id, base_event_id)
VALUES ($1, $2)
RETURNING *;
-- name: CreateTaggedRso :one
INSERT INTO TaggedRso (tag_id, rso_id)
VALUES ($1, $2)
RETURNING *;
-- name: CreateComment :one
INSERT INTO Comment (body, student_id, base_event_id)
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateRating :one
INSERT INTO Rating (stars, student_id, base_event_id)
VALUES ($1, $2, $3)
RETURNING *;
-- name: ReadEvents :many
SELECT sqlc.embed(BaseEvent),
    sqlc.embed(PublicEvent),
    sqlc.embed(PrivateEvent),
    sqlc.embed(RsoEvent)
FROM BaseEvent
    LEFT JOIN PublicEvent ON BaseEvent.id = PublicEvent.id
    LEFT JOIN PrivateEvent ON BaseEvent.id = PrivateEvent.id
    LEFT JOIN RsoEvent ON BaseEvent.id = RsoEvent.id;
-- name: ReadPublicEvents :many
SELECT sqlc.embed(BaseEvent),
    sqlc.embed(PublicEvent)
FROM BaseEvent
    JOIN PublicEvent ON BaseEvent.id = PublicEvent.id;
-- name: ReadPrivateEvents :many
SELECT sqlc.embed(BaseEvent),
    sqlc.embed(PrivateEvent)
FROM BaseEvent
    JOIN PrivateEvent ON BaseEvent.id = PrivateEvent.id;
-- name: ReadRsoEvents :many
SELECT sqlc.embed(BaseEvent),
    sqlc.embed(RsoEvent)
FROM BaseEvent
    JOIN RsoEvent ON BaseEvent.id = RsoEvent.id;
-- name: ReadUniversities :many
SELECT *
FROM University;
-- name: ReadBaseUsers :many
SELECT sqlc.embed(BaseUser),
    sqlc.embed(StudentVW),
    sqlc.embed(SuperAdminVW),
    sqlc.embed(RsoMemberVW)
FROM BaseUser
    FULL OUTER JOIN StudentVW ON BaseUser.id = StudentVW.id
    FULL OUTER JOIN SuperAdminVW ON BaseUser.id = SuperAdminVW.id
    FULL OUTER JOIN RsoMemberVW ON BaseUser.id = RsoMemberVW.id;
-- name: ReadStudents :many
SELECT sqlc.embed(StudentVW)
FROM StudentVW;
-- name: ReadSuperAdmins :many
SELECT *
FROM SuperAdminVW;
-- name: ReadRsoMembers :many
SELECT sqlc.embed(BaseUser),
    sqlc.embed(RsoMember)
FROM BaseUser
    JOIN RsoMember ON BaseUser.id = RsoMember.id;