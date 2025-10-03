package api

import (
	"errors"
	"net/http"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	auth "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
	authuser "github.com/CakeForKit/CraftPlace.git/internal/services/auth/auth_user"
	postservice "github.com/CakeForKit/CraftPlace.git/internal/services/post_service"
	productservice "github.com/CakeForKit/CraftPlace.git/internal/services/product_service"
	"github.com/CakeForKit/CraftPlace.git/internal/services/searcher"
	shopservice "github.com/CakeForKit/CraftPlace.git/internal/services/shop_service"
	userselfservice "github.com/CakeForKit/CraftPlace.git/internal/services/user_self_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserSelfRouter struct {
	userSelfServ userselfservice.UserSelfServ
	authz        auth.AuthZ
	searcherServ searcher.Searcher
	shopServ     shopservice.ShopServ
	productServ  productservice.ProductServ
	postServ     postservice.PostServ
}

func NewUserSelfRouter(
	router *gin.RouterGroup,
	userSelfServ userselfservice.UserSelfServ,
	authz auth.AuthZ,
	searcherServ searcher.Searcher,
	shopServ shopservice.ShopServ,
	productServ productservice.ProductServ,
	postServ postservice.PostServ,
) UserSelfRouter {
	r := UserSelfRouter{
		userSelfServ: userSelfServ,
		authz:        authz,
		searcherServ: searcherServ,
		shopServ:     shopServ,
		productServ:  productServ,
		postServ:     postServ,
	}
	gr := router.Group("user")
	gr.GET("/:id_user", r.GetUserByID)
	gr.PATCH("/update-login", r.UpdateLogin)
	gr.PATCH("/update-password", r.UpdatePassword)
	return r
}

// GetUserByID godoc
// @Summary Получить пользователя по ID
// @Description Возвращает информацию о пользователе по его идентификатору
// @Tags Пользователь
// @Accept json
// @Produce json
// @Param id_user path string true "ID пользователя" format(uuid)
// @Success 200 {object} reqresp.UserResponse "Информация о пользователе"
// @Failure 400 {object} map[string]interface{} "Неверный формат ID пользователя"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /user/{id_user} [get]
func (r *UserSelfRouter) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := uuid.Parse(c.Param("id_user"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}
	user, err := r.userSelfServ.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, authuser.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateLogin godoc
// @Summary Обновить логин пользователя
// @Description Изменяет логин текущего авторизованного пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.UpdateLoginRequest true "Данные для обновления логина"
// @Success 200 {object} map[string]interface{} "Успешное обновление"
// @Router /user/update-login [patch]
func (r *UserSelfRouter) UpdateLogin(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.UpdateLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newLogin := req.Login
	if err := r.userSelfServ.ChangeLogin(ctx, newLogin); err != nil {
		if errors.Is(err, auth.ErrNotAuthZ) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// UpdatePassword godoc
// @Summary Обновить пароль пользователя
// @Description Изменяет пароль текущего авторизованного пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.UpdateUserPasswordRequest true "Данные для обновления пароля"
// @Success 200 {object} map[string]interface{} "Успешное обновление"
// @Router /user/update-password [patch]
func (r *UserSelfRouter) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.UpdateUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPassword := req.Password
	if err := r.userSelfServ.ChangePassword(ctx, newPassword); err != nil {
		if errors.Is(err, auth.ErrNotAuthZ) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
