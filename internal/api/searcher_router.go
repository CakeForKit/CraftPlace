package api

import (
	"errors"
	"net/http"

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
	gr.GET("/shops/:id/posts", r.GetShopPosts)
	gr.GET("/shops/:id/products", r.GetShopProducts)
	return r
}

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

func (r *SearcherRouter) GetCategoryProducts(c *gin.Context) {
	ctx := c.Request.Context()
	categoryID, err := uuid.Parse(c.Param("id_category"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
		return
	}

	category, err := r.searcherServ.GetCategoruByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, searcher.ErrCategoryNotFpund) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, category.ToResponse())
}
