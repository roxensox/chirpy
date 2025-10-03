package chirpyserver

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) FServerHits(writer http.ResponseWriter, req *http.Request) {
	// Handles hit to /metrics

	// Builds body string using string formatting
	body := fmt.Sprintf("<html>\n<body>\n\t<h1>Welcome, Chirpy Admin</h1>\n\t<p>Chirpy has been visited %d times!</p>\n</body>\n</html>", cfg.FileserverHits.Load())

	// Writes the request
	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte(body))
}
