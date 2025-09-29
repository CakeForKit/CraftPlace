package reqresp

import (
	"time"

	"github.com/google/uuid"
)

type AddPostRequest struct {
	Description string    `json:"description" binding:"required,max=255" example:"Магазин сережек"`
	ShopID      uuid.UUID `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
}

type UpdatePostRequest struct {
	Description string    `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
	ShopID      uuid.UUID `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
}

type PostResponse struct {
	ID              string    `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	Description     string    `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
	TimePublication time.Time `json:"timePublication" example:"2023-06-15T10:00:00Z"`
	ShopID          uuid.UUID `json:"shopID" binding:"required,uuid" example:"bb2e8400-e29b-41d4-a716-446655442222"`
}

type DeletePostRequest struct {
	ID string `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
}
