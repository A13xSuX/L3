package main

import (
	"context"
	"l3/CommentTree/internal/appcfg"
	"log"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func main() {

	cfg, err := appcfg.NewAppConfig()
	if err != nil {
		log.Fatal(err)
	}
	zlog.Init()
	err = zlog.SetLevel(cfg.LoggerConfig.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	opts := &dbpg.Options{
		MaxOpenConns:    cfg.PostgresConfig.MaxOpenConns,
		MaxIdleConns:    cfg.PostgresConfig.MaxIdleConns,
		ConnMaxLifetime: cfg.PostgresConfig.ConnMaxLifetime,
	}
	db, err := dbpg.New(cfg.PostgresConfig.MasterDSN, cfg.PostgresConfig.SlaveDSNs, opts)
	if err != nil {
		log.Fatal(err)
	}

	r := ginext.New("release")
	r.Use(ginext.Logger(), ginext.Recovery())
	r.GET("/health", func(c *ginext.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var res int
		err := db.QueryRowContext(ctx, "select 1").Scan(&res)
		if err != nil {
			c.JSON(500, ginext.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(200, ginext.H{"status": "ok"})
	})
	zlog.Logger.Info().Str("addr", cfg.ServerConfig.Addr).Msg("starting server")
	log.Fatal(r.Run(cfg.ServerConfig.Addr))
}
