package chirpyserver

import (
	"net/http"
)

func Healthz(writer http.ResponseWriter, req *http.Request) {
	// Manually writes a response for the /healthz endpoint

	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}
