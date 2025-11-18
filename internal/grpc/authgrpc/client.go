package authgrpc

import (
	"context"
	"fmt"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	tokenmaker "github.com/CakeForKit/CraftPlace.git/internal/services/auth/token_maker"
	"github.com/CakeForKit/CraftPlace.git/proto/pb/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client auth.AuthServiceClient
}

func NewClient(addr string) (*Client, error) {
	maxMsgSize := 1024 * 1024 * 100 // 100 MB
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: auth.NewAuthServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) LoginUser(ctx context.Context, lur models.LoginUserRequest) (string, error) {
	resp, err := c.client.LoginUser(ctx, &auth.LoginUserRequest{
		Login:    lur.Login,
		Password: lur.Password,
	})
	if err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}

func (c *Client) RegisterUser(ctx context.Context, rur models.RegisterUserRequest) error {
	_, err := c.client.RegisterUser(ctx, &auth.RegisterUserRequest{
		Login:    rur.Login,
		Password: rur.Password,
	})
	return err
}

func (c *Client) VerifyByToken(token string) (*tokenmaker.Payload, error) {
	fmt.Printf("VerifyByToken: \n\n")
	resp, err := c.client.VerifyToken(context.Background(), &auth.VerifyTokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	var role tokenmaker.RoleAuth
	if resp.Payload.Role == auth.Role_ROLE_USER {
		role = tokenmaker.UserRole
	} else {
		return nil, fmt.Errorf("wrong user role")
	}
	payload := &tokenmaker.Payload{
		PersonID:  uuid.MustParse(resp.Payload.PersonId),
		Role:      role,
		ExpiredAt: resp.Payload.ExpiredAt.AsTime(),
	}

	return payload, nil
}
