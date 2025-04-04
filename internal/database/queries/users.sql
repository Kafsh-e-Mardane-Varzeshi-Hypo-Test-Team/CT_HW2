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
  @username, @encrypted_password, @role::user_role
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = @username, 
    encrypted_password = @encrypted_password, 
    role = @role::user_role
WHERE id = @id
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;