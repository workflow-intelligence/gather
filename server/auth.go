package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/julienschmidt/httprouter"
	"github.com/workflow-intelligence/gather/user"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Error *ErrorMessage `json:"error"`
	Token string        `json:"token"`
}

func AuthLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var l LoginRequest
	err := decoder.Decode(&l)
	if err != nil {
		errorOut(w, 500, "ERR00010010", "Login error", err)
		return
	}
	u, err := user.Login(l.Username, l.Password, jwtSecretKey)
	if err != nil {
		errorOut(w, 500, "ERR00010011", "Login error", err)
		return
	}
	if u.Token == "" {
		errorOut(w, 403, "ERR00010012", "Invalid login", errors.New("Invalid login"))
		return
	}

	response := LoginResponse{Error: nil, Token: u.Token}

	log.Debug().Str("user", u.Name).Str("token", u.Token).Str("id", "DBG00010010").Msg("Login")
	json, err := json.Marshal(response)
	if err != nil {
		errorOut(w, 500, "ERR00010011", "Could marshal login info", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(json[:]))
}
