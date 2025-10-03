package models

import (
	"errors"
	"fmt"
	"strings"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

const (
	MaxLenUsername  = 50
	MaxLenUserLogin = 50
)

type User struct {
	id             uuid.UUID
	login          string // unique
	hashedPassword string
}

var (
	ErrUserValidate = errors.New("model user validate error")
)

func NewUser(id uuid.UUID, login string, hashedPassword string) (User, error) {
	user := User{
		id:             id,
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
	if u.login == "" || len(u.login) > MaxLenUserLogin {
		return fmt.Errorf("%w login", ErrUserValidate)
	} else if u.hashedPassword == "" {
		return fmt.Errorf("%w hashedPassword", ErrUserValidate)
	}
	return nil
}

func (p *User) ToResponse() reqresp.UserResponse {
	return reqresp.UserResponse{
		Login: p.GetLogin(),
	}
}

func (u *User) GetID() uuid.UUID {
	return u.id
}

func (u *User) GetLogin() string {
	return u.login
}

func (u *User) GetHashedPassword() string {
	return u.hashedPassword
}
