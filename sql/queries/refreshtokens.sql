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

-- name: GetToken :one
SELECT * 
FROM refresh_tokens
WHERE token = $1
AND expires_at > $2
AND revoked_at IS NULL;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET 
	revoked_at = $2,
	updated_at = $2
WHERE token = $1;
