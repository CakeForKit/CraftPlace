package productrep

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type ProductRep interface {
	GetByFilter(ctx context.Context, filterOps *reqresp.ProductFilter) ([]*models.Product, error)
	GetByID(ctx context.Context, productID uuid.UUID) (*models.Product, error)
	Add(ctx context.Context, m *models.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, productID uuid.UUID, funcUpdate func(*models.Product) (*models.Product, error)) (*models.Product, error)
	Ping(ctx context.Context) error
	Close()
}

var (
	ErrProductNotFound = errors.New("the Product was not found")
	ErrFailedToAddUser = errors.New("failed to add the Product")
	ErrUpdateProduct   = errors.New("failed to update the Product in the repository")
)
