-- name: AddRefreshToken :exec
INSERT INTO refresh_tokens (
	token,
	created_at,
	updated_at,
	expires_at,
	user_id
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
);
