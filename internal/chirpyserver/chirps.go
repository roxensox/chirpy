package chirpyserver

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/roxensox/chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *ApiConfig) POSTChirps(writer http.ResponseWriter, req *http.Request) {
	// Handles a POST request to chirps endpoint, returns newly created chirp

	// Creates anonymous struct for receiving input
	inObj := struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}{}

	// Creates a JSON decoder for the request and decodes it into inObj
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&inObj)

	// Parses user_id as UUID
	user_id, err := uuid.Parse(inObj.UserID)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to parse user ID"))
		return
	}

	// Builds query param object
	params := database.CreateChirpParams{
		UserID:    user_id,
		Body:      inObj.Body,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ID:        uuid.New(),
	}

	// Queries the database to insert parameters
	dbResp, err := cfg.DBConn.CreateChirp(req.Context(), params)
	if err != nil {
		writer.WriteHeader(404)
		writer.Write([]byte("User not found"))
		return
	}

	// Casts response to output object
	outObj := Chirp{
		ID:        dbResp.ID,
		CreatedAt: dbResp.CreatedAt,
		UpdatedAt: dbResp.UpdatedAt,
		Body:      dbResp.Body,
		UserID:    dbResp.UserID,
	}

	// Marshals output object to JSON
	outJson, err := json.Marshal(outObj)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to marshal data"))
		return
	}

	// Writes success response
	writer.WriteHeader(201)
	writer.Write(outJson)
}

func (cfg *ApiConfig) GETChirps(writer http.ResponseWriter, req *http.Request) {
	// Handles a GET request to the chirps endpoint, returns all chirps

	// Initializes an empty slice of Chirp objects
	out := make([]Chirp, 0)

	// Queries the chirps from the database
	allChirps, err := cfg.DBConn.GetChirps(req.Context())
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Unable to get chirps"))
		return
	}

	// Iterates through the returned chirps
	for _, c := range allChirps {
		// Casts the db chirps to output object with appropriate JSON fields
		ctoJSON := Chirp{
			ID:        c.ID,
			CreatedAt: c.UpdatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
		// Adds the chirp to out slice
		out = append(out, ctoJSON)
	}

	// Marshals the out slice to JSON
	outJson, err := json.Marshal(out)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to marshal data"))
		return
	}

	// Writes http response
	writer.WriteHeader(200)
	writer.Write(outJson)
}

func (cfg *ApiConfig) GETChirpByID(writer http.ResponseWriter, req *http.Request) {
	// Handles a GET request to chirps/{chirpID} endpoint, returns chirp with matching ID

	// Extracts chirp ID from the path
	chirpID := req.PathValue("chirpID")

	// Parses the ID to a UUID
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to parse chirp ID"))
		return
	}

	// Queries the database for the matching chirp
	dbResp, err := cfg.DBConn.GetExactChirp(req.Context(), chirpUUID)
	if err != nil {
		writer.WriteHeader(404)
		writer.Write([]byte("Chirp not found"))
		return
	}

	// Casts the db response to a Chirp object for JSON marshaling
	out := Chirp{
		ID:        dbResp.ID,
		CreatedAt: dbResp.CreatedAt,
		UpdatedAt: dbResp.UpdatedAt,
		Body:      dbResp.Body,
		UserID:    dbResp.UserID,
	}

	// Marshals chirp to JSON
	outJson, err := json.Marshal(out)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to marshal chirp data"))
		return
	}

	// Writes success response
	writer.WriteHeader(200)
	writer.Write(outJson)
}
