package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/CakeForKit/CraftPlace.git/internal/models/models"
	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	categoryrep "github.com/CakeForKit/CraftPlace.git/internal/repository/category_rep"
	shoprep "github.com/CakeForKit/CraftPlace.git/internal/repository/shop_rep"
	"github.com/CakeForKit/CraftPlace.git/internal/services/searcher"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SearcherRouter struct {
	searcherServ searcher.Searcher
}

func NewSearcherRouter(router *gin.RouterGroup, searcherServ searcher.Searcher) SearcherRouter {
	r := SearcherRouter{
		searcherServ: searcherServ,
	}
	gr := router.Group("craftplace")
	gr.GET("/categories", r.GetCategories)
	gr.GET("/categories/:id_category", r.GetCategoryByID)
	gr.GET("/shops", r.GetShops)
	gr.GET("/shops/:id_shop", r.GetShopByID)
	gr.GET("/products", r.GetProducts)
	gr.GET("/posts", r.GetPosts)

	return r
}

// GetCategories godoc
// @Summary Получить категории
// @Description Возвращает список категорий с возможностью фильтрации
// @Tags Поиск
// @Accept json
// @Produce json
// @Param title query string false "Фильтр по названию категории"
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param size query int false "Размер страницы" default(20) minimum(1) maximum(100)
// @Success 200 {object} reqresp.CategoriesResponse "Успешный ответ с категориями"
// @Failure 400 {object} ErrorResponse "Неверный формат параметров пагинации"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /craftplace/categories [get]
func (r *SearcherRouter) GetCategories(c *gin.Context) {
	ctx := c.Request.Context()

	page, err := strconv.ParseUint(c.Query("page"), 10, 64)
	if err != nil || page == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	size, err := strconv.ParseUint(c.Query("size"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset := (page - 1) * size
	limit := size
	filterOps := models.CategoryFilter{
		Title:  c.Query("title"),
		Offset: offset,
		Limit:  limit,
	}

	caterories, err := r.searcherServ.GetCategories(ctx, &filterOps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]reqresp.CategoryResponse, len(caterories))
	for i, v := range caterories {
		resp[i] = v.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"Categories": reqresp.CategoriesResponse{
		Categories: resp,
	}})
}

// GetCategoryByID godoc
// @Summary Получить категорию по ID
// @Description Возвращает информацию о категории по её идентификатору
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id_category path string true "ID категории" format(uuid)
// @Success 200 {object} reqresp.CategoryResponse "Информация о категории"
// @Failure 400 {object} ErrorResponse "Неверный формат ID категории"
// @Failure 404 {object} ErrorResponse "Категория не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /craftplace/categories/{id_category} [get]
func (r *SearcherRouter) GetCategoryByID(c *gin.Context) {
	ctx := c.Request.Context()
	categoryID, err := uuid.Parse(c.Param("id_category"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	category, err := r.searcherServ.GetCategoryByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, categoryrep.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, category.ToResponse())
}

// GetShops godoc
// @Summary Получить магазины
// @Description Возвращает список магазинов с возможностью фильтрации
// @Tags Поиск
// @Accept json
// @Produce json
// @Param title query string false "Фильтр по названию магазина"
// @Param id_user query string false "Фильтр по ID пользователя" format(uuid) default(00000000-0000-0000-0000-000000000000)
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param size query int false "Размер страницы" default(20) minimum(1) maximum(100)
// @Success 200 {object} reqresp.ShopsResponse
// @Failure 400 {object} ErrorResponse "Неверный формат параметров пагинации или ID пользователя"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /craftplace/shops [get]
func (r *SearcherRouter) GetShops(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := uuid.Parse(c.Query("id_user")) // default = uuid.Nil
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, err := strconv.ParseUint(c.Query("page"), 10, 64)
	if err != nil || page == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	size, err := strconv.ParseUint(c.Query("size"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset := (page - 1) * size
	limit := size
	filterOps := models.ShopFilter{
		Title:  c.Query("title"),
		UserID: userID,
		Offset: offset,
		Limit:  limit,
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
	c.JSON(http.StatusOK, reqresp.ShopsResponse{
		Shops: resp,
	})
}

// GetShopByID godoc
// @Summary Получить магазин по ID
// @Description Возвращает информацию о магазине по его идентификатору
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id_shop path string true "ID магазина" format(uuid)
// @Success 200 {object} reqresp.ShopResponse "Информация о магазине"
// @Failure 400 {object} ErrorResponse "Неверный формат ID магазина"
// @Failure 404 {object} ErrorResponse "Магазин не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /craftplace/shops/{id_shop} [get]
func (r *SearcherRouter) GetShopByID(c *gin.Context) {
	ctx := c.Request.Context()
	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	shop, err := r.searcherServ.GetShopByID(ctx, shopID)
	if err != nil {
		if errors.Is(err, shoprep.ErrShopNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, shop.ToResponse())
}

// GetProducts godoc
// @Summary Получить товары
// @Description Возвращает список товаров с возможностью фильтрации по различным параметрам
// @Tags Поиск
// @Accept json
// @Produce json
// @Param title query string false "Фильтр по названию товара"
// @Param min_cost query integer false "Минимальная цена товара" default(0)
// @Param max_cost query integer false "Максимальная цена товара" default(100000)
// @Param id_shop query string false "Фильтр по ID магазина" format(uuid) default(00000000-0000-0000-0000-000000000000)
// @Param id_category query string false "Фильтр по ID категории" format(uuid) default(00000000-0000-0000-0000-000000000000)
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param size query int false "Размер страницы" default(20) minimum(1) maximum(100)
// @Success 200 {object} reqresp.ProductsResponse "Список товаров"
// @Failure 400 {object} ErrorResponse "Неверный формат параметров (цена, ID, пагинация)"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /craftplace/products [get]
func (r *SearcherRouter) GetProducts(c *gin.Context) {
	ctx := c.Request.Context()

	minCost, err := strconv.ParseUint(c.Query("min_cost"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	maxCost, err := strconv.ParseUint(c.Query("max_cost"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shopID, err := uuid.Parse(c.Query("id_shop")) // default = uuid.Nil
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	categoryID, err := uuid.Parse(c.Query("id_category")) // default = uuid.Nil
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, err := strconv.ParseUint(c.Query("page"), 10, 64)
	if err != nil || page == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	size, err := strconv.ParseUint(c.Query("size"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset := (page - 1) * size
	limit := size
	filterOps := models.ProductFilter{
		Title:      c.Query("title"), // default = ""
		MaxCost:    maxCost,
		MinCost:    minCost,
		ShopID:     shopID,
		CategoryID: categoryID,
		Offset:     offset,
		Limit:      limit,
	}

	products, err := r.searcherServ.GetProducts(ctx, &filterOps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]reqresp.ProductResponse, len(products))
	for i, v := range products {
		resp[i] = v.ToResponse()
	}
	c.JSON(http.StatusOK, reqresp.ProductsResponse{
		Products: resp,
	})
}

// GetPosts godoc
// @Summary Получить посты
// @Description Возвращает список постов с возможностью фильтрации по магазину
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id_shop query string false "Фильтр по ID магазина" format(uuid) default(00000000-0000-0000-0000-000000000000)
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param size query int false "Размер страницы" default(20) minimum(1) maximum(100)
// @Success 200 {object} reqresp.PostsResponse "Список постов"
// @Failure 400 {object} ErrorResponse "Неверный формат ID магазина или параметров пагинации"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /craftplace/posts [get]
func (r *SearcherRouter) GetPosts(c *gin.Context) {
	ctx := c.Request.Context()

	shopID, err := uuid.Parse(c.Query("id_shop")) // default = uuid.Nil
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, err := strconv.ParseUint(c.Query("page"), 10, 64)
	if err != nil || page == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	size, err := strconv.ParseUint(c.Query("size"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset := (page - 1) * size
	limit := size
	filterOps := models.PostFilter{
		ShopID: shopID,
		Offset: offset,
		Limit:  limit,
	}

	posts, err := r.searcherServ.GetPosts(ctx, &filterOps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]reqresp.PostResponse, len(posts))
	for i, v := range posts {
		resp[i] = v.ToResponse()
	}
	c.JSON(http.StatusOK, reqresp.PostsResponse{
		Posts: resp,
	})
}
