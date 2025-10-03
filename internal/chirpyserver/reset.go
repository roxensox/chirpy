package chirpyserver

import (
	"net/http"
)

func (cfg *ApiConfig) Reset(writer http.ResponseWriter, req *http.Request) {
	// Handles hit to /reset (resets hit counter to 0)

	// Resets hit counter to 0
	cfg.FileserverHits.Store(0)

	// Deletes all records from the users table
	cfg.DBConn.ResetUsers(req.Context())

	// Writes the response
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("RESET"))
}
