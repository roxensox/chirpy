package chirpyserver

import (
	"github.com/google/uuid"
	"github.com/roxensox/chirpy/internal/database"
	"sync/atomic"
	"time"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DBConn         *database.Queries
	Secret         string
}

type ValidateResponse struct {
	Valid       bool   `json:"valid"`
	Error       string `json:"error"`
	CleanedBody string `json:"cleaned_body"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

type User struct {
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ID           uuid.UUID `json:"id"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}
