package appCfg

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
)

type AppConfig struct {
	ServerConfig   ServerConfig
	LoggerConfig   LoggerConfig
	PostgresConfig PostgresConfig
	JWTConfig      JWTConfig
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
	ConnMaxLifetime time.Duration
	Port            int
	MasterDSN       string
	SlaveDSN        []string
}

type JWTConfig struct {
	SecretKey string
	TTL       time.Duration
}

func NewAppConfig() (*AppConfig, error) {
	appConfigFilePath := "../config.yaml"

	cfg := config.New()

	if err := cfg.LoadConfigFiles(appConfigFilePath); err != nil {
		return nil, fmt.Errorf("failed to load config files: %w", err)
	}

	var appConfig AppConfig
	appConfig.ServerConfig.Addr = cfg.GetString("server.addr")
	appConfig.LoggerConfig.LogLevel = cfg.GetString("logger.level")
	appConfig.PostgresConfig.Port = cfg.GetInt("postgres.port")
	appConfig.PostgresConfig.MaxOpenConns = cfg.GetInt("postgres.max_open_conns")
	appConfig.PostgresConfig.MaxIdleConns = cfg.GetInt("postgres.max_idle_conns")
	appConfig.PostgresConfig.ConnMaxLifetime = cfg.GetDuration("postgres.conn_max_lifetime")
	appConfig.PostgresConfig.MasterDSN = cfg.GetString("postgres.master_dsn")
	appConfig.PostgresConfig.SlaveDSN = cfg.GetStringSlice("postgres.slave_dsn")
	appConfig.JWTConfig.SecretKey = cfg.GetString("jwt.secret_key")
	appConfig.JWTConfig.TTL = cfg.GetDuration("jwt.ttl")
	return &appConfig, nil
}
