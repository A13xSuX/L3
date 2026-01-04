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
	ConnMaxLifetime time.Duration
	Port            int
	MasterDSN       string
	SlaveDSNs       []string
}

func NewAppConfig() (*AppConfig, error) {
	envFilePath := "../../.env"                            //"./.env-example"
	appConfigFilePath := "../../app-config.yaml"           //"./config-example1.yaml"
	postgresConfigFilePath := "../../postgres-config.yaml" //"./config-example2.yaml"

	cfg := config.New()

	//Загрузка .env файлов
	if err := cfg.LoadEnvFiles(envFilePath); err != nil {
		return nil, fmt.Errorf("failed to load env files: %w", err)
	}

	// Включение поддержки переменных окружения
	cfg.EnableEnv("APP")

	// Загрузка файлов конфигурации
	if err := cfg.LoadConfigFiles(appConfigFilePath, postgresConfigFilePath); err != nil {
		return nil, fmt.Errorf("failed to load config files: %w", err)
	}

	// Определение флагов командной строки
	err := cfg.DefineFlag("p", "srvport", "transport.http.port", 7777, "HTTP server port")
	if err != nil {
		return nil, fmt.Errorf("failed to define flags: %w", err)

	}
	if err := cfg.ParseFlags(); err != nil {
		return nil, fmt.Errorf("failed to pars flags: %w", err)
	}

	// Распаковка в структуру
	var appConfig AppConfig
	appConfig.ServerConfig.Addr = cfg.GetString("server.addr")
	appConfig.LoggerConfig.LogLevel = cfg.GetString("logger.level")
	appConfig.PostgresConfig.MaxOpenConns = cfg.GetInt("postgres.max_open_conns")
	appConfig.PostgresConfig.MaxIdleConns = cfg.GetInt("postgres.max_idle_conns")
	appConfig.PostgresConfig.ConnMaxLifetime = cfg.GetDuration("postgres.conn_max_lifetime")
	appConfig.PostgresConfig.Port = cfg.GetInt("postgres.port") // из переменной окружения (из файла .env)
	appConfig.PostgresConfig.MasterDSN = cfg.GetString("postgres.master_dsn")
	appConfig.PostgresConfig.SlaveDSNs = cfg.GetStringSlice("postgres.slave_dsns")

	return &appConfig, nil
}
