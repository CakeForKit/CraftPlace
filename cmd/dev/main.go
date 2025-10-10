// @title CraftPlace
// @version 1.0
// @description API для платформы для мастеров ручной работы
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "github.com/CakeForKit/CraftPlace.git/docs"
	"github.com/CakeForKit/CraftPlace.git/internal/api"
	"github.com/CakeForKit/CraftPlace.git/internal/cnfg"
	"github.com/CakeForKit/CraftPlace.git/internal/middleware"
	postrep "github.com/CakeForKit/CraftPlace.git/internal/repository/post_rep"
	productrep "github.com/CakeForKit/CraftPlace.git/internal/repository/product_rep"
	shoprep "github.com/CakeForKit/CraftPlace.git/internal/repository/shop_rep"
	userrep "github.com/CakeForKit/CraftPlace.git/internal/repository/user_rep"
	auth "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
	authuser "github.com/CakeForKit/CraftPlace.git/internal/services/auth/auth_user"
	"github.com/CakeForKit/CraftPlace.git/internal/services/auth/hasher"
	tokenmaker "github.com/CakeForKit/CraftPlace.git/internal/services/auth/token_maker"
	postservice "github.com/CakeForKit/CraftPlace.git/internal/services/post_service"
	productservice "github.com/CakeForKit/CraftPlace.git/internal/services/product_service"
	"github.com/CakeForKit/CraftPlace.git/internal/services/searcher"
	shopservice "github.com/CakeForKit/CraftPlace.git/internal/services/shop_service"
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
	appCnfg, err := cnfg.LoadAppConfig("./configs/", "app_config", "yaml")
	if err != nil {
		panic(fmt.Errorf("cannot load AppConfig: %v", err))
	}
	pgCredentials, err := cnfg.LoadPgCredentials("./configs/", "db_config", "env")
	if err != nil {
		panic(fmt.Errorf("cannot load PgCredentials: %v", err))
	}
	dbConnCnfg, err := cnfg.LoadDatebaseConnConfig("./configs/", "app_config", "yaml")
	if err != nil {
		panic(fmt.Errorf("cannot load DatebaseConnConfig: %v", err))
	}
	// -------------------

	// для Swagger - НЕ ТРОГАТЬ
	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", appCnfg.Port))
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// ----- Repositories -----
	ctx := context.Background()
	userRep, err := userrep.NewPgUserRep(ctx, pgCredentials, dbConnCnfg)
	if err != nil {
		panic(err.Error())
	}
	shopRep, err := shoprep.NewPgShopRep(ctx, pgCredentials, dbConnCnfg)
	if err != nil {
		panic(err.Error())
	}
	postRep, err := postrep.NewPgPostRep(ctx, pgCredentials, dbConnCnfg)
	if err != nil {
		panic(err.Error())
	}
	productRep, err := productrep.NewPgProductRep(ctx, pgCredentials, dbConnCnfg)
	if err != nil {
		panic(err.Error())
	}
	// --------------------
	// ----- Services -----
	tokenMaker, err := tokenmaker.NewTokenMaker(appCnfg.TokenSymmetricKey)
	if err != nil {
		panic(err.Error())
	}
	hasher, err := hasher.NewHasher()
	if err != nil {
		panic(err.Error())
	}
	authUserServ := authuser.NewAuthUser(tokenMaker, hasher, appCnfg, userRep)
	authz, err := auth.NewAuthZ()
	if err != nil {
		panic(err.Error())
	}
	shopServ := shopservice.NewShopServ(shopRep, authz)
	searcherServ := searcher.NewSearcher()
	postServ := postservice.NewPostServ(postRep, authz, shopRep)
	productServ := productservice.NewProductServ(productRep, authz, shopRep)
	// --------------------

	// ----- Groups -----
	apiGroup := engine.Group("/api/v1")
	usersGroup := apiGroup.Group("/")
	usersGroup.Use(middleware.AuthMiddleware(authUserServ, authz))
	// ------------------
	searcherRouter := api.NewSearcherRouter(apiGroup, searcherServ)
	_ = searcherRouter
	authUserRouter := api.NewAuthUserRouter(apiGroup, authUserServ)
	_ = authUserRouter
	shopRouter := api.NewShopRouter(usersGroup, shopServ)
	_ = shopRouter
	postRouter := api.NewPostRouter(usersGroup, postServ)
	_ = postRouter
	productRouter := api.NewProductRouter(usersGroup, productServ)
	_ = productRouter
	// userSelfRouter := api.NewUserSelfRouter(apiGroup)

	engine.Run(fmt.Sprintf(":%d", appCnfg.Port))
}
