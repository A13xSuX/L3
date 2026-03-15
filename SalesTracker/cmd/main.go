package main

import (
	"context"
	"l3/SalesTracker/internal/appCfg"
	"l3/SalesTracker/internal/handlers"
	"l3/SalesTracker/internal/repository"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.InitConsole()
	cfg, err := appCfg.NewAppConfig()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to get config")
		return
	}
	zlog.Logger.Info().Msg("Config upload")

	if err := zlog.SetLevel(cfg.LoggerConfig.Level); err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to set level of logger")
	}
	options := dbpg.Options{
		MaxOpenConns:    cfg.PostgresConfig.MaxOpenConns,
		MaxIdleConns:    cfg.PostgresConfig.MaxIdleConns,
		ConnMaxLifetime: cfg.PostgresConfig.ConnMaxLifetime,
	}

	db, err := dbpg.New(cfg.PostgresConfig.MasterDSN, cfg.PostgresConfig.SlaveDSN, &options)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed connect to db")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, "SELECT 1")

	var one int
	if err := row.Scan(&one); err != nil {
		zlog.Logger.Error().Err(err).Msg("Database healthcheck failed")
		return
	}
	zlog.Logger.Info().Msg("Success connect to db")

	salesRepo := repository.NewSalesRepo(db)
	handler := handlers.NewSaleHandler(salesRepo)
	//change to release
	router := ginext.New("debug")

	router.StaticFile("/", "../web/index.html")

	router.POST("/items", handler.Create)
	router.PUT("/items/:id", handler.Update)
	router.DELETE("/items/:id", handler.Delete)
	router.GET("/items", handler.GetAll)
	router.GET("/items/:id", handler.GetByID)
	router.GET("/analytics", handler.Analytics)

	if err := router.Run(cfg.ServerConfig.Addr); err != nil {
		zlog.Logger.Error().Err(err).Msg("server failed")
		return
	}
}
