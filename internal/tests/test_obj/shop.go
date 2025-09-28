package testobj

import (
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
	)
	return shop
}
