package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"net/http"
)

type ErrorMessage struct {
	Message string `json:"message"`
	Err     error  `json:"error"`
}

type ErrorResponse struct {
	Error *ErrorMessage `json:"error"`
}

func errorOut(w http.ResponseWriter, code int, id string, message string, err error) {
	var response ErrorResponse
	response.Error = &ErrorMessage{message, err}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	log.Error().Err(err).Str("id", id).Msg(message)
	json, _ := json.Marshal(response)
	fmt.Fprintf(w, string(json[:]))
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		log.Error().Err(err).Str("id", "ERR00010004").Msg("Could not parse the token")
		return err
	}
	if !token.Valid {
		return errors.New("Invalid token")
	}

	return nil
}

func JWTAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			errorOut(w, 401, "ERR00010401", "Unauthorized", errors.New("Missing Authorization header"))
			return
		}
		tokenString = tokenString[len("Bearer "):]
		err := verifyToken(tokenString)
		if err != nil {
			errorOut(w, 401, "ERR00010401", "Unauthorized", err)
			return
		}
		h(w, r, ps)
	}
}

func BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}
