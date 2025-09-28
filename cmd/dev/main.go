// @title CraftPlace
// @version 1.0
// @description API для платформы для мастеров ручной работы
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/CakeForKit/CraftPlace.git/docs"
	"github.com/CakeForKit/CraftPlace.git/internal/api"
	"github.com/CakeForKit/CraftPlace.git/internal/services/searcher"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	engine := gin.New()
	// Настройка CORS
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Можно указать конкретные домены вместо "*"
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	engine.OPTIONS("/*any", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// ----- Config ------
	appCnfg_Port := 8080
	// -------------------

	// для Swagger - НЕ ТРОГАТЬ
	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", appCnfg_Port))
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// ----- Services -----
	searcherServ := searcher.NewSearcher()
	// --------------------

	// ----- Groups -----
	apiGroup := engine.Group("/api/v1")
	// ------------------
	searcherRouter := api.NewSearcherRouter(apiGroup, searcherServ)
	_ = searcherRouter

	engine.Run(fmt.Sprintf(":%d", appCnfg_Port))
}
