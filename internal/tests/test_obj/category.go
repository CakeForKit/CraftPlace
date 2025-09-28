package testobj

import (
	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	"github.com/google/uuid"
)

type CategoryMother interface {
	CategoryP() *models.Category
}

func NewCategoryMother() CategoryMother {
	return &categoryMother{}
}

type categoryMother struct{}

func (um *categoryMother) CategoryP() *models.Category {
	category, _ := models.NewCategory(
		uuid.New(),
		"test-title"+uuid.New().String(),
		"test-desription",
	)
	return category
}
