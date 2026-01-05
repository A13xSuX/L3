package main

import (
	"l3/CommentTree/internal/appcfg"
	"l3/CommentTree/internal/handler"
	"log"

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

	handler.Register(r, handler.Deps{DB: db})
	
	zlog.Logger.Info().Str("addr", cfg.ServerConfig.Addr).Msg("starting server")
	log.Fatal(r.Run(cfg.ServerConfig.Addr))
}
