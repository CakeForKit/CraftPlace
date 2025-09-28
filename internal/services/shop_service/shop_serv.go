package shopservice

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type ShopServ interface {
	GetShops(ctx context.Context, filterOps *reqresp.ShopFilter) ([]*models.Shop, error)
	Add(ctx context.Context, addReq reqresp.AddShopRequest) error
	Delete(ctx context.Context, shopID uuid.UUID) error
	Update(ctx context.Context, updateReq reqresp.UpdateShopRequest) error
}

var (
	ErrShopServ = errors.New("ShopServ")
)
