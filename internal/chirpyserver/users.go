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
		Email:       dbResp.Email,
		CreatedAt:   dbResp.CreatedAt,
		UpdatedAt:   dbResp.UpdatedAt,
		ID:          dbResp.ID,
		IsChirpyRed: dbResp.IsChirpyRed,
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

func (cfg *ApiConfig) PUTUsers(writer http.ResponseWriter, req *http.Request) {
	// Handles PUT requests at users endpoint, takes in new email and password

	// Gets access token from request header
	tkn, err := auth.GetBearerToken(req.Header)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte("Could not find token"))
		return
	}

	// Gets the user ID by validating access token
	UID, err := auth.ValidateJWT(tkn, cfg.Secret)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte("Invalid token"))
		return
	}

	// Prepares object to receive request body
	rcv := struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}

	// Decodes request body into object
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&rcv)

	// Hashes the new password
	hashed, err := auth.HashPassword(rcv.Password)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Unable to hash password"))
		return
	}

	// Prepares parameteres object for the query
	params := database.UpdateUserParams{
		Email:          rcv.Email,
		HashedPassword: hashed,
		UpdatedAt:      time.Now().UTC(),
		ID:             UID,
	}

	// Runs the query
	resp, err := cfg.DBConn.UpdateUser(req.Context(), params)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to update user"))
		return
	}

	// Transfers query response to JSON-able object
	userObj := User{
		Email:       resp.Email,
		ID:          resp.ID,
		UpdatedAt:   resp.UpdatedAt,
		CreatedAt:   resp.CreatedAt,
		IsChirpyRed: resp.IsChirpyRed,
	}

	// Marshals object to JSON
	userJson, err := json.Marshal(userObj)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to marshal output"))
		return
	}

	// Writes success response
	writer.WriteHeader(200)
	writer.Write(userJson)
}
