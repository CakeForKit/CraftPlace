package api

import (
	"errors"
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
	gr := router.Group("/shops/:id_shop/posts")
	gr.POST("/", r.AddPostToShop)
	gr.DELETE("/:id_post", r.DeletePost)
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
// @Param id_shop path string true "ID магазина" format(uuid)
// @Param request body reqresp.AddPostRequest true "Данные нового поста"
// @Success 201 "Пост успешно добавлен"
// @Failure 400 {object} ErrorResponse  "Неверный формат данных или ID, неверный магазин"
// @Failure 401 {object} ErrorResponse  "Неавторизованный доступ"
// @Failure 500 {object} ErrorResponse  "Внутренняя ошибка сервера"
// @Router /shops/{id_shop}/posts [post]
func (r *PostRouter) AddPostToShop(c *gin.Context) {
	ctx := c.Request.Context()

	var req reqresp.AddPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	dataToAdd := postservice.AddPostData{
		Description: req.Description,
		ShopID:      shopID,
	}
	if _, err := r.postServ.Add(ctx, dataToAdd); err != nil {
		if errors.Is(err, postservice.ErrWrongShop) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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
// @Param id_shop path string true "ID магазина" format(uuid)
// @Param id_post path string true "ID поста" format(uuid)
// @Success 204 "Пост успешно удален"
// @Failure 400 {object} ErrorResponse  "Неверный формат ID, неверный магазин"
// @Failure 401 {object} ErrorResponse  "Неавторизованный доступ"
// @Failure 404 {object} ErrorResponse  "Пост не найден"
// @Failure 500 {object} ErrorResponse  "Внутренняя ошибка сервера"
// @Router /shops/{id_shop}/posts/{id_post} [delete]
func (r *PostRouter) DeletePost(c *gin.Context) {
	ctx := c.Request.Context()

	shopID, err := uuid.Parse(c.Param("id_shop"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	postID, err := uuid.Parse(c.Param("id_post"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	if err := r.postServ.Delete(ctx, postID, shopID); err != nil {
		if errors.Is(err, postservice.ErrWrongShop) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
