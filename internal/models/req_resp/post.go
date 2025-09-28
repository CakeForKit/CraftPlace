package reqresp

type AddPostRequest struct {
	Description string `json:"description" binding:"required,max=255" example:"Магазин сережек"`
}

type UpdatePostRequest struct {
	Description string `json:"description" binding:"required,max=255" example:"Лучший магазин сережек"`
}
