-- name: CreateUser :exec
INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: UpdateUser :exec
UPDATE users SET name=$2, email=$3, password=$4, updated_at=$5 WHERE id=$1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;