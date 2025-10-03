package chirpyserver

import (
	"database/sql"
	"github.com/roxensox/chirpy/internal/auth"
	"github.com/roxensox/chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *ApiConfig) POSTRevoke(writer http.ResponseWriter, req *http.Request) {
	tkn, err := auth.GetBearerToken(req.Header)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte("Token not found"))
		return
	}
	params := database.RevokeTokenParams{
		Token: tkn,
		RevokedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}
	err = cfg.DBConn.RevokeToken(req.Context(), params)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to revoke token"))
		return
	}
	writer.WriteHeader(204)
}
