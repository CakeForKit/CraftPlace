package productservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	productrep "github.com/CakeForKit/CraftPlace.git/internal/repository/product_rep"
	shoprep "github.com/CakeForKit/CraftPlace.git/internal/repository/shop_rep"
	auth "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
	"github.com/google/uuid"
)

type ProductServ interface {
	Add(ctx context.Context, addReq models.AddProductData) (*models.Product, error)
	Delete(ctx context.Context, productID uuid.UUID, shopID uuid.UUID) error
	Update(ctx context.Context, productID uuid.UUID, updateReq models.UpdateProductData) (*models.Product, error)
}

var (
	ErrProductServ = errors.New("ProductServ")
	ErrWrongShop   = errors.New("product does not belong to the shop")
)

type productServ struct {
	productRep productrep.ProductRep
	authz      auth.AuthZ
	shopRep    shoprep.ShopRep
}

func NewProductServ(productRep productrep.ProductRep, authz auth.AuthZ, shopRep shoprep.ShopRep) ProductServ {
	return &productServ{
		productRep: productRep,
		authz:      authz,
		shopRep:    shopRep,
	}
}

func (s *productServ) checkUserRights(ctx context.Context, shopID uuid.UUID) error {
	baseErr := fmt.Errorf("checkUserRights %w", ErrWrongShop)
	userID, err := s.authz.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	shop, err := s.shopRep.GetByID(ctx, shopID)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	if shop.GetUserID() != userID {
		return fmt.Errorf("%w: %w", baseErr, auth.ErrHasNoRights)
	}
	return nil
}

func (s *productServ) Add(ctx context.Context, addReq models.AddProductData) (*models.Product, error) {
	baseErr := fmt.Errorf("%w Add", ErrProductServ)
	err := s.checkUserRights(ctx, addReq.ShopID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	newProduct, err := models.NewProduct(
		uuid.New(),
		addReq.Title,
		addReq.Description,
		addReq.Cost,
		addReq.ShopID,
		addReq.CategoryIDs,
		time.Now().UTC(),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	err = s.productRep.Add(ctx, newProduct)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return newProduct, nil
}

func (s *productServ) Delete(ctx context.Context, productID uuid.UUID, shopID uuid.UUID) error {
	baseErr := fmt.Errorf("%w Delete", ErrProductServ)
	product, err := s.productRep.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	err = s.checkUserRights(ctx, product.GetShopID())
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	if shopID != product.GetShopID() {
		return fmt.Errorf("%w:%w", baseErr, ErrWrongShop)
	}

	err = s.productRep.Delete(ctx, productID)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	return nil
}

func (s *productServ) Update(ctx context.Context, productID uuid.UUID, updateReq models.UpdateProductData) (*models.Product, error) {
	baseErr := fmt.Errorf("%w Update", ErrProductServ)
	product, err := s.productRep.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	err = s.checkUserRights(ctx, product.GetShopID())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	if updateReq.ShopID != product.GetShopID() {
		return nil, fmt.Errorf("%w:%w", baseErr, ErrWrongShop)
	}

	updated, err := s.productRep.Update(ctx, productID, func(*models.Product) (*models.Product, error) {
		err := product.Update(&updateReq)
		return product, err
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return updated, nil
}
