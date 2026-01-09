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
	KafkaConfig    KafkaConfig
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

// TODO
type KafkaConfig struct{}

func NewAppConfig() (*AppConfig, error) {
	//only env file
	envFilePath := "../.env"

	cfg := config.New()

	if err := cfg.LoadEnvFiles(envFilePath); err != nil {
		return nil, fmt.Errorf("failed to load env files: %w", err)
	}

	cfg.EnableEnv("")

	var appConfig AppConfig
	appConfig.ServerConfig.Addr = cfg.GetString("server.addr")
	appConfig.LoggerConfig.LogLevel = cfg.GetString("logger.level")
	appConfig.PostgresConfig.MaxOpenConns = cfg.GetInt("postgres.max_open_conns")
	appConfig.PostgresConfig.MaxIdleConns = cfg.GetInt("postgres.max_idle_conns")
	appConfig.PostgresConfig.ConnMaxLifeTime = cfg.GetDuration("postgres.conn_max_lifetime")
	appConfig.PostgresConfig.Port = cfg.GetInt("postgres.port")
	appConfig.PostgresConfig.MasterDSN = cfg.GetString("postgres.master_dsn")
	appConfig.PostgresConfig.SlaveDSN = cfg.GetStringSlice("postgres.slave")
	//kafka
	return &appConfig, nil
}
