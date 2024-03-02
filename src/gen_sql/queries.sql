-- name: CreateCoordinate :one
INSERT INTO Coordinate (title, latitude, longitude)
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateUniversity :one
INSERT INTO University (title, coordinate, about)
VALUES ($1, $2, $3)
RETURNING *;
-- name: CreateUniversityMember :one
INSERT INTO UniversityMember
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;