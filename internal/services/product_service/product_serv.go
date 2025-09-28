package productservice

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type ProductServ interface {
	GetProducts(ctx context.Context, filterOps *reqresp.ProductFilter) ([]*models.Product, error)
	Add(ctx context.Context, addReq reqresp.AddProductRequest) error
	Delete(ctx context.Context, productID uuid.UUID) error
	Update(ctx context.Context, updateReq reqresp.UpdateProductRequest) error
}

var (
	ErrProductServ = errors.New("ProductServ")
)
