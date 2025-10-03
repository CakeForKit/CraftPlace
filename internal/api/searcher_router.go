package api

import (
	"errors"
	"net/http"
	"strconv"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
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
	gr := router.Group("/")
	gr.GET("/categories", r.GetCategories)
	gr.GET("/categories/:id_category", r.GetCategoryByID)
	gr.GET("/shops", r.GetShops)
	gr.GET("/shops/:id_shop", r.GetShopByID)
	gr.GET("/products", r.GetProducts)
	gr.GET("/posts", r.GetPosts)

	// gr.GET("/shops/:id_shop/posts", r.GetShopPosts)
	// gr.GET("/shops/:id_shop/products", r.GetShopProducts)
	return r
}

// GetCategories godoc
// @Summary Получить категории
// @Description Возвращает список категорий с возможностью фильтрации
// @Tags Поиск
// @Accept json
// @Produce json
// @Param title query string false "Фильтр по названию категории"
// @Success 200 {array} reqresp.CategoryResponse
// @Router /categories [get]
func (r *SearcherRouter) GetCategories(c *gin.Context) {
	ctx := c.Request.Context()

	filterOps := reqresp.CategoryFilter{
		Title: c.Query("title"),
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
	c.JSON(http.StatusOK, resp)
}

// GetCategoryByID godoc
// @Summary Получить категорию по ID
// @Description Возвращает информацию о категории по её идентификатору
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id_category path string true "ID категории" format(uuid)
// @Success 200 {object} reqresp.CategoryResponse "Информация о категории"
// @Failure 400 {object} map[string]interface{} "Неверный формат ID категории"
// @Failure 404 {object} map[string]interface{} "Категория не найдена"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /categories/{id_category} [get]
func (r *SearcherRouter) GetCategoryByID(c *gin.Context) {
	ctx := c.Request.Context()
	categoryID, err := uuid.Parse(c.Param("id_category"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	category, err := r.searcherServ.GetCategoruByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, searcher.ErrCategoryNotFound) {
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
// @Success 200 {array} reqresp.ShopResponse
// @Router /shops [get]
func (r *SearcherRouter) GetShops(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := uuid.Parse(c.Query("id_user")) // default = uuid.Nil
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filterOps := reqresp.ShopFilter{
		Title:  c.Query("title"),
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

// GetShopByID godoc
// @Summary Получить магазин по ID
// @Description Возвращает информацию о магазине по его идентификатору
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id_shop path string true "ID магазина" format(uuid)
// @Success 200 {object} reqresp.ShopResponse "Информация о магазине"
// @Failure 400 {object} map[string]interface{} "Неверный формат ID магазина"
// @Failure 404 {object} map[string]interface{} "Магазин не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /shops/{id_shop} [get]
func (r *SearcherRouter) GetShopByID(c *gin.Context) {
	ctx := c.Request.Context()
	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	shop, err := r.searcherServ.GetShopByID(ctx, shopID)
	if err != nil {
		if errors.Is(err, searcher.ErrShopNotFound) {
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
// @Success 200 {array} reqresp.ProductResponse "Список товаров"
// @Failure 400 {object} map[string]interface{} "Неверный формат параметров"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /products [get]
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
	filterOps := reqresp.ProductFilter{
		Title:      c.Query("title"), // default = ""
		MaxCost:    maxCost,
		MinCost:    minCost,
		ShopID:     shopID,
		CategoryID: categoryID,
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
	c.JSON(http.StatusOK, resp)
}

// GetPosts godoc
// @Summary Получить посты
// @Description Возвращает список постов с возможностью фильтрации по магазину
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id_shop query string false "Фильтр по ID магазина" format(uuid)
// @Success 200 {array} reqresp.PostResponse "Список постов"
// @Failure 400 {object} map[string]interface{} "Неверный формат ID магазина"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /posts [get]
func (r *SearcherRouter) GetPosts(c *gin.Context) {
	ctx := c.Request.Context()

	shopID, err := uuid.Parse(c.Query("id_shop")) // default = uuid.Nil
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filterOps := reqresp.PostFilter{
		ShopID: shopID,
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
	c.JSON(http.StatusOK, resp)
}

// // GetShopPosts godoc
// // @Summary Получить посты магазина
// // @Description Возвращает список постов указанного магазина
// // @Tags Поиск
// // @Accept json
// // @Produce json
// // @Param id_shop path string true "ID магазина" format(uuid)
// // @Success 200 {array} reqresp.PostResponse
// // @Router /shops/{id_shop}/posts [get]
// func (r *SearcherRouter) GetShopPosts(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	shopID, err := uuid.Parse(c.Param("id_shop"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id_shop format"})
// 		return
// 	}

// 	posts, err := r.searcherServ.GetPosts(ctx, shopID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	resp := make([]reqresp.PostResponse, len(posts))
// 	for i, v := range posts {
// 		resp[i] = v.ToResponse()
// 	}
// 	c.JSON(http.StatusOK, resp)
// }

// // GetShopProducts godoc
// // @Summary Получить товары магазина
// // @Description Возвращает список товаров указанного магазина с возможностью фильтрации
// // @Tags Поиск
// // @Accept json
// // @Produce json
// // @Param id_shop path string true "ID магазина" format(uuid)
// // @Param title query string false "Фильтр по названию товара"
// // @Param max_cost query integer false "Максимальная цена товара" default(100000)
// // @Param min_cost query integer false "Минимальная цена товара" default(0)
// // @Success 200 {array} reqresp.ProductResponse
// // @Router /shops/{id_shop}/products [get]
// func (r *SearcherRouter) GetShopProducts(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	shopID, err := uuid.Parse(c.Param("id_shop"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id_shop format"})
// 		return
// 	}

// 	maxCost, err := strconv.ParseUint(c.Query("max_cost"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	minCost, err := strconv.ParseUint(c.Query("min_cost"), 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	filterOps := reqresp.ProductFilter{
// 		Title:      c.Query("title"),
// 		MaxCost:    maxCost,
// 		MinCost:    minCost,
// 		ShopID:     shopID,
// 		CategoryID: uuid.Nil,
// 	}

// 	products, err := r.searcherServ.GetProducts(ctx, &filterOps)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	resp := make([]reqresp.ProductResponse, len(products))
// 	for i, v := range products {
// 		resp[i] = v.ToResponse()
// 	}
// 	c.JSON(http.StatusOK, resp)
// }
