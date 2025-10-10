package userselfservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	userrep "github.com/CakeForKit/CraftPlace.git/internal/repository/user_rep"
	auth "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
	"github.com/CakeForKit/CraftPlace.git/internal/services/auth/hasher"
	"github.com/google/uuid"
)

type UserSelfServ interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	ChangeLogin(ctx context.Context, userID uuid.UUID, newLogin string) error
	ChangePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
}

var (
	ErrUserSelfServ = errors.New("UserSelfServ")
	ErrUserIDInCtx  = errors.New("userID != userID in context")
)

func NewUserSelfServ(authz auth.AuthZ, userRep userrep.UserRep, hasher hasher.Hasher) UserSelfServ {
	return &userSelfServ{
		authz:   authz,
		userRep: userRep,
		hasher:  hasher,
	}
}

type userSelfServ struct {
	authz   auth.AuthZ
	userRep userrep.UserRep
	hasher  hasher.Hasher
}

func (s *userSelfServ) checkUser(ctx context.Context, userID uuid.UUID) error {
	baseErr := fmt.Errorf("checkUser")
	ctxUserID, err := s.authz.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	if ctxUserID != userID {
		return fmt.Errorf("%w: %w", baseErr, ErrUserIDInCtx)
	}
	return nil
}

func (s *userSelfServ) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	baseErr := fmt.Errorf("%w GetUserByID", ErrUserSelfServ)
	if err := s.checkUser(ctx, userID); err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	user, err := s.userRep.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return user, nil
}

func (s *userSelfServ) ChangeLogin(ctx context.Context, userID uuid.UUID, newLogin string) error {
	baseErr := fmt.Errorf("%w ChangeLogin", ErrUserSelfServ)
	if err := s.checkUser(ctx, userID); err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	_, err := s.userRep.Update(ctx, userID, func(u *models.User) (*models.User, error) {
		err := u.UpdateLogin(newLogin)
		return u, err
	})
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	return nil
}
func (s *userSelfServ) ChangePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	baseErr := fmt.Errorf("%w ChangePassword", ErrUserSelfServ)
	if err := s.checkUser(ctx, userID); err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	newHashedPassword, err := s.hasher.HashPassword(newPassword)
	if err != nil {
		return err
	}
	_, err = s.userRep.Update(ctx, userID, func(u *models.User) (*models.User, error) {
		err := u.UpdatePassword(newHashedPassword)
		return u, err
	})
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	return nil
}
