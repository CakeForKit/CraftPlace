package reqresp

import "github.com/google/uuid"

type AddShopRequest struct {
	Title       string    `json:"title" binding:"required,max=255" example:"Звезды"`
	Description string    `json:"description" binding:"required,max=255" example:"Магазин сережек"`
	UserID      uuid.UUID `json:"userID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
}

type UpdateShopRequest struct {
	Title       string    `json:"title" binding:"required,max=255" example:"Лучшие звезды"`
	Description string    `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
	UserID      uuid.UUID `json:"userID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
}

type ShopResponse struct {
	ID          string    `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	Title       string    `json:"title" example:"Eco"`
	Description string    `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
	UserID      uuid.UUID `json:"userID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
}
