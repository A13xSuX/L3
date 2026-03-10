package main

import (
	"fmt"
	"l3/EventBooker/internal/appcfg"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	cfg, err := appcfg.NewAppConfig()
	if err != nil {
		fmt.Printf("%w", err)
		return
	}
	zlog.InitConsole()
	//TODO уровень надо нормально передавать с конфига
	zlog.Logger.Level(0)
	opts := &dbpg.Options{
		MaxOpenConns: cfg.PostgresConfig.MaxOpenConns,
		MaxIdleConns: cfg.PostgresConfig.MaxIdleConns,
	}
	db, err := dbpg.New(cfg.PostgresConfig.MasterDSN, cfg.PostgresConfig.SlaveDSN, opts)
	if err != nil {
		zlog.Logger.Debug().Err(err)
		return
	}

	fmt.Println(cfg)
}
