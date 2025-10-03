package shopservice

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type ShopServ interface {
	// UserID в контексте
	Add(ctx context.Context, addReq reqresp.AddShopRequest) (*models.Shop, error)
	Delete(ctx context.Context, shopID uuid.UUID) error
	Update(ctx context.Context, updateReq reqresp.UpdateShopRequest) (*models.Shop, error)
}

var (
	ErrShopServ     = errors.New("ShopServ")
	ErrShopNotFound = errors.New("shop not found")
)
