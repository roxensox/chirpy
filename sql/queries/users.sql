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
