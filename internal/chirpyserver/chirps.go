package chirpyserver

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/roxensox/chirpy/internal/auth"
	"github.com/roxensox/chirpy/internal/database"
	"log"
	"net/http"
	"sort"
	"time"
)

func (cfg *ApiConfig) POSTChirps(writer http.ResponseWriter, req *http.Request) {
	// Handles a POST request to chirps endpoint, returns newly created chirp

	// Creates anonymous struct for receiving input
	inObj := struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}{}

	tkn, err := auth.GetBearerToken(req.Header)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte("Must be logged in"))
		return
	}

	UID, err := auth.ValidateJWT(tkn, cfg.Secret)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte("Unauthorized"))
		log.Println(err)
		log.Println(tkn)
		return
	}

	// Creates a JSON decoder for the request and decodes it into inObj
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&inObj)

	chirpID, err := uuid.NewUUID()
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to generate post ID"))
		return
	}

	// Builds query param object
	params := database.CreateChirpParams{
		UserID:    UID,
		Body:      inObj.Body,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ID:        chirpID,
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

	// Sets content type in header
	writer.Header().Set("Content-Type", "application/json")

	// Initializes an empty slice of Chirp objects
	out := make([]Chirp, 0)

	// Gets the author ID as a string from the query, if there is one
	auth_id := req.URL.Query().Get("author_id")

	// Declares variables used in conditional scope
	var allChirps []database.Chirp
	var err error
	var auth_uuid uuid.UUID

	// Proceeds as usual if no ID is provided
	if auth_id != "" {
		// Parses author ID to UUID
		auth_uuid, err = uuid.Parse(auth_id)
		if err != nil {
			http.Error(writer, "Invalid author ID", http.StatusBadRequest)
			return
		}

		allChirps, err = cfg.DBConn.GetChirpsByAuthor(req.Context(), auth_uuid)
		if err != nil {
			http.Error(writer, "Unable to get chirps", http.StatusInternalServerError)
			return
		}
	} else {
		// Queries the DB with nullable auth_uuid
		allChirps, err = cfg.DBConn.GetChirps(req.Context())
		if err != nil {
			http.Error(writer, "Unable to get chirps", http.StatusInternalServerError)
			return
		}
	}

	if req.URL.Query().Get("sort") == "desc" {
		sort.Slice(allChirps, func(i, j int) bool {
			return allChirps[i].CreatedAt.After(allChirps[j].CreatedAt)
		})
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
		http.Error(writer, "Failed to marshal data", http.StatusInternalServerError)
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

func (cfg *ApiConfig) DELETEChirpByID(writer http.ResponseWriter, req *http.Request) {
	// Handles DELETE requests at the chirps/{chirpID} endpoint, deleting a chirp by ID

	// Gets the chirp ID as a string from the endpoint
	chirpID := req.PathValue("chirpID")

	// Reads the ID into a UUID
	CID, err := uuid.Parse(chirpID)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Unable to parse chirp ID"))
		return
	}

	// Queries the chirp from the DB
	chirp, err := cfg.DBConn.GetExactChirp(req.Context(), CID)
	if err != nil {
		writer.WriteHeader(404)
		writer.Write([]byte("Chirp not found"))
		return
	}

	// Gets the access token from the request header
	tkn, err := auth.GetBearerToken(req.Header)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte("Token not found"))
		return
	}

	// Gets the user's ID by validating the token
	UID, err := auth.ValidateJWT(tkn, cfg.Secret)
	if err != nil {
		writer.WriteHeader(401)
		writer.Write([]byte("Invalid token"))
		return
	}

	// Compares the chirp's user ID with the token's
	if chirp.UserID != UID {
		writer.WriteHeader(403)
		writer.Write([]byte("Unauthorized"))
		return
	}

	// Deletes the chirp
	err = cfg.DBConn.DeleteChirp(req.Context(), CID)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to delete chirp"))
		return
	}

	// Writes success response
	writer.WriteHeader(204)
}
