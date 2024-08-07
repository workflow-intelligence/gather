package user

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type User struct {
	Name         string
	PasswordHash string `yaml:"PasswordHash"`
	EMail        string `yaml:"EMail"`
	Token        string
}

type UserList map[string]User

func Login(username string, password string, jwtsecret string) (*User, error) {
	if username == "" {
		return nil, errors.New("No username provided")
	}
	if password == "" {
		return nil, errors.New("No password provided")
	}

	var userlist UserList
	f, err := os.ReadFile("users.yaml")
	if err != nil {
		log.Error().Str("file", "users.yaml").Err(err).Str("id", "ERR00020001").Msg("Could not read user database")
		return nil, err
	}

	err = yaml.Unmarshal(f, &userlist)
	if err != nil {
		log.Error().Err(err).Str("id", "ERR00020002").Str("file", "users.yaml").Msg("Could not unmarshall yaml")
		return nil, err
	}

	Login := userlist.Find(username)
	if Login == nil {
		log.Error().Err(err).Str("id", "ERR00020003").Str("file", "users.yaml").Str("user", username).Msg("User not found")
		return Login, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(Login.PasswordHash), []byte(password))
	if err != nil {
		log.Error().Err(err).Str("id", "ERR00020004").Str("user", username).Msg("Password mismatch")
		return Login, nil
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": Login.Name,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	token, err := t.SignedString([]byte(jwtsecret))
	if err != nil {
		log.Error().Err(err).Str("id", "ERR00020005").Str("user", username).Msg("Could not get signed token")
		return Login, err
	}
	Login.Token = token

	return Login, nil
}

func (List UserList) Find(Name string) *User {
	for n, u := range List {
		if n == Name {
			entry := new(User)
			entry.Name = Name
			entry.PasswordHash = u.PasswordHash
			entry.EMail = u.EMail
			return entry
		}
	}
	return nil
}
