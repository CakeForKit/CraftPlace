package models

import (
	"errors"
	"fmt"
	"strings"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

const (
	MaxLenUserLogin = 50
)

type User struct {
	id             uuid.UUID
	login          string // unique
	hashedPassword string
}

type LoginUserRequest struct {
	Login    string
	Password string
}

type RegisterUserRequest struct {
	Login    string
	Password string
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

func (u *User) ToResponse() reqresp.UserResponse {
	return reqresp.UserResponse{
		Login: u.GetLogin(),
	}
}

func (u *User) UpdateLogin(newLogin string) error {
	copyU := *u
	copyU.login = newLogin
	if err := copyU.validate(); err != nil {
		return err
	}
	*u = copyU
	return nil
}

func (u *User) UpdatePassword(newHashedPassword string) error {
	copyU := *u
	copyU.hashedPassword = newHashedPassword
	if err := copyU.validate(); err != nil {
		return err
	}
	*u = copyU
	return nil
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
