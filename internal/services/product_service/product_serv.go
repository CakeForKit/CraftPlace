package productservice

import (
	"context"
	"errors"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type ProductServ interface {
	Add(ctx context.Context, addReq reqresp.AddProductRequest) error
	Delete(ctx context.Context, productID uuid.UUID) error
	Update(ctx context.Context, updateReq reqresp.UpdateProductRequest) error
}

var (
	ErrProductServ = errors.New("ProductServ")
)
