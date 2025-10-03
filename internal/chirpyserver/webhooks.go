package chirpyserver

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/roxensox/chirpy/internal/auth"
	"net/http"
)

func (cfg *ApiConfig) POSTPolkaWebhooks(writer http.ResponseWriter, req *http.Request) {
	// Handles POST request to polka/webhooks endpoint

	// Gets the API Key from the header
	apiKey, err := auth.GetAPIKey(req.Header)
	// Returns error code if API key isn't found or doesn't match config
	if err != nil || apiKey != cfg.APIKey {
		writer.WriteHeader(401)
		writer.Write([]byte("Invalid/Missing API Key"))
		return
	}

	// Creates object to receive webhook input
	rcv := struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}{}

	// Decodes input into object
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&rcv)

	// Returns success code if event isn't appropriate to prevent repeat requests
	if rcv.Event != "user.upgraded" {
		writer.WriteHeader(204)
		return
	}

	// Pulls UID from input and parses to UUID
	UID, err := uuid.Parse(rcv.Data.UserID)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Unable to parse user ID as UUID"))
		return
	}

	// Upgrades user chirpy red status
	err = cfg.DBConn.UpgradeUser(req.Context(), UID)
	if err != nil {
		writer.WriteHeader(404)
		writer.Write([]byte("User not found"))
		return
	}

	// Returns success code
	writer.WriteHeader(204)
}
