package chirpyserver

import (
	"encoding/json"
	"github.com/roxensox/chirpy/internal/auth"
	"github.com/roxensox/chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *ApiConfig) POSTRefresh(writer http.ResponseWriter, req *http.Request) {
	// Handles POST request to refresh endpoint, authenticates and returns a new access token

	// Gets the user's refresh token
	tkn, err := auth.GetBearerToken(req.Header)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte("Token not found in header"))
	}

	// Builds parameters for refresh token DB query
	params := database.GetTokenParams{
		Token:     tkn,
		ExpiresAt: time.Now().UTC(),
	}

	// Runs the query
	resp, err := cfg.DBConn.GetToken(req.Context(), params)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte("Invalid token"))
		return
	}

	// Creates a new access token if the user validated successfully
	accTkn, err := auth.MakeJWT(resp.UserID, cfg.Secret, 1*time.Hour)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to create new access token"))
		return
	}

	// Builds an output object
	respObj := struct {
		Token string `json:"token"`
	}{
		Token: accTkn,
	}

	// Marshals output object to json
	respJson, err := json.Marshal(respObj)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to marshal response to JSON"))
		return
	}

	// Returns output json with success code
	writer.WriteHeader(200)
	writer.Write(respJson)
}
