package main

import (
	"context"
	"l3/WarehouseControl/internal/appCfg"
	"l3/WarehouseControl/internal/auth"
	"l3/WarehouseControl/internal/handlers"
	"l3/WarehouseControl/internal/middleware"
	"l3/WarehouseControl/internal/repository"
	"l3/WarehouseControl/internal/service"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.InitConsole()

	cfg, err := appCfg.NewAppConfig()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to get config")
		return
	}
	zlog.Logger.Info().Msg("Success get config")

	err = zlog.SetLevel(cfg.LoggerConfig.LogLevel)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to set level logger")
	}

	opts := dbpg.Options{
		MaxOpenConns:    cfg.PostgresConfig.MaxOpenConns,
		MaxIdleConns:    cfg.PostgresConfig.MaxIdleConns,
		ConnMaxLifetime: cfg.PostgresConfig.ConnMaxLifetime,
	}
	db, err := dbpg.New(cfg.PostgresConfig.MasterDSN, cfg.PostgresConfig.SlaveDSN, &opts)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed connect to db")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, "SELECT 1")
	var one int
	if err := row.Scan(&one); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed connection to db")
		return
	}
	zlog.Logger.Info().Msg("Success connect to db")

	userRepo := repository.NewUserRepo(db)
	jwtService := auth.NewJWT(cfg.JWTConfig.SecretKey, cfg.JWTConfig.TTL)
	loginService := service.NewLoginService(userRepo, jwtService)
	loginHandler := handlers.NewLoginHandler(loginService)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	authHandler := handlers.NewAuthHandler()

	router := ginext.New("debug")

	router.POST("/auth/login", loginHandler.Login)
	router.GET("me", authMiddleware.Auth(), middleware.RequireRoles("admin"), authHandler.Me)

	if err := router.Run(cfg.ServerConfig.Addr); err != nil {
		zlog.Logger.Error().Err(err).Msg("server failed")
		return
	}
}
