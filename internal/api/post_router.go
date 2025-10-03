package api

import (
	"net/http"

	reqresp "github.com/CakeForKit/CraftPlace.git/internal/models/req_resp"
	postservice "github.com/CakeForKit/CraftPlace.git/internal/services/post_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostRouter struct {
	postServ postservice.PostServ
}

func NewPostRouter(
	router *gin.RouterGroup,
	postServ postservice.PostServ,
) PostRouter {
	r := PostRouter{
		postServ: postServ,
	}
	gr := router.Group("posts")
	gr.POST("/", r.AddPostToShop)
	gr.DELETE("/", r.DeletePost)
	return r
}

// AddPostToShop godoc
// @Summary Добавить пост в магазин
// @Description Добавляет новый пост в указанный магазин пользователя
// @Tags Посты
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.AddPostRequest true "Данные нового поста"
// @Success 201 {object} map[string]interface{} "Пост успешно добавлен"
// @Router /user/user-posts [post]
func (r *PostRouter) AddPostToShop(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.AddPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.postServ.Add(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// DeletePost godoc
// @Summary Удалить пост
// @Description Удаляет пост пользователя
// @Tags Посты
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer токен"
// @Param request body reqresp.DeletePostRequest true "Данные для удаления поста"
// @Success 200 {object} map[string]interface{} "Пост успешно удален"
// @Router /user/user-posts [delete]
func (r *PostRouter) DeletePost(c *gin.Context) {
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
