-- name: CreateProfile :exec
INSERT INTO profiles (id)
VALUES ($1)
ON CONFLICT (id) DO NOTHING;

-- name: UpdateLastLogin :exec
UPDATE profiles
SET last_login_at = now()
WHERE id = $1;

-- name: GetProfile :one
SELECT * FROM profiles
WHERE id = $1;