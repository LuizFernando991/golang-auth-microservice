package server

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/LuizFernando991/golang-auth-microservice/internal/config"
	"github.com/LuizFernando991/golang-auth-microservice/internal/handler"
	"github.com/LuizFernando991/golang-auth-microservice/internal/middleware"
	"github.com/go-redis/redis/v8"
)

type App struct {
	cfg     *config.Config
	logger  *config.Logger
	handler *handler.AuthHandler
	engine  *gin.Engine
	rdb     *redis.Client
}

func NewServer(cfg *config.Config, logger *config.Logger, authHandler *handler.AuthHandler, rdb *redis.Client) *App {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	app := &App{
		cfg:     cfg,
		logger:  logger,
		handler: authHandler,
		engine:  r,
		rdb:     rdb,
	}

	app.routes()
	return app
}

func (a *App) routes() {
	v1 := a.engine.Group("/v1")
	{
		v1.POST("/register", a.handler.Register)
		v1.POST("/login", a.handler.Login)
		v1.POST("/refresh", a.handler.Refresh)
		v1.POST("/logout", a.handler.Logout)

		// rate limiter applied to auth endpoints
		rl := middleware.NewRedisRateLimiter(a.rdb, a.logger, a.cfg.RateLimitRequests, a.cfg.RateLimitWindow)
		v1.Use(rl.Middleware())

		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(a.cfg.JwtSecret))
		protected.GET("/me", a.handler.Me)
	}
}

func (a *App) Run() error {
	addr := fmt.Sprintf(":%s", a.cfg.Port)
	return a.engine.Run(addr)
}
