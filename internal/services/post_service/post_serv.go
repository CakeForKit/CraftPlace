package postservice

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type PostServ interface {
	GetPostsByFilter(ctx context.Context, filterOps *reqresp.PostFilter) ([]*models.Post, error)
	Add(ctx context.Context, addReq AddPostData) error
	Delete(ctx context.Context, postID uuid.UUID) error
	// Update(ctx context.Context, updateReq reqresp.UpdatePostRequest) error
}

var (
	ErrPostServ = errors.New("PostServ")
)

type AddPostData struct {
	Description string    `json:"description" binding:"required,max=255" example:"Магазин сережек"`
	ShopID      uuid.UUID `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
}
