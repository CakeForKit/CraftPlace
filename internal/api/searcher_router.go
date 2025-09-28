package api

import (
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
	gr.GET("/categories/:id_category/products", r.GetCategoryProducts)
	gr.GET("/shops", r.GetShops)
	gr.GET("/shops/:id_shop/posts", r.GetShopPosts)
	gr.GET("/shops/:id_shop/products", r.GetShopProducts)
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

// GetCategoryProducts godoc
// @Summary Получить товары категории
// @Description Возвращает список товаров указанной категории с возможностью фильтрации
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id_category path string true "ID категории" format(uuid)
// @Param title query string false "Фильтр по названию товара"
// @Param max_cost query integer false "Максимальная цена товара" default(100000)
// @Param min_cost query integer false "Минимальная цена товара" default(0)
// @Success 200 {array} reqresp.ProductResponse
// @Router /categories/{id_category}/products [get]
func (r *SearcherRouter) GetCategoryProducts(c *gin.Context) {
	ctx := c.Request.Context()
	categoryID, err := uuid.Parse(c.Param("id_category"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	maxCost, err := strconv.ParseUint(c.Query("max_cost"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	minCost, err := strconv.ParseUint(c.Query("min_cost"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filterOps := reqresp.ProductFilter{
		Title:      c.Query("title"),
		MaxCost:    maxCost,
		MinCost:    minCost,
		ShopID:     uuid.Nil,
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

// GetShops godoc
// @Summary Получить магазины
// @Description Возвращает список магазинов с возможностью фильтрации
// @Tags Поиск
// @Accept json
// @Produce json
// @Param title query string false "Фильтр по названию магазина"
// @Success 200 {array} reqresp.ShopResponse
// @Router /shops [get]
func (r *SearcherRouter) GetShops(c *gin.Context) {
	ctx := c.Request.Context()

	filterOps := reqresp.ShopFilter{
		Title:  c.Query("title"),
		UserID: uuid.Nil,
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

// GetShopPosts godoc
// @Summary Получить посты магазина
// @Description Возвращает список постов указанного магазина
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id_shop path string true "ID магазина" format(uuid)
// @Success 200 {array} reqresp.PostResponse
// @Router /shops/{id_shop}/posts [get]
func (r *SearcherRouter) GetShopPosts(c *gin.Context) {
	ctx := c.Request.Context()

	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id_shop format"})
		return
	}

	posts, err := r.searcherServ.GetPosts(ctx, shopID)
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

// GetShopProducts godoc
// @Summary Получить товары магазина
// @Description Возвращает список товаров указанного магазина с возможностью фильтрации
// @Tags Поиск
// @Accept json
// @Produce json
// @Param id_shop path string true "ID магазина" format(uuid)
// @Param title query string false "Фильтр по названию товара"
// @Param max_cost query integer false "Максимальная цена товара" default(100000)
// @Param min_cost query integer false "Минимальная цена товара" default(0)
// @Success 200 {array} reqresp.ProductResponse
// @Router /shops/{id_shop}/products [get]
func (r *SearcherRouter) GetShopProducts(c *gin.Context) {
	ctx := c.Request.Context()

	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id_shop format"})
		return
	}

	maxCost, err := strconv.ParseUint(c.Query("max_cost"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	minCost, err := strconv.ParseUint(c.Query("min_cost"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filterOps := reqresp.ProductFilter{
		Title:      c.Query("title"),
		MaxCost:    maxCost,
		MinCost:    minCost,
		ShopID:     shopID,
		CategoryID: uuid.Nil,
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
