package testobj

import (
	"time"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	"github.com/google/uuid"
)

type ShopMother interface {
	ShopP() *models.Shop
}

func NewShopMother() ShopMother {
	return &shopMother{}
}

type shopMother struct{}

func (um *shopMother) ShopP() *models.Shop {
	shop, _ := models.NewShop(
		uuid.New(),
		"test-title"+uuid.New().String(),
		"test-desription",
		uuid.New(),
		time.Date(2023, time.October, 1, 15, 30, 0, 0, time.UTC),
	)
	return shop
}
