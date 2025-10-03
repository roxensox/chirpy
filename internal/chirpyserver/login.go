package chirpyserver

import (
	"encoding/json"
	"github.com/roxensox/chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *ApiConfig) POSTLogin(writer http.ResponseWriter, req *http.Request) {
	// Handles post request to login endpoint

	// Creates an anonymous struct instance to receive input
	inObj := struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
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

	t := func() time.Duration {
		secs := time.Duration(inObj.ExpiresInSeconds)
		if inObj.ExpiresInSeconds == 0 || secs*time.Second > time.Hour {
			return time.Hour
		}
		return secs * time.Second
	}()

	tkn, err := auth.MakeJWT(user.ID, cfg.Secret, t)

	out := User{
		Email:     user.Email,
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token:     tkn,
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
