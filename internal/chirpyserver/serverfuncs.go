package chirpyserver

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	// Middleware to track hits to the main server

	return http.HandlerFunc(
		// Newly defined handler function
		func(writer http.ResponseWriter, req *http.Request) {
			// Increments FileserverHits
			cfg.FileserverHits.Add(1)
			// Calls ServerHTTP method on input handler
			next.ServeHTTP(writer, req)
		},
	)
}

func removeProfanity(s string) string {
	// Replaces occurrences of flagged words in an input string and returns it

	// List of flagged words
	flagged := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	// Gets each individual word from the input string in a slice
	words := strings.Fields(s)

	for i, v := range words {
		// If the word is in the flagged list, replaces it with ****
		if slices.Contains(flagged, strings.ToLower(v)) {
			words[i] = "****"
		}
	}

	// Joins the words and returns the result
	return strings.Join(words, " ")
}

func ValidateChirp(writer http.ResponseWriter, req *http.Request) {
	// Receives a "chirp" and validates it by specified conditions

	// Creates instance of Chirp object and decodes request into it
	chirp := Chirp{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)

	// If that produced an error, sends an error 500 code
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		writer.WriteHeader(500)
	}

	// Ensures the input text is at most 140 characters
	if len(chirp.Body) <= 140 {
		// Instantiates an output object
		out := ValidateResponse{
			Valid: true,
		}

		// Cleans the message
		out.CleanedBody = removeProfanity(chirp.Body)

		// Marshals output object to json
		outjson, err := json.Marshal(out)

		// If an error occurs while marshalling, manually writes an error response and returns
		if err != nil {
			writer.WriteHeader(400)
			resp := []byte("{\"error\":\"something went wrong\"")
			writer.Write(resp)
			return
		}

		// Writes the response
		writer.WriteHeader(200)
		writer.Write(outjson)
	} else {
		// Writes an error response
		writer.WriteHeader(400)
		// Creates output object with error message
		out := ValidateResponse{
			Error: "chirp is too long",
		}

		// Marshals output object
		resp, err := json.Marshal(out)

		// If an error occurs while marshalling, manually writes an error response and returns
		if err != nil {
			resp := []byte("{\"error\":\"something went wrong\"")
			writer.Write(resp)
		}

		// Writes the response
		writer.Write(resp)
	}
}
