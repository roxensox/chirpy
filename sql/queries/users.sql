-- name: CreateUser :one
INSERT INTO users (
	id, 
	created_at, 
	updated_at, 
	hashed_password,
	email
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
) RETURNING id, created_at, updated_at, email;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * 
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET 
	email = $1,
	hashed_password = $2,
	updated_at = $3
WHERE id = $4
RETURNING id, email, created_at, updated_at;
