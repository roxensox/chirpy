package chirpyserver

import (
	"encoding/json"
	"github.com/roxensox/chirpy/internal/auth"
	"net/http"
)

func (cfg *ApiConfig) POSTLogin(writer http.ResponseWriter, req *http.Request) {
	inObj := struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&inObj)
	user, err := cfg.DBConn.GetUserByEmail(req.Context(), inObj.Email)
	validPass, err2 := auth.CheckPasswordHash(inObj.Password, user.HashedPassword)
	if err != nil || !validPass {
		writer.WriteHeader(401)
		writer.Write([]byte("Incorrect email or password"))
		return
	}
	if err2 != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Failed to compare passwords"))
		return
	}
	out := User{
		Email:     user.Email,
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	outJson, err := json.Marshal(out)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("Something went wrong"))
		return
	}
	writer.WriteHeader(200)
	writer.Write(outJson)
}
