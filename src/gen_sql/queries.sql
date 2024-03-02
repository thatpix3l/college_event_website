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