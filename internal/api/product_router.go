package api

import (
	"net/http"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	productservice "github.com/CakeForKit/CraftPlace.git/internal/services/product_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductRouter struct {
	productServ productservice.ProductServ
}

func NewProductRouter(
	router *gin.RouterGroup,
	productServ productservice.ProductServ,
) ProductRouter {
	r := ProductRouter{
		productServ: productServ,
	}
	gr := router.Group("products")
	gr.POST("/", r.AddProductToShop)
	gr.PUT("/", r.UpdateProduct)
	gr.DELETE("/", r.DeleteProduct)

	return r
}

// AddProductToShop godoc
// @Summary Добавить товар в магазин
// @Description Добавляет новый товар в указанный магазин пользователя
// @Tags Изделия
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.AddProductRequest true "Данные нового товара"
// @Success 201 {object} map[string]interface{} "Товар успешно добавлен"
// @Router /user/user-products [post]
func (r *ProductRouter) AddProductToShop(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.AddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.productServ.Add(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// UpdateProduct godoc
// @Summary Обновить товар
// @Description Обновляет данные товара пользователя
// @Tags Изделия
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.UpdateProductRequest true "Данные для обновления товара"
// @Success 200 {object} map[string]interface{} "Товар успешно обновлен"
// @Router /user/user-products [put]
func (r *ProductRouter) UpdateProduct(c *gin.Context) {
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

// DeleteProduct godoc
// @Summary Удалить товар
// @Description Удаляет товар пользователя
// @Tags Изделия
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.DeleteProductRequest true "Данные для удаления товара"
// @Success 200 {object} map[string]interface{} "Товар успешно удален"
// @Router /user/user-products [delete]
func (r *ProductRouter) DeleteProduct(c *gin.Context) {
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
