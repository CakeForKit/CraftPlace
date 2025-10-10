package shopservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	shoprep "github.com/CakeForKit/CraftPlace.git/internal/repository/shop_rep"
	auth "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
	"github.com/google/uuid"
)

type ShopServ interface {
	// UserID в контексте
	Add(ctx context.Context, addReq reqresp.AddShopRequest) (*models.Shop, error)
	Delete(ctx context.Context, shopID uuid.UUID) error
	Update(ctx context.Context, shopID uuid.UUID, updateReq reqresp.UpdateShopRequest) (*models.Shop, error)
}

var (
	ErrShopServ     = errors.New("ShopServ")
	ErrShopNotFound = errors.New("shop not found")
)

type shopServ struct {
	shopRep shoprep.ShopRep
	authz   auth.AuthZ
}

func NewShopServ(shopRep shoprep.ShopRep, authz auth.AuthZ) ShopServ {
	return &shopServ{
		shopRep: shopRep,
		authz:   authz,
	}
}

func (s *shopServ) Add(ctx context.Context, addReq reqresp.AddShopRequest) (*models.Shop, error) {
	baseErr := fmt.Errorf("%w Add", ErrShopServ)
	userID, err := s.authz.UserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	newShop, err := models.NewShop(
		uuid.New(),
		addReq.Title,
		addReq.Description,
		userID,
		time.Now().UTC(),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	err = s.shopRep.Add(ctx, newShop)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return newShop, nil
}

func (s *shopServ) Delete(ctx context.Context, shopID uuid.UUID) error {
	baseErr := fmt.Errorf("%w Delete", ErrShopServ)
	userID, err := s.authz.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	delShop, err := s.shopRep.GetByID(ctx, shopID)
	if err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	if delShop.GetUserID() != userID {
		return fmt.Errorf("%w: %w", baseErr, auth.ErrHasNoRights)
	}

	if err = s.shopRep.Delete(ctx, shopID); err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	return nil
}

func (s *shopServ) Update(ctx context.Context, shopID uuid.UUID, updateReq reqresp.UpdateShopRequest) (*models.Shop, error) {
	baseErr := fmt.Errorf("%w Update", ErrShopServ)
	userID, err := s.authz.UserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	shop, err := s.shopRep.GetByID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	if shop.GetUserID() != userID {
		return nil, fmt.Errorf("%w: %w", baseErr, auth.ErrHasNoRights)
	}

	updated, err := s.shopRep.Update(ctx, shopID, func(shop *models.Shop) (*models.Shop, error) {
		err := shop.Update(&updateReq)
		return shop, err
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return updated, nil
}
