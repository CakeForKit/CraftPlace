package reqresp

type UpdateLoginRequest struct {
	Login string `json:"login" binding:"required,max=50" example:"uname"`
}

type UpdateUserPasswordRequest struct {
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}

type UserResponse struct {
	Login string `json:"login" binding:"required,max=50" example:"uname"`
}
