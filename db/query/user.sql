-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- -- name: ListUsers :many
-- SELECT * FROM users
-- ORDER BY id
-- LIMIT $1 OFFSET $2;

-- -- name: UpdateUser :one
-- UPDATE users
-- set full_name = $2, email = $3, hashed_password = $4
-- WHERE id = $1 RETURNING *;

-- -- name: DeleteUser :exec
-- DELETE FROM users
-- WHERE id = $1;