package chirpyserver

import (
	"encoding/json"
	"github.com/roxensox/chirpy/internal/auth"
	"github.com/roxensox/chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *ApiConfig) POSTRefresh(writer http.ResponseWriter, req *http.Request) {
	tkn, err := auth.GetBearerToken(req.Header)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte("Token not found in header"))
	}
	params := database.GetTokenParams{
		Token:     tkn,
		ExpiresAt: time.Now().UTC(),
	}
	resp, err := cfg.DBConn.GetToken(req.Context(), params)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte("Invalid token"))
		return
	}

	accTkn, err := auth.MakeJWT(resp.UserID, cfg.Secret, 1*time.Hour)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to create new access token"))
		return
	}

	respObj := struct {
		Token string `json:"token"`
	}{
		Token: accTkn,
	}
	respJson, err := json.Marshal(respObj)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to marshal response to JSON"))
		return
	}
	writer.WriteHeader(200)
	writer.Write(respJson)
}
