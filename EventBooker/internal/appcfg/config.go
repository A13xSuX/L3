package appcfg

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
)

type AppConfig struct {
	ServerConfig   ServerConfig
	LoggerConfig   LoggerConfig
	PostgresConfig PostgresConfig
}

type ServerConfig struct {
	Addr string
}

type LoggerConfig struct {
	LogLevel string
}

type PostgresConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifeTime time.Duration
	Port            int
	MasterDSN       string
	SlaveDSN        []string
}

func NewAppConfig() (*AppConfig, error) {
	envFilePath := "../.env"
	postgresConfigFilePath := "../config-postgres.yaml"

	cfg := config.New()

	if err := cfg.LoadEnvFiles(envFilePath); err != nil {
		return nil, fmt.Errorf("failed to load env files: %w", err)
	}

	cfg.EnableEnv("")

	if err := cfg.LoadConfigFiles(postgresConfigFilePath); err != nil {
		return nil, fmt.Errorf("failed to load config files: %w", err)
	}

	cfg.DefineFlag("p", "srvport", "transport.port.http", 8080, "HTTP server port")

	var appConfig AppConfig
	appConfig.ServerConfig.Addr = cfg.GetString("server.addr")
	appConfig.LoggerConfig.LogLevel = cfg.GetString("logger.level")
	appConfig.PostgresConfig.MaxOpenConns = cfg.GetInt("postgres.max_open_conns")
	appConfig.PostgresConfig.MaxIdleConns = cfg.GetInt("postgres.max_idle_conns")
	appConfig.PostgresConfig.Port = cfg.GetInt("postgres.port")
	appConfig.PostgresConfig.ConnMaxLifeTime = cfg.GetDuration("postgres.conn_max_life_time")
	appConfig.PostgresConfig.MasterDSN = cfg.GetString("postgres.master_DSN")
	appConfig.PostgresConfig.SlaveDSN = cfg.GetStringSlice("postgres.slave_DSN")
	return &appConfig, nil
}
