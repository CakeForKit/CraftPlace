package api

import (
	"net/http"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	shopservice "github.com/CakeForKit/CraftPlace.git/internal/services/shop_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShopRouter struct {
	shopServ shopservice.ShopServ
}

func NewShopRouter(
	router *gin.RouterGroup,
	shopServ shopservice.ShopServ,
) ShopRouter {
	r := ShopRouter{
		shopServ: shopServ,
	}
	gr := router.Group("/user-shops")
	gr.POST("/", r.AddUserShop)
	gr.PUT("/", r.UpdateShop)
	gr.DELETE("/", r.DeleteShop)
	return r
}

// AddUserShop godoc
// @Summary Добавить магазин
// @Description Создает новый магазин для текущего авторизованного пользователя
// @Tags Магазины
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.AddShopRequest true "Данные нового магазина"
// @Success 201 {object} map[string]interface{} "Магазин успешно создан"
// @Router /user/user-shops [post]
func (r *ShopRouter) AddUserShop(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.AddShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := r.shopServ.Add(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// UpdateShop godoc
// @Summary Обновить магазин
// @Description Обновляет данные указанного магазина пользователя
// @Tags Магазины
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.UpdateShopRequest true "Данные для обновления магазина"
// @Success 200 {object} map[string]interface{} "Магазин успешно обновлен"
// @Router /user/user-shops [put]
func (r *ShopRouter) UpdateShop(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.UpdateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := r.shopServ.Update(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteShop godoc
// @Summary Удалить магазин
// @Description Удаляет магазин пользователя
// @Tags Магазины
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.DeleteShopRequest true "Данные для удаления магазина"
// @Success 200 {object} map[string]interface{} "Магазин успешно удален"
// @Router /user/user-shops [delete]
func (r *ShopRouter) DeleteShop(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.DeleteShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shopID, err := uuid.Parse(req.ShopID)
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
