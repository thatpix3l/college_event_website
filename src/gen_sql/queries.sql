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
-- name: ReadUniversities :many
SELECT *
FROM University;
-- name: ReadUserEmail :one
SELECT *
FROM BaseUser
WHERE email = $1;
-- name: ReadUsers :many
SELECT *
FROM BaseUser;
-- name: CreateSuperAdmin :one
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
    RETURNING *
)
INSERT INTO SuperAdmin(id)
VALUES (
        (
            SELECT id
            FROM base_user
        )
    )
RETURNING *;
-- name: CreateStudent :one
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
    RETURNING *
)
INSERT INTO Student(id, university)
VALUES (
        (
            SELECT id
            FROM base_user
        ),
        $5
    )
RETURNING *;
-- name: ReadStudents :many
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
    AND S.university = U.id;
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
-- name: ReadEvents :many
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
    FULL OUTER JOIN RsoEvent RE ON BE.id = RE.id;