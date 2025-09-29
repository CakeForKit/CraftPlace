package reqresp

import "github.com/google/uuid"

type AddProductRequest struct {
	Title       string      `json:"title" binding:"required,max=255" example:"Звезды"`
	Description string      `json:"description" binding:"required,max=255" example:"Магазин сережек"`
	Cost        uint64      `json:"cost" binding:"required,min=0" example:"100"`
	ShopID      uuid.UUID   `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	CategoryIDs []uuid.UUID `json:"categoryIDs" binding:"required,dive,uuid"`
}

type UpdateProductRequest struct {
	ID          string      `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	Title       string      `json:"title" binding:"required,max=255" example:"Лучшие звезды"`
	Description string      `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
	Cost        uint64      `json:"cost" binding:"required,min=0" example:"200"`
	ShopID      uuid.UUID   `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	CategoryIDs []uuid.UUID `json:"categoryIDs" binding:"required,dive,uuid"`
}

type ProductResponse struct {
	ID          string      `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	Title       string      `json:"title" example:"Eco"`
	Description string      `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
	Cost        uint64      `json:"cost" binding:"required,min=0" example:"200"`
	ShopID      uuid.UUID   `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	CategoryIDs []uuid.UUID `json:"categoryIDs" binding:"required,dive,uuid"`
}

type DeleteProductRequest struct {
	ID string `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
}
