package categoryservice

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type CategoryServ interface {
	GetCategorys(ctx context.Context, filterOps *reqresp.CategoryFilter) ([]*models.Category, error)
	Add(ctx context.Context, addReq reqresp.AddCategoryRequest) error
	Delete(ctx context.Context, categoryID uuid.UUID) error
	Update(ctx context.Context, updateReq reqresp.UpdateCategoryRequest) error
}

var (
	ErrCategoryServ = errors.New("CategoryServ")
)
