package reqresp

type UpdateUsernameRequest struct {
	Username string `json:"username" binding:"required,max=50" example:"uname"`
}

type UpdateUserPasswordRequest struct {
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}
