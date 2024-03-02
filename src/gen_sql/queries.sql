-- name: CreateCoordinate :one
INSERT INTO Coordinate (title, latitude, longitude)
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateUniversity :one
INSERT INTO University (title, coordinate, about)
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateUniversityMember :one
INSERT INTO UniversityMember (
        university,
        name_first,
        name_middle,
        name_last,
        email,
        password_hash
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: CreateSuperAdmin :one
INSERT INTO SuperAdmin (id)
VALUES ($1)
RETURNING *;
-- name: CreateUniversityEvent :one
INSERT INTO UniversityEvent (
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
INSERT INTO PublicEvent (id, approved)
VALUES ($1, $2)
RETURNING *;
-- name: CreatePrivateEvent :one
INSERT INTO PrivateEvent (id)
VALUES ($1)
RETURNING *;
-- name: CreateRso :one
INSERT INTO Rso (title, university)
VALUES ($1, $2)
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
INSERT INTO TaggedEvent (tag, university_event)
VALUES ($1, $2)
RETURNING *;
-- name: CreateTaggedRso :one
INSERT INTO TaggedRso (tag, rso)
VALUES ($1, $2)
RETURNING *;
-- name: CreateRsoMember :one
INSERT INTO RsoMember (rso, university_member, is_admin)
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateComment :one
INSERT INTO Comment (body, posted_by, university_event)
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateRating :one
INSERT INTO Rating (stars, posted_by, university_event)
VALUES ($1, $2, $3)
RETURNING *;