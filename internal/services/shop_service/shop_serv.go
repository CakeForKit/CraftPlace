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
	ErrShopServ  = errors.New("ShopServ")
	ErrWrongShop = errors.New("shop does not belong to the user")
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

func (s *shopServ) checkUserRights(ctx context.Context, shopID uuid.UUID) error {
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
	if err := s.checkUserRights(ctx, shopID); err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	if err := s.shopRep.Delete(ctx, shopID); err != nil {
		return fmt.Errorf("%w: %w", baseErr, err)
	}
	return nil
}

func (s *shopServ) Update(ctx context.Context, shopID uuid.UUID, updateReq reqresp.UpdateShopRequest) (*models.Shop, error) {
	baseErr := fmt.Errorf("%w Update", ErrShopServ)
	if err := s.checkUserRights(ctx, shopID); err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	updated, err := s.shopRep.Update(ctx, shopID, func(sh *models.Shop) (*models.Shop, error) {
		err := sh.Update(&updateReq)
		return sh, err
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return updated, nil
}
