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
	return r
}

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

func (r *UserSelfRouter) UpdateProduct(c *gin.Context) {
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
