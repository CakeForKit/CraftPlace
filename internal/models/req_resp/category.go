package reqresp

type AddCategoryRequest struct {
	Title       string `json:"title" binding:"required,max=255" example:"Звезды"`
	Description string `json:"description" binding:"required,max=255" example:"Магазин сережек"`
}

type UpdateCategoryRequest struct {
	Title       string `json:"title" binding:"required,max=255" example:"Лучшие звезды"`
	Description string `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
}

type CategoryResponse struct {
	ID          string `json:"id" example:"bb2e8400-e29b-41d4-a716-446655442222"`
	Title       string `json:"title" example:"Eco"`
	Description string `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
}
