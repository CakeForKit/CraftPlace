package userselfservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	auth "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
	testobj "github.com/CakeForKit/CraftPlace.git/internal/tests/test_obj"
	"github.com/google/uuid"
)

type UserSelfServ interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	ChangeLogin(ctx context.Context, newLogin string) error
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

func (s *userSelfServ) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	creator := testobj.NewUserMother()
	return creator.DefaultUserP(userID), nil
}

func (s *userSelfServ) ChangeLogin(ctx context.Context, newLogin string) error {
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
