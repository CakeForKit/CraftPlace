package reqresp

type AddShopRequest struct {
	Title       string `json:"title" binding:"required,max=255" example:"Звезды"`
	Description string `json:"description" binding:"required,max=255" example:"Магазин сережек"`
}

type UpdateShopRequest struct {
	Title       string `json:"title" binding:"required,max=255" example:"Лучшие звезды"`
	Description string `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
}
