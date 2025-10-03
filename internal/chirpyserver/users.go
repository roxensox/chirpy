package chirpyserver

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/roxensox/chirpy/internal/auth"
	"github.com/roxensox/chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *ApiConfig) POSTUsers(writer http.ResponseWriter, req *http.Request) {
	// Handles POST request to users endpoint, returns newly created user

	// Creates anonymous struct instance for receiving input and umarshals into it
	in := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&in)
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to hash password"))
	}

	// Builds query param object
	params := database.CreateUserParams{
		Email:          in.Email,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		HashedPassword: hash,
		ID:             uuid.New(),
	}

	// Queries database to insert new data
	dbResp, err := cfg.DBConn.CreateUser(req.Context(), params)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to find user"))
		return
	}

	// Builds new User object to umarshal to JSON
	jsonResp := User{
		Email:     dbResp.Email,
		CreatedAt: dbResp.CreatedAt,
		UpdatedAt: dbResp.UpdatedAt,
		ID:        dbResp.ID,
	}

	// Umarshals User object
	resp, err := json.Marshal(jsonResp)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to marshal results"))
		return
	}

	// Writes success response
	writer.WriteHeader(201)
	writer.Write(resp)
}
