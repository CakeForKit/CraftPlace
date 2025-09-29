package api

import (
	"errors"
	"net/http"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	auth "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
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
	gr.PUT("/update-username", r.UpdateUsername)
	gr.PUT("/update-password", r.UpdatePassword)
	gr.GET("/user-shops", r.GetUserShops)
	gr.POST("/user-shops", r.AddUserShop)

	gr.PUT("/user-products/:id_shop", r.AddProductToShop)
	gr.PUT("/user-posts/:id_shop", r.AddPostToShop)

	gr.POST("/user-shops/:id_shop", r.UpdateShop)
	gr.POST("/user-products", r.UpdateProduct)

	gr.DELETE("/user-shops", r.DeleteShop)
	gr.DELETE("/user-products", r.DeleteProduct)
	gr.DELETE("/user-posts", r.DeletePost)
	return r
}

// UpdateUsername godoc
// @Summary Обновить имя пользователя
// @Description Изменяет имя текущего авторизованного пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.UpdateUsernameRequest true "Данные для обновления имени"
// @Success 200 {object} map[string]interface{} "Успешное обновление"
// @Router /user/update-username [put]
func (r *UserSelfRouter) UpdateUsername(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.UpdateUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUsername := req.Username
	if err := r.userSelfServ.ChangeName(ctx, newUsername); err != nil {
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
// @Router /user/update-password [put]
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

// GetUserShops godoc
// @Summary Получить магазины пользователя
// @Description Возвращает список магазинов текущего авторизованного пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Success 200 {array} reqresp.ShopResponse "Список магазинов пользователя"
// @Router /user/user-shops [get]
func (r *UserSelfRouter) GetUserShops(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := r.authz.UserIDFromContext(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNotAuthZ) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	filterOps := reqresp.ShopFilter{
		Title:  "",
		UserID: userID,
	}

	shops, err := r.searcherServ.GetShops(ctx, &filterOps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]reqresp.ShopResponse, len(shops))
	for i, v := range shops {
		resp[i] = v.ToResponse()
	}
	c.JSON(http.StatusOK, resp)
}

// AddUserShop godoc
// @Summary Добавить магазин пользователя
// @Description Создает новый магазин для текущего авторизованного пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.AddShopRequest true "Данные нового магазина"
// @Success 201 {object} map[string]interface{} "Магазин успешно создан"
// @Router /user/user-shops [post]
func (r *UserSelfRouter) AddUserShop(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.AddShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.shopServ.Add(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// AddProductToShop godoc
// @Summary Добавить товар в магазин
// @Description Добавляет новый товар в указанный магазин пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param id_shop path string true "ID магазина" format(uuid)
// @Param request body reqresp.AddProductRequest true "Данные нового товара"
// @Success 201 {object} map[string]interface{} "Товар успешно добавлен"
// @Router /user/user-products/{id_shop} [put]
func (r *UserSelfRouter) AddProductToShop(c *gin.Context) {
	ctx := c.Request.Context()

	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id_shop format"})
		return
	}

	var req reqresp.AddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ShopID = shopID

	if err := r.productServ.Add(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// AddPostToShop godoc
// @Summary Добавить пост в магазин
// @Description Добавляет новый пост в указанный магазин пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param id_shop path string true "ID магазина" format(uuid)
// @Param request body reqresp.AddPostRequest true "Данные нового поста"
// @Success 201 {object} map[string]interface{} "Пост успешно добавлен"
// @Router /user/user-posts/{id_shop} [put]
func (r *UserSelfRouter) AddPostToShop(c *gin.Context) {
	ctx := c.Request.Context()

	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id_shop format"})
		return
	}

	var req reqresp.AddPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ShopID = shopID

	if err := r.postServ.Add(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// UpdateShop godoc
// @Summary Обновить магазин
// @Description Обновляет данные указанного магазина пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param id_shop path string true "ID магазина" format(uuid)
// @Param request body reqresp.UpdateShopRequest true "Данные для обновления магазина"
// @Success 200 {object} map[string]interface{} "Магазин успешно обновлен"
// @Router /user/user-shops/{id_shop} [post]
func (r *UserSelfRouter) UpdateShop(c *gin.Context) {
	ctx := c.Request.Context()

	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id_shop format"})
		return
	}

	var req reqresp.UpdateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.shopServ.Update(ctx, shopID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// UpdateProduct godoc
// @Summary Обновить товар
// @Description Обновляет данные товара пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.UpdateProductRequest true "Данные для обновления товара"
// @Success 200 {object} map[string]interface{} "Товар успешно обновлен"
// @Router /user/user-products [post]
func (r *UserSelfRouter) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.productServ.Update(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteShop godoc
// @Summary Удалить магазин
// @Description Удаляет магазин пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.DeleteShopRequest true "Данные для удаления магазина"
// @Success 200 {object} map[string]interface{} "Магазин успешно удален"
// @Router /user/user-shops [delete]
func (r *UserSelfRouter) DeleteShop(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.DeleteShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shopID, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := r.shopServ.Delete(ctx, shopID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteProduct godoc
// @Summary Удалить товар
// @Description Удаляет товар пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.DeleteProductRequest true "Данные для удаления товара"
// @Success 200 {object} map[string]interface{} "Товар успешно удален"
// @Router /user/user-products [delete]
func (r *UserSelfRouter) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.DeleteProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := r.productServ.Delete(ctx, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeletePost godoc
// @Summary Удалить пост
// @Description Удаляет пост пользователя
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.DeletePostRequest true "Данные для удаления поста"
// @Success 200 {object} map[string]interface{} "Пост успешно удален"
// @Router /user/user-posts [delete]
func (r *UserSelfRouter) DeletePost(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.DeletePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postID, err := uuid.Parse(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := r.postServ.Delete(ctx, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
