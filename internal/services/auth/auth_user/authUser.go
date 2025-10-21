package authuser

import (
	"context"

	"github.com/CakeForKit/CraftPlace.git/internal/cnfg"
	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	userrep "github.com/CakeForKit/CraftPlace.git/internal/repository/user_rep"
	"github.com/CakeForKit/CraftPlace.git/internal/services/auth/hasher"
	tokenmaker "github.com/CakeForKit/CraftPlace.git/internal/services/auth/token_maker"
	"github.com/google/uuid"
)

type AuthUser interface {
	LoginUser(ctx context.Context, lur models.LoginUserRequest) (string, error)
	RegisterUser(ctx context.Context, rur models.RegisterUserRequest) error
	VerifyByToken(token string) (*tokenmaker.Payload, error)
}

type authUser struct {
	tokenMaker tokenmaker.TokenMaker
	hasher     hasher.Hasher
	config     *cnfg.AppConfig
	userrep    userrep.UserRep
}

func NewAuthUser(tokenMaker tokenmaker.TokenMaker, hasher hasher.Hasher, config *cnfg.AppConfig, urep userrep.UserRep) AuthUser {
	server := &authUser{
		tokenMaker: tokenMaker,
		hasher:     hasher,
		config:     config,
		userrep:    urep,
	}
	return server
}

func (s *authUser) LoginUser(ctx context.Context, lur models.LoginUserRequest) (string, error) {
	user, err := s.userrep.GetByLogin(ctx, lur.Login)
	if err != nil {
		return "", err
	}

	err = s.hasher.CheckPassword(lur.Password, user.GetHashedPassword())
	if err != nil {
		return "", err
	}
	userID := user.GetID()
	accessToken, err := s.tokenMaker.CreateToken(
		userID,
		tokenmaker.UserRole,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *authUser) RegisterUser(ctx context.Context, rur models.RegisterUserRequest) error {
	hashedPassword, err := s.hasher.HashPassword(rur.Password)
	if err != nil {
		return err
	}
	user, err := models.NewUser(
		uuid.New(),
		rur.Login,
		hashedPassword,
	)
	if err != nil {
		return nil
	}
	err = s.userrep.Add(ctx, &user)
	return err
}

func (s *authUser) VerifyByToken(tokenStr string) (*tokenmaker.Payload, error) {
	return s.tokenMaker.VerifyToken(tokenStr, tokenmaker.UserRole)
}
