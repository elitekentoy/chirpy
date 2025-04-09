-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)
RETURNING *;


-- name: DeleteUsers :exec
DELETE FROM users;


-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET updated_at = $1, email = $2, hashed_password = $3
WHERE id = $4
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateChirpyRed :exec
UPDATE users SET is_chirpy_red = $1 WHERE id = $2;