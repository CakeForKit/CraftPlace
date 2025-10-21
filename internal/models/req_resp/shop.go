package reqresp

import (
	"time"

	"github.com/google/uuid"
)

type AddShopRequest struct {
	Title       string `json:"title" binding:"required,max=255" example:"Звезды"`
	Description string `json:"description" binding:"required,max=255" example:"Магазин сережек"`
}

type UpdateShopRequest struct {
	Title       string `json:"title" binding:"required,max=255" example:"Лучшие звезды"`
	Description string `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
}

type ShopResponse struct {
	ShopID      string    `json:"id_shop" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	Title       string    `json:"title" example:"Eco"`
	Description string    `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
	UserID      uuid.UUID `json:"userID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	UpdateTime  time.Time `json:"updateTime" example:"2023-06-15T14:30:00Z"`
}

type ShopsResponse struct {
	Shops []ShopResponse `json:"shops"`
}
