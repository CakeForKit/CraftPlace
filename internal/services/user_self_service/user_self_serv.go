package userselfservice

import (
	"context"
	"errors"
	"fmt"

	auth "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
)

type UserSelfServ interface {
	ChangeName(ctx context.Context, newName string) error
	ChangePassword(ctx context.Context, newPassword string) error
}

var (
	ErrUserSelfServ = errors.New("UserSelfServ")
)

func NewUserSelfServ(authz auth.AuthZ) UserSelfServ {
	return &userSelfServ{
		authz: authz,
	}
}

type userSelfServ struct {
	authz auth.AuthZ
}

func (s *userSelfServ) ChangeName(ctx context.Context, newName string) error {
	userID, err := s.authz.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUserSelfServ, err)
	}
	_ = userID
	// TODO UserRep.Update
	return nil
}
func (s *userSelfServ) ChangePassword(ctx context.Context, newPassword string) error {
	userID, err := s.authz.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUserSelfServ, err)
	}
	_ = userID
	// TODO UserRep.Update
	return nil
}
