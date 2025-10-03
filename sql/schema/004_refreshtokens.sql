-- +goose Up
CREATE TABLE refresh_tokens (
	token TEXT NOT NULL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	revoked_at TIMESTAMP,
	user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_tokens;
