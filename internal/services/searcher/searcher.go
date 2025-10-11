package searcher

import (
	"context"
	"errors"
	"fmt"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	categoryrep "github.com/CakeForKit/CraftPlace.git/internal/repository/category_rep"
	postrep "github.com/CakeForKit/CraftPlace.git/internal/repository/post_rep"
	productrep "github.com/CakeForKit/CraftPlace.git/internal/repository/product_rep"
	shoprep "github.com/CakeForKit/CraftPlace.git/internal/repository/shop_rep"
	"github.com/google/uuid"
)

type Searcher interface {
	GetCategories(ctx context.Context, filterOps *reqresp.CategoryFilter) ([]*models.Category, error)
	GetShops(ctx context.Context, filterOps *reqresp.ShopFilter) ([]*models.Shop, error)
	GetPosts(ctx context.Context, filterOps *reqresp.PostFilter) ([]*models.Post, error)
	GetProducts(ctx context.Context, filterOps *reqresp.ProductFilter) ([]*models.Product, error)

	GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*models.Category, error)
	GetShopByID(ctx context.Context, shopID uuid.UUID) (*models.Shop, error)
}

var (
	ErrSearcher = errors.New("searcher")
)

func NewSearcher(categoryRep categoryrep.CategoryRep,
	shopRep shoprep.ShopRep,
	postRep postrep.PostRep,
	productRep productrep.ProductRep,
) Searcher {
	return &searcher{
		shopRep:     shopRep,
		categoryRep: categoryRep,
		postRep:     postRep,
		productRep:  productRep,
	}
}

type searcher struct {
	categoryRep categoryrep.CategoryRep
	shopRep     shoprep.ShopRep
	postRep     postrep.PostRep
	productRep  productrep.ProductRep
}

func (s *searcher) GetCategories(ctx context.Context, filterOps *reqresp.CategoryFilter) ([]*models.Category, error) {
	baseErr := fmt.Errorf("%w, GetCategories", ErrSearcher)
	res, err := s.categoryRep.GetByFilter(ctx, filterOps)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return res, nil
}

func (s *searcher) GetPosts(ctx context.Context, filterOps *reqresp.PostFilter) ([]*models.Post, error) {
	baseErr := fmt.Errorf("%w, GetPosts", ErrSearcher)
	res, err := s.postRep.GetByFilter(ctx, filterOps)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return res, nil
}

func (s *searcher) GetProducts(ctx context.Context, filterOps *reqresp.ProductFilter) ([]*models.Product, error) {
	baseErr := fmt.Errorf("%w, GetProducts", ErrSearcher)
	res, err := s.productRep.GetByFilter(ctx, filterOps)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return res, nil
}

func (s *searcher) GetShops(ctx context.Context, filterOps *reqresp.ShopFilter) ([]*models.Shop, error) {
	baseErr := fmt.Errorf("%w, GetShops", ErrSearcher)
	res, err := s.shopRep.GetByFilter(ctx, filterOps)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return res, nil
}

func (s *searcher) GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*models.Category, error) {
	baseErr := fmt.Errorf("%w, GetCategoryByID", ErrSearcher)
	res, err := s.categoryRep.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return res, nil
}

func (s *searcher) GetShopByID(ctx context.Context, shopID uuid.UUID) (*models.Shop, error) {
	baseErr := fmt.Errorf("%w, GetShopByID", ErrSearcher)
	res, err := s.shopRep.GetByID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return res, nil
}
