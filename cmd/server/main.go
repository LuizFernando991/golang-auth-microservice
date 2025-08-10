package main

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/LuizFernando991/golang-auth-microservice/internal/config"
	"github.com/LuizFernando991/golang-auth-microservice/internal/handler"
	"github.com/LuizFernando991/golang-auth-microservice/internal/repository"
	"github.com/LuizFernando991/golang-auth-microservice/internal/server"
	"github.com/LuizFernando991/golang-auth-microservice/internal/service"
)

func main() {
	cfg, err := config.LoadEnv(".env")
	if err != nil {
		panic(err)
	}

	logger := config.NewLogger("application")

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to db", err)
	}
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	//redis config
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisURL})

	if _, err := rdb.Ping(rdb.Context()).Result(); err != nil {
		log.Println("warning: redis ping failed:", err)
	}

	userRepo := repository.NewUserRepo(db)
	authSvc := service.NewAuthService(userRepo, cfg.JwtSecret, cfg.AccessTTL, cfg.RefreshTTL, cfg.BcryptCost)
	authHandler := handler.NewAuthHandler(authSvc)

	app := server.NewServer(cfg, logger, authHandler, rdb)
	if err := app.Run(); err != nil {
		logger.Error("server exited", err)
	}
}
