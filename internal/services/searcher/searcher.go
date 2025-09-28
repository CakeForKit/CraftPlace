package searcher

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	"github.com/google/uuid"
)

type Searcher interface {
	GetCategories(ctx context.Context, filterOps *reqresp.CategoryFilter) ([]*models.Category, error)
	GetPosts(ctx context.Context) ([]*models.Post, error)
	GetProducts(ctx context.Context, filterOps *reqresp.ProductFilter) ([]*models.Product, error)
	GetShops(ctx context.Context, filterOps *reqresp.ShopFilter) ([]*models.Shop, error)
	GetCategoruByID(ctx context.Context, categoryID uuid.UUID) (models.Category, error)
}

var (
	ErrCategoryNotFpund = errors.New("category not found")
)
