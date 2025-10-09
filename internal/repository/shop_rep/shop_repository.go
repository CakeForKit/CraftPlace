package shoprep

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type ShopRep interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Shop, error)
	GetByFilter(ctx context.Context, filterOps *reqresp.ShopFilter) ([]*models.Shop, error)

	Add(ctx context.Context, m *models.Shop) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, funcUpdate func(*models.Shop) (*models.Shop, error)) error
	Ping(ctx context.Context) error
	Close()
}

var (
	ErrShopNotFound    = errors.New("the Shop was not found")
	ErrFailedToAddUser = errors.New("failed to add the Shop")
	ErrUpdateShop      = errors.New("failed to update the Shop in the repository")
)
