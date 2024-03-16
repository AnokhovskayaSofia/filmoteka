package db

import (
	"errors"
	"log/slog"

	"github.com/go-pg/pg/v10"
)

const (
	Admin  = "admin"
	Client = "client"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"date" validate:"oneof admin client"`
}

func GetUser(db *pg.DB, username string, password string) (string, error) {
	user := &User{}

	err := db.Model(user).Where("username = ?", username).
		Select()

	if err != nil {
		slog.Error("no user with such username")
		err = errors.New("no user with such username")
		return "", err
	}

	if user.Password != password {
		slog.Error("wrong password for use")
		err = errors.New("wrong password for use")
		return "", err
	}

	switch user.Role {
	case Admin:
		return Admin, err
	case Client:
		return Client, err
	}
	return "", err
}
