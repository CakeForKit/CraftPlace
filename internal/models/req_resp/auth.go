package reqresp

type LoginUserRequest struct {
	Login    string `json:"login" binding:"required,min=4,max=50" example:"ulogin"`
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required,max=50" example:"uname"`
	Login    string `json:"login" binding:"required,min=4,max=50" example:"ulogin"`
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}
