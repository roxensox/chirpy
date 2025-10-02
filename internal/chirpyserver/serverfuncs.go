package chirpyserver

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/roxensox/chirpy/internal/database"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

func Healthz(writer http.ResponseWriter, req *http.Request) {
	// Manually writes a response for the /healthz endpoint
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	code, err := writer.Write([]byte("OK"))
	if err != nil {
		fmt.Printf("Failed to write with code: %v", code)
	}
}

func (cfg *ApiConfig) POSTChirps(writer http.ResponseWriter, req *http.Request) {
	type inChirp struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}
	inObj := inChirp{}
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&inObj)
	user_id, err := uuid.Parse(inObj.UserID)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to parse user ID"))
		return
	}

	params := database.CreateChirpParams{
		UserID:    user_id,
		Body:      inObj.Body,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ID:        uuid.New(),
	}

	dbResp, err := cfg.DBConn.CreateChirp(req.Context(), params)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("User not found"))
		return
	}
	outObj := Chirp{
		ID:        dbResp.ID,
		CreatedAt: dbResp.CreatedAt,
		UpdatedAt: dbResp.UpdatedAt,
		Body:      dbResp.Body,
		UserID:    dbResp.UserID,
	}
	outJson, err := json.Marshal(outObj)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to marshal data"))
		return
	}
	writer.WriteHeader(201)
	writer.Write(outJson)
}

func (cfg *ApiConfig) GETChirps(writer http.ResponseWriter, req *http.Request) {
	out := make([]Chirp, 0)
	allChirps, err := cfg.DBConn.GetChirps(req.Context())
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Unable to get chirps"))
	}
	for _, c := range allChirps {
		ctoJSON := Chirp{
			ID:        c.ID,
			CreatedAt: c.UpdatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
		out = append(out, ctoJSON)
	}
	outJson, err := json.Marshal(out)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to marshal data"))
	}
	writer.WriteHeader(200)
	writer.Write(outJson)
}

func (cfg *ApiConfig) GETChirpByID(writer http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to parse chirp ID"))
		return
	}
	dbResp, err := cfg.DBConn.GetExactChirp(req.Context(), chirpUUID)
	if err != nil {
		writer.WriteHeader(404)
		writer.Write([]byte("Chirp not found"))
		return
	}
	out := Chirp{
		ID:        dbResp.ID,
		CreatedAt: dbResp.CreatedAt,
		UpdatedAt: dbResp.UpdatedAt,
		Body:      dbResp.Body,
		UserID:    dbResp.UserID,
	}
	outJson, err := json.Marshal(out)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to marshal chirp data"))
		return
	}
	writer.WriteHeader(200)
	writer.Write(outJson)
}

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

func (cfg *ApiConfig) POSTUsers(writer http.ResponseWriter, req *http.Request) {
	type ReqPayload struct {
		Email string `json:"email"`
	}
	in := ReqPayload{}
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&in)

	params := database.CreateUserParams{
		Email:     in.Email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ID:        uuid.New(),
	}

	dbResp, err := cfg.DBConn.CreateUser(req.Context(), params)

	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to find user"))
		return
	}

	jsonResp := User{
		Email:     dbResp.Email,
		CreatedAt: dbResp.CreatedAt,
		UpdatedAt: dbResp.UpdatedAt,
		ID:        dbResp.ID,
	}

	resp, err := json.Marshal(jsonResp)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to marshal results"))
		return
	}

	writer.WriteHeader(201)
	writer.Write(resp)
}

func ValidateChirp(writer http.ResponseWriter, req *http.Request) {
	// Receives a "chirp" and validates it by specified conditions

	// Creates a payload struct to decode the json into
	type Chirp struct {
		Body string `json:"body"`
	}

	// Instantiates the payload and decodes the json into it
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

func (cfg *ApiConfig) FServerHits(writer http.ResponseWriter, req *http.Request) {
	// Handles hit to /metrics

	// Builds body string using string formatting
	body := fmt.Sprintf("<html>\n<body>\n\t<h1>Welcome, Chirpy Admin</h1>\n\t<p>Chirpy has been visited %d times!</p>\n</body>\n</html>", cfg.FileserverHits.Load())

	// Writes the request
	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte(body))
}
