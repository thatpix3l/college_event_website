-- name: CreateCoordinate :one
INSERT INTO Coordinate (title, latitude, longitude)
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateUniversity :one
WITH new_coord AS (
    INSERT INTO Coordinate (title, latitude, longitude)
    VALUES (@coord_title, @coord_latitude, @coord_longitude)
    RETURNING *
)
INSERT INTO University (title, coordinate, about)
VALUES (
        @university_title,
        (
            SELECT id
            FROM new_coord
        ),
        @university_about
    )
RETURNING *;
-- name: CreateBaseUser :one
INSERT INTO BaseUser (
        name_first,
        name_middle,
        name_last,
        email,
        password_hash
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
RETURNING *;
-- name: CreateSuperAdmin :one
INSERT INTO SuperAdmin (id)
VALUES (@id)
RETURNING *;
-- name: CreateMember :one
INSERT INTO Member (id, university)
VALUES ($1, $2)
RETURNING *;
-- name: CreateRso :one
INSERT INTO Rso (title, university)
VALUES ($1, $2)
RETURNING *;
-- name: CreateRsoMember :one
INSERT INTO RsoMember (rso)
VALUES ($1)
RETURNING *;
-- name: CreateBaseEvent :one
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
RETURNING *;
-- name: CreatePublicEvent :one
INSERT INTO PublicEvent (id)
VALUES ($1)
RETURNING *;
-- name: CreatePrivateEvent :one
INSERT INTO PrivateEvent (id)
VALUES ($1)
RETURNING *;
-- name: CreateRsoEvent :one
INSERT INTO RsoEvent (id, rso)
VALUES ($1, $2)
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