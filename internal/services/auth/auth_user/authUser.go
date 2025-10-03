package authuser

import (
	"context"
	"errors"
	"time"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/CakeForKit/CraftPlace.git/internal/services/auth/hasher"
	tokenmaker "github.com/CakeForKit/CraftPlace.git/internal/services/auth/token_maker"
	"github.com/google/uuid"
)

type AuthUser interface {
	LoginUser(ctx context.Context, lur reqresp.LoginUserRequest) (string, error)
	RegisterUser(ctx context.Context, rur reqresp.RegisterUserRequest) error
	VerifyByToken(token string) (*tokenmaker.Payload, error)
}

var (
	ErrDuplicateLoginUser = errors.New("err Duplicate Login User")
	ErrUserNotFound       = errors.New("err User Not Found")
)

type authUser struct {
	tokenMaker tokenmaker.TokenMaker
	hasher     hasher.Hasher
	// config     cnfg.AppConfig
	// userrep    userrep.UserRep

}

// config cnfg.AppConfig, urep userrep.UserRep,
func NewAuthUser(tokenMaker tokenmaker.TokenMaker, hasher hasher.Hasher) AuthUser {
	server := &authUser{
		tokenMaker: tokenMaker,
		hasher:     hasher,
		// config:     config,
		// userrep:    urep,
	}
	return server
}

func (s *authUser) LoginUser(ctx context.Context, lur reqresp.LoginUserRequest) (string, error) {
	// user, err := s.userrep.GetByLogin(ctx, lur.Login)
	// if err != nil {
	// 	return "", err
	// }

	// err = s.hasher.CheckPassword(lur.Password, user.GetHashedPassword())
	// if err != nil {
	// 	return "", err
	// }
	userID := uuid.New()
	s_config_AccessTokenDuration := time.Hour
	accessToken, err := s.tokenMaker.CreateToken(
		userID,
		tokenmaker.UserRole,
		s_config_AccessTokenDuration,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *authUser) RegisterUser(ctx context.Context, rur reqresp.RegisterUserRequest) error {
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
	_ = user
	// err = s.userrep.Add(ctx, &user)
	return nil
}

func (s *authUser) VerifyByToken(tokenStr string) (*tokenmaker.Payload, error) {
	return s.tokenMaker.VerifyToken(tokenStr, tokenmaker.UserRole)
}
