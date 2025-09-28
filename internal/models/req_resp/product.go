package reqresp

type AddProductRequest struct {
	Title       string `json:"title" binding:"required,max=255" example:"Звезды"`
	Description string `json:"description" binding:"required,max=255" example:"Магазин сережек"`
	Cost        uint64 `json:"cost" binding:"required,min=0" example:"100"`
}

type UpdateProductRequest struct {
	Title       string `json:"title" binding:"required,max=255" example:"Лучшие звезды"`
	Description string `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
	Cost        uint64 `json:"cost" binding:"required,min=0" example:"200"`
}
