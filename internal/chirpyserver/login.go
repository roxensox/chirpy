package chirpyserver

import (
	"encoding/json"
	"github.com/roxensox/chirpy/internal/auth"
	"github.com/roxensox/chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *ApiConfig) POSTLogin(writer http.ResponseWriter, req *http.Request) {
	// Handles post request to login endpoint

	// Creates an anonymous struct instance to receive input
	inObj := struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}

	// Creates a new decoder and decodes the request body into struct
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&inObj)

	// Queries the user from the database
	user, err := cfg.DBConn.GetUserByEmail(req.Context(), inObj.Email)
	// Validates the password
	validPass, err2 := auth.CheckPasswordHash(inObj.Password, user.HashedPassword)

	if err != nil || !validPass {
		writer.WriteHeader(401)
		writer.Write([]byte("Incorrect email or password"))
		return
	}

	if err2 != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to compare passwords"))
		return
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.Secret, time.Hour)

	ref_token, err := auth.MakeRefreshToken()
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to generate refresh token"))
		return
	}

	refTokenParams := database.AddRefreshTokenParams{
		UserID:    user.ID,
		Token:     ref_token,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
	}

	err = cfg.DBConn.AddRefreshToken(req.Context(), refTokenParams)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to add refresh token"))
		return
	}

	out := User{
		Email:        user.Email,
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Token:        jwt,
		RefreshToken: ref_token,
		IsChirpyRed:  user.IsChirpyRed,
	}

	outJson, err := json.Marshal(out)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Something went wrong"))
		return
	}
	writer.WriteHeader(200)
	writer.Write(outJson)
}
