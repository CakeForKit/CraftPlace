package auth

import (
	"context"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	tokenmaker "github.com/CakeForKit/CraftPlace.git/internal/services/auth/token_maker"
)

type AuthUser interface {
	LoginUser(ctx context.Context, lur reqresp.LoginUserRequest) (string, error)
	RegisterUser(ctx context.Context, rur reqresp.RegisterUserRequest) error
	VerifyByToken(token string) (*tokenmaker.Payload, error)
}

/*
type authUser struct {
	tokenMaker tokenmaker.TokenMaker
	config     cnfg.AppConfig
	userrep    userrep.UserRep
	hasher     hasher.Hasher
}

func NewAuthUser(config cnfg.AppConfig, urep userrep.UserRep, tokenMaker tokenmaker.TokenMaker, hasher hasher.Hasher) (AuthUser, error) {
	server := &authUser{
		tokenMaker: tokenMaker,
		config:     config,
		userrep:    urep,
		hasher:     hasher,
	}

	return server, nil
}

func (s *authUser) LoginUser(ctx context.Context, lur LoginUserRequest) (string, error) {
	user, err := s.userrep.GetByLogin(ctx, lur.Login)
	if err != nil {
		return "", err
	}

	err = s.hasher.CheckPassword(lur.Password, user.GetHashedPassword())
	if err != nil {
		return "", err
	}

	accessToken, err := s.tokenMaker.CreateToken(
		user.GetID(),
		tokenmaker.UserRole,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *authUser) RegisterUser(ctx context.Context, rur RegisterUserRequest) error {
	hashedPassword, err := s.hasher.HashPassword(rur.Password)
	if err != nil {
		return err
	}
	user, err := models.NewUser(
		uuid.New(),
		rur.Username,
		rur.Login,
		hashedPassword,
		time.Now(),
		rur.Email,
		rur.SubscribeEmail,
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
*/
