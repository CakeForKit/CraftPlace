package testobj

import (
	"time"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	"github.com/google/uuid"
)

type ProductMother interface {
	ProductP() *models.Product
}

func NewProductMother() ProductMother {
	return &productMother{}
}

type productMother struct{}

func (um *productMother) ProductP() *models.Product {
	product, _ := models.NewProduct(
		uuid.New(),
		"test-title"+uuid.New().String(),
		"test-desription",
		1000,
		uuid.New(),
		uuid.UUIDs{uuid.New(), uuid.New()},
		time.Date(2023, time.October, 1, 15, 30, 0, 0, time.UTC),
	)
	return product
}
