-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
  username, encrypted_password, role
) VALUES (
  $1, $2, $3::user_role
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = $2, encrypted_password = $3, role = $4::user_role
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;