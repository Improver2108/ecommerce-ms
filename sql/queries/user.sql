-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users(name,email,password_hsh,phone)
VALUES($1,$2,$3,$4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = $1, email = $2, password_hsh = $3, phone = $4
WHERE id = $5
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;