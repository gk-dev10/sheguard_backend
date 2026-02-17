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

-- name: UpdateProfile :one
UPDATE profiles
SET
  name = COALESCE($2, name),
  phone_number = COALESCE($3, phone_number),
  profile_image_url = COALESCE($4, profile_image_url),
  blood_group = COALESCE($5, blood_group),
  allergies = COALESCE($6, allergies),
  medications = COALESCE($7, medications),
  updated_at = now()
WHERE id = $1
RETURNING *;