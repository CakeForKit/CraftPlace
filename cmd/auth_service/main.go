// auth_service/main.go
// @title Auth Service
// @version 1.0
// @description Authentication microservice
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/CakeForKit/CraftPlace.git/docs"
	"github.com/CakeForKit/CraftPlace.git/internal/api"
	"github.com/CakeForKit/CraftPlace.git/internal/cnfg"
	"github.com/CakeForKit/CraftPlace.git/internal/grpc/authgrpc"
	"github.com/CakeForKit/CraftPlace.git/internal/grpc/loggergrpc"
	"github.com/CakeForKit/CraftPlace.git/internal/middleware"
	userrep "github.com/CakeForKit/CraftPlace.git/internal/repository/user_rep"
	authz "github.com/CakeForKit/CraftPlace.git/internal/services/auth/authZ"
	authuser "github.com/CakeForKit/CraftPlace.git/internal/services/auth/auth_user"
	"github.com/CakeForKit/CraftPlace.git/internal/services/auth/hasher"
	tokenmaker "github.com/CakeForKit/CraftPlace.git/internal/services/auth/token_maker"
	"github.com/CakeForKit/CraftPlace.git/proto/pb/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

func main() {
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
	GRPCPort, err := cnfg.LoadGRPCConfig()
	if err != nil {
		panic(fmt.Errorf("cannot load GRPCPort: %v", err))
	}
	// -------------------

	// ----- Repositories -----
	ctx := context.Background()
	userRep, err := userrep.NewPgUserRep(ctx, pgCredentials, dbConnCnfg)
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
	authz, err := authz.NewAuthZ()
	if err != nil {
		panic(err.Error())
	}
	// --------------------

	// Создаем канал для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запускаем gRPC и HTTP серверы в отдельных горутинах
	fmt.Printf("GRPC linstening on %s\n", fmt.Sprintf(":%d", GRPCPort))
	grpcServer := startGRPCServer(authUserServ, fmt.Sprintf(":%d", GRPCPort))
	httpServer := startHTTPServer(authUserServ, authz, appCnfg)

	// Ожидаем сигнал завершения
	<-stop
	log.Println("Shutting down servers...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Servers stopped gracefully")
}

func startGRPCServer(authService authuser.AuthUser, grpcPort string) *grpc.Server {
	logger := middleware.StructuredLogger("CraftPlaceServ")
	maxMsgSize := 1024 * 1024 * 100 // 100 MB
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
		grpc.ChainUnaryInterceptor(
			loggergrpc.GRPCLoggingInterceptor(logger),
		),
	)
	authGrpcServer := authgrpc.NewServer(authService)

	auth.RegisterAuthServiceServer(grpcServer, authGrpcServer)

	go func() {
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("failed to listen for gRPC: %v", err)
		}

		log.Printf("gRPC Auth Service listening on %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	return grpcServer
}

func startHTTPServer(authService authuser.AuthUser, authz authz.AuthZ, appCnfg *cnfg.AppConfig) *http.Server {
	engine := gin.New()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	engine.OPTIONS("/*any", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})

	// Logger
	logger := middleware.StructuredLogger("CraftPlaceServ")
	engine.Use(middleware.LogMiddleware(logger))
	fmt.Printf("--------\n\n\n")

	contextPathGroup := engine.Group(appCnfg.ContextPath)

	// Swagger
	swaggerURL := fmt.Sprintf("http://localhost:%d%sswagger/doc.json", appCnfg.SwaggerPort, appCnfg.ContextPath)
	url := ginSwagger.URL(swaggerURL)
	fmt.Printf("SWAGGER url: %s\n\n", swaggerURL)
	contextPathGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// ----- Groups -----
	apiGroup := contextPathGroup.Group("/api/v1")
	usersGroup := apiGroup.Group("/")
	usersGroup.Use(middleware.AuthMiddleware(authService, authz))
	// ------------------

	authUserRouter := api.NewAuthUserRouter(apiGroup, authService)
	_ = authUserRouter

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", appCnfg.Port),
		Handler: engine,
	}

	go func() {
		log.Printf("HTTP REST API listening on :%d", appCnfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	return server
}

func main2() {
	engine := gin.New()
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

	// Logger
	logger := middleware.StructuredLogger("CraftPlaceServ")
	engine.Use(middleware.LogMiddleware(logger))
	fmt.Printf("--------\n\n\n")

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

	// ----- Repositories -----
	ctx := context.Background()
	userRep, err := userrep.NewPgUserRep(ctx, pgCredentials, dbConnCnfg)
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
	authz, err := authz.NewAuthZ()
	if err != nil {
		panic(err.Error())
	}
	// --------------------

	contextPathGroup := engine.Group(appCnfg.ContextPath)
	// для Swagger - НЕ ТРОГАТЬ
	swaggerURL := fmt.Sprintf("http://localhost:%d%sswagger/doc.json", appCnfg.SwaggerPort, appCnfg.ContextPath)
	url := ginSwagger.URL(swaggerURL)
	fmt.Printf("SWAGGER url: %s\n\n", swaggerURL)
	contextPathGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// ----- Groups -----
	apiGroup := contextPathGroup.Group("/api/v1")
	usersGroup := apiGroup.Group("/")
	usersGroup.Use(middleware.AuthMiddleware(authUserServ, authz))
	// ------------------

	authUserRouter := api.NewAuthUserRouter(apiGroup, authUserServ)
	_ = authUserRouter

	engine.Run(fmt.Sprintf(":%d", appCnfg.Port))
}
