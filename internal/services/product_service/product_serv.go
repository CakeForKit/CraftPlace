package productservice

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ProductServ interface {
	Add(ctx context.Context, addReq AddProductData) error
	Delete(ctx context.Context, productID uuid.UUID) error
	Update(ctx context.Context, productID uuid.UUID, updateReq UpdateProductData) error
}

var (
	ErrProductServ = errors.New("ProductServ")
)

type AddProductData struct {
	Title       string      `json:"title" binding:"required,max=255" example:"Звезды"`
	Description string      `json:"description" binding:"required,max=255" example:"Магазин сережек"`
	Cost        uint64      `json:"cost" binding:"required,min=0" example:"100"`
	ShopID      uuid.UUID   `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	CategoryIDs []uuid.UUID `json:"categoryIDs" binding:"required,dive,uuid"`
}

type UpdateProductData struct {
	Title       string      `json:"title" binding:"required,max=255" example:"Лучшие звезды"`
	Description string      `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
	Cost        uint64      `json:"cost" binding:"required,min=0" example:"200"`
	ShopID      uuid.UUID   `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	CategoryIDs []uuid.UUID `json:"categoryIDs" binding:"required,dive,uuid"`
}
