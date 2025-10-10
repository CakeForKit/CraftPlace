package api

import (
	"errors"
	"net/http"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
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
	gr := router.Group("shops/:id_shop/products")
	gr.POST("/", r.AddProductToShop)
	gr.PUT("/:id_product", r.UpdateProduct)
	gr.DELETE("/:id_product", r.DeleteProduct)

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
// @Param id_shop path string true "ID магазина" format(uuid)
// @Param request body reqresp.AddProductRequest true "Данные нового товара"
// @Success 201 "Товар успешно добавлен"
// @Failure 400 {object} ErrorResponse  "Неверный формат данных или ID, неверный магазин"
// @Failure 401 {object} ErrorResponse  "Неавторизованный доступ"
// @Failure 500 {object} ErrorResponse  "Внутренняя ошибка сервера"
// @Router /shops/{id_shop}/products [post]
func (r *ProductRouter) AddProductToShop(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.AddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	dataToAdd := models.AddProductData{
		Title:       req.Title,
		Description: req.Description,
		Cost:        req.Cost,
		ShopID:      shopID,
		CategoryIDs: req.CategoryIDs,
	}
	if _, err := r.productServ.Add(ctx, dataToAdd); err != nil {
		if errors.Is(err, productservice.ErrWrongShop) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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
// @Param id_shop path string true "ID магазина" format(uuid)
// @Param id_product path string true "ID изделия" format(uuid)
// @Param request body reqresp.UpdateProductRequest true "Данные для обновления товара"
// @Success 200 "Товар успешно обновлен"
// @Failure 400 {object} ErrorResponse  "Неверный формат данных или ID, неверный магазин"
// @Failure 401 {object} ErrorResponse  "Неавторизованный доступ"
// @Failure 404 {object} ErrorResponse  "Товар не найден"
// @Failure 500 {object} ErrorResponse  "Внутренняя ошибка сервера"
// @Router /shops/{id_shop}/products/{id_product} [put]
func (r *ProductRouter) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	productID, err := uuid.Parse(c.Param("id_product"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	dataToUpdate := models.UpdateProductData{
		Title:       req.Title,
		Description: req.Description,
		Cost:        req.Cost,
		ShopID:      shopID,
		CategoryIDs: req.CategoryIDs,
	}
	if _, err := r.productServ.Update(ctx, productID, dataToUpdate); err != nil {
		if errors.Is(err, productservice.ErrWrongShop) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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
// @Param id_shop path string true "ID магазина" format(uuid)
// @Param id_product path string true "ID изделия" format(uuid)
// @Success 200 "Товар успешно удален"
// @Failure 400 {object} ErrorResponse  "Неверный формат ID, неверный магазин"
// @Failure 401 {object} ErrorResponse  "Неавторизованный доступ"
// @Failure 404 {object} ErrorResponse  "Товар не найден"
// @Failure 500 {object} ErrorResponse  "Внутренняя ошибка сервера"
// @Router /shops/{id_shop}/products/{id_product} [delete]
func (r *ProductRouter) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()
	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	productID, err := uuid.Parse(c.Param("id_product"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	if err := r.productServ.Delete(ctx, productID, shopID); err != nil {
		if errors.Is(err, productservice.ErrWrongShop) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
