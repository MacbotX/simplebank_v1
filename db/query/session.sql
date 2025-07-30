-- name: CreateSession :one
INSERT INTO sessions (
  id,
  username,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- -- name: ListSessions :many
-- SELECT * FROM sessions
-- ORDER BY id
-- LIMIT $1 OFFSET $2;

-- -- name: UpdateSession :one
-- UPDATE sessions
-- set full_name = $2, email = $3, hashed_password = $4
-- WHERE id = $1 RETURNING *;

-- -- name: DeleteSession :exec
-- DELETE FROM sessions
-- WHERE id = $1;