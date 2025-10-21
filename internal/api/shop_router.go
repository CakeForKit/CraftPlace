package api

import (
	"errors"
	"net/http"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	shoprep "github.com/CakeForKit/CraftPlace.git/internal/repository/shop_rep"
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
	gr := router.Group("/shops")
	gr.POST("/", r.AddUserShop)
	gr.PUT("/:id_shop", r.UpdateShop)
	gr.DELETE("/:id_shop", r.DeleteShop)
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
// @Success 201 "Магазин успешно создан"
// @Failure 400 {object} ErrorResponse "Неверные данные запроса, неверный магазин"
// @Failure 401 {object} ErrorResponse "Неавторизованный доступ"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /shops [post]
func (r *ShopRouter) AddUserShop(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.AddShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dataAdd := models.AddShopData{
		Title:       req.Title,
		Description: req.Description,
	}
	_, err := r.shopServ.Add(ctx, dataAdd)
	if err != nil {
		if errors.Is(err, shopservice.ErrWrongShop) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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
// @Param id_shop path string true "ID магазина" format(uuid)
// @Param request body reqresp.UpdateShopRequest true "Данные для обновления магазина"
// @Success 204 "Магазин успешно обновлен"
// @Failure 400 {object} ErrorResponse "Неверный формат ID, неверные данные запроса, неверный магазин"
// @Failure 401 {object} ErrorResponse "Неавторизованный доступ"
// @Failure 404 {object} ErrorResponse "Магазин не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /shops/{id_shop} [put]
func (r *ShopRouter) UpdateShop(c *gin.Context) {
	ctx := c.Request.Context()
	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	var req reqresp.UpdateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateData := models.UpdateShopData{
		Title:       req.Title,
		Description: req.Description,
	}

	if _, err := r.shopServ.Update(ctx, shopID, updateData); err != nil {
		if errors.Is(err, shoprep.ErrShopNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, shopservice.ErrWrongShop) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

// DeleteShop godoc
// @Summary Удалить магазин
// @Description Удаляет магазин пользователя
// @Tags Магазины
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param id_shop path string true "ID магазина" format(uuid)
// @Success 204 "Магазин успешно удален"
// @Failure 400 {object} ErrorResponse "Неверный формат ID, неверный магазин"
// @Failure 401 {object} ErrorResponse "Неавторизованный доступ"
// @Failure 404 {object} ErrorResponse "Магазин не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /shops/{id_shop} [delete]
func (r *ShopRouter) DeleteShop(c *gin.Context) {
	ctx := c.Request.Context()
	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := r.shopServ.Delete(ctx, shopID); err != nil {
		if errors.Is(err, shoprep.ErrShopNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, shopservice.ErrWrongShop) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
