package api

import (
	"errors"
	"net/http"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	authuser "github.com/CakeForKit/CraftPlace.git/internal/services/auth/auth_user"
	"github.com/gin-gonic/gin"
)

type AuthUserRouter struct {
	authu authuser.AuthUser
}

func NewAuthUserRouter(router *gin.RouterGroup, authu authuser.AuthUser) AuthUserRouter {
	r := AuthUserRouter{
		authu: authu,
	}
	gr := router.Group("auth-user")
	gr.POST("/register", r.Register)
	gr.POST("/login", r.Login)
	return r
}

// Register Handler
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя
// @Tags аутентификация
// @Accept json
// @Param request body reqresp.RegisterUserRequest true "Данные для регистрации"
// @Success 200 "Пользователь зарегистрирован"
// @Failure 400 "Неверные входные параметры"
// @Failure 401 "Ошибка аутентификации"
// @Failure 409 "Попытка повторной регистрации"
// @Router /auth-user/register [post]
func (r *AuthUserRouter) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.authu.RegisterUser(ctx, req); err != nil {
		if errors.Is(err, authuser.ErrDuplicateLoginUser) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// Login Handler
// @Summary Вход пользователя
// @Description Аутентифицирует пользователя и возвращает токен доступа
// @Tags аутентификация
// @Accept json
// @Param request body reqresp.LoginUserRequest true "Учетные данные для входа"
// @Success 200 "Пользователь успешно аутентифицирован"
// @Failure 400 "Неверные входные параметры"
// @Failure 401 "Ошибка аутентификации"
// @Router /auth-user/login [post]
func (r *AuthUserRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := r.authu.LoginUser(ctx, req)
	if err != nil {
		if errors.Is(err, authuser.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	rsp := reqresp.LoginUserResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, rsp)
}
