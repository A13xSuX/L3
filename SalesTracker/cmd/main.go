package main

import (
	"fmt"
	"l3/SalesTracker/internal/appCfg"

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
	fmt.Println(cfg)
}
