package searcher

import (
	"context"
	"errors"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	testobj "github.com/CakeForKit/CraftPlace.git/internal/tests/test_obj"
	"github.com/google/uuid"
)

type Searcher interface {
	GetCategories(ctx context.Context, filterOps *reqresp.CategoryFilter) ([]*models.Category, error)
	GetShops(ctx context.Context, filterOps *reqresp.ShopFilter) ([]*models.Shop, error)
	GetPosts(ctx context.Context, filterOps *reqresp.PostFilter) ([]*models.Post, error)
	GetProducts(ctx context.Context, filterOps *reqresp.ProductFilter) ([]*models.Product, error)

	GetCategoruByID(ctx context.Context, categoryID uuid.UUID) (*models.Category, error)
	GetShopByID(ctx context.Context, shopID uuid.UUID) (*models.Shop, error)
}

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrShopNotFound     = errors.New("shop not found")
)

func NewSearcher() Searcher {
	return &searcher{}
}

type searcher struct {
}

func (s *searcher) GetCategories(ctx context.Context, filterOps *reqresp.CategoryFilter) ([]*models.Category, error) {
	categoryCreator := testobj.NewCategoryMother()
	res := make([]*models.Category, 3)
	for i := range res {
		res[i] = categoryCreator.CategoryP()
	}
	return res, nil
}

func (s *searcher) GetPosts(ctx context.Context, filterOps *reqresp.PostFilter) ([]*models.Post, error) {
	postCreator := testobj.NewPostMother()
	res := make([]*models.Post, 3)
	for i := range res {
		res[i] = postCreator.PostP()
	}
	return res, nil
}

func (s *searcher) GetProducts(ctx context.Context, filterOps *reqresp.ProductFilter) ([]*models.Product, error) {
	creator := testobj.NewProductMother()
	res := make([]*models.Product, 3)
	for i := range res {
		res[i] = creator.ProductP()
	}
	return res, nil
}

func (s *searcher) GetShops(ctx context.Context, filterOps *reqresp.ShopFilter) ([]*models.Shop, error) {
	creator := testobj.NewShopMother()
	res := make([]*models.Shop, 3)
	for i := range res {
		res[i] = creator.ShopP()
	}
	return res, nil
}

func (s *searcher) GetCategoruByID(ctx context.Context, categoryID uuid.UUID) (*models.Category, error) {
	categoryCreator := testobj.NewCategoryMother()
	return categoryCreator.CategoryP(), nil
}

func (s *searcher) GetShopByID(ctx context.Context, shopID uuid.UUID) (*models.Shop, error) {
	creator := testobj.NewShopMother()
	return creator.ShopP(), nil
}
