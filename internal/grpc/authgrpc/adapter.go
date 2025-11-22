package authgrpc

import (
	"context"
	"fmt"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	authuser "github.com/CakeForKit/CraftPlace.git/internal/services/auth/auth_user"
	"github.com/CakeForKit/CraftPlace.git/proto/pb/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	auth.UnimplementedAuthServiceServer
	authService authuser.AuthUser
}

func NewServer(authService authuser.AuthUser) *Server {
	return &Server{
		authService: authService,
	}
}

func (s *Server) LoginUser(ctx context.Context, req *auth.LoginUserRequest) (*auth.LoginUserResponse, error) {
	loginReq := models.LoginUserRequest{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}
	accessToken, err := s.authService.LoginUser(ctx, loginReq)
	if err != nil {
		// Конвертируем доменные ошибки в gRPC status codes
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &auth.LoginUserResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *Server) RegisterUser(ctx context.Context, req *auth.RegisterUserRequest) (*emptypb.Empty, error) {
	registerReq := models.RegisterUserRequest{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}
	err := s.authService.RegisterUser(ctx, registerReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest) (*auth.VerifyTokenResponse, error) {
	fmt.Print("Server VerifyToken\n")
	payload, err := s.authService.VerifyByToken(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	protoPayload := &auth.TokenPayload{
		PersonId:  payload.PersonID.String(),
		Role:      auth.Role(auth.Role_value["ROLE_USER"]),
		ExpiredAt: timestamppb.New(payload.ExpiredAt),
	}

	return &auth.VerifyTokenResponse{
		Payload: protoPayload,
	}, nil
}
