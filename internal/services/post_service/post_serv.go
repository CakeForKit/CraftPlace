package postservice

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type PostServ interface {
	GetPosts(ctx context.Context) ([]*models.Post, error)
	Add(ctx context.Context, addReq reqresp.AddPostRequest) error
	Delete(ctx context.Context, postID uuid.UUID) error
	Update(ctx context.Context, updateReq reqresp.UpdatePostRequest) error
}

var (
	ErrPostServ = errors.New("PostServ")
)
