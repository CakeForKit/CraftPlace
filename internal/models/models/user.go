package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	MaxLenUsername  = 50
	MaxLenUserLogin = 50
)

type User struct {
	id             uuid.UUID
	username       string
	login          string // unique
	hashedPassword string
}

var (
	ErrUserValidate = errors.New("model user validate error")
)

func NewUser(id uuid.UUID, username string, login string, hashedPassword string) (User, error) {
	user := User{
		id:             id,
		username:       strings.TrimSpace(username),
		login:          strings.TrimSpace(login),
		hashedPassword: hashedPassword,
	}
	err := user.validate()
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *User) validate() error {
	if u.username == "" || len(u.username) > MaxLenUsername {
		return fmt.Errorf("%w username", ErrUserValidate)
	} else if u.login == "" || len(u.login) > MaxLenUserLogin {
		return fmt.Errorf("%w username", ErrUserValidate)
	} else if u.hashedPassword == "" {
		return fmt.Errorf("%w hashedPassword", ErrUserValidate)
	}
	return nil
}

func (u *User) GetID() uuid.UUID {
	return u.id
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) GetLogin() string {
	return u.login
}

func (u *User) GetHashedPassword() string {
	return u.hashedPassword
}
