package main

import (
	"context"
	"fmt"
	"l3/CommentTree/internal/appcfg"
	"log"

	"github.com/wb-go/wbf/dbpg"
)

func main() {
	cfg, err := appcfg.NewAppConfig()
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
	ctx := context.Background()

}
