-- name: CreateContact :one
INSERT INTO contacts (
  user_id,
  name,
  phone_number,
  image_uri,
  type,
  is_pinned
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetContacts :many
SELECT * FROM contacts
WHERE user_id = $1
ORDER BY is_pinned DESC, name ASC;

-- name: GetContactById :one
SELECT * FROM contacts
WHERE id = $1 AND user_id = $2;

-- name: UpdateContact :one
UPDATE contacts
SET
  name = COALESCE(sqlc.narg('name'), name),
  phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
  image_uri = COALESCE(sqlc.narg('image_uri'), image_uri),
  type = COALESCE(sqlc.narg('type'), type),
  is_pinned = COALESCE(sqlc.narg('is_pinned'), is_pinned),
  updated_at = now()
WHERE id = sqlc.arg('id') AND user_id = sqlc.arg('user_id')
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contacts
WHERE id = $1 AND user_id = $2;

-- name: ToggleContactType :one
UPDATE contacts
SET
    type = CASE WHEN type = 'Trusted' THEN 'Casual' ELSE 'Trusted' END,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: ToggleContactPin :one
UPDATE contacts
SET
    is_pinned = NOT is_pinned,
    updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;
