-- name: CreateUser :one
INSERT INTO users (id, email, name, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET name = $2, updated_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
