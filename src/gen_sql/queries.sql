-- name: CreateUniversity :one
INSERT INTO University (title, about, latitude, longitude)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: CreateBaseUser :one
INSERT INTO BaseUser(
        name_first,
        name_last,
        email,
        password_hash,
        is_super_admin,
        university
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: CreateRso :one
INSERT INTO Rso (title, about, university)
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateRsoMember :one
INSERT INTO RsoMember (id, rso)
VALUES ($1, $2)
RETURNING *;
-- name: CreateBaseEvent :one
INSERT INTO BaseEvent (
        title,
        about,
        university,
        start_time,
        contact_phone,
        contact_email,
        event_type,
        latitude,
        longitude
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;
-- name: CreateTag :one
INSERT INTO Tag (title)
VALUES ($1)
RETURNING *;
-- name: CreateTaggedEvent :one
INSERT INTO TaggedEvent (tag, base_event)
VALUES ($1, $2)
RETURNING *;
-- name: CreateTaggedRso :one
INSERT INTO TaggedRso (tag, rso)
VALUES ($1, $2)
RETURNING *;
-- name: CreateComment :one
INSERT INTO Comment (body, posted_by, base_event)
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateRating :one
INSERT INTO Rating (stars, posted_by, base_event)
VALUES ($1, $2, $3)
RETURNING *;
-- name: ReadEvents :many
SELECT *
FROM BaseEvent;
-- name: ReadUniversities :many
SELECT *
FROM University;
-- name: ReadBaseUsers :many
WITH uni AS (
    SELECT id AS university_id,
        title AS university_title,
        latitude AS university_latitude,
        longitude AS university_longitude,
        about AS university_about
    FROM University
)
SELECT *
FROM BaseUser
    LEFT JOIN uni ON BaseUser.university = uni.university_id;