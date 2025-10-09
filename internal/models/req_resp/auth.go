package reqresp

type LoginUserRequest struct {
	Login    string `json:"login" binding:"required,min=4,max=50" example:"ulogin"`
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token" example:"7249ede8-5083-4bd6-ad09-0b5fa3c5f2de"`
}

type RegisterUserRequest struct {
	Login    string `json:"login" binding:"required,min=4,max=50" example:"ulogin"`
	Password string `json:"password" binding:"required,min=4" example:"12345678"`
}
