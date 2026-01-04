package appcfg

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
)

type appConfig struct {
	serverConfig   serverConfig
	loggerConfig   loggerConfig
	postgresConfig postgresConfig
}

type serverConfig struct {
	addr string
}

type loggerConfig struct {
	logLevel string
}

type postgresConfig struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	port            int
}

func NewAppConfig() (*appConfig, error) {
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
	cfg.DefineFlag("p", "srvport", "transport.http.port", 7777, "HTTP server port")
	if err := cfg.ParseFlags(); err != nil {
		return nil, fmt.Errorf("failed to pars flags: %w", err)
	}

	// Распаковка в структуру
	var appConfig appConfig
	appConfig.serverConfig.addr = cfg.GetString("server.addr")
	appConfig.loggerConfig.logLevel = cfg.GetString("logger.level")
	appConfig.loggerConfig.logLevel = cfg.GetString("logger.level")
	appConfig.postgresConfig.maxOpenConns = cfg.GetInt("postgres.max_open_conns")
	appConfig.postgresConfig.maxIdleConns = cfg.GetInt("postgres.max_idle_conns")
	appConfig.postgresConfig.connMaxLifetime = cfg.GetDuration("postgres.conn_max_lifetime")
	appConfig.postgresConfig.port = cfg.GetInt("postgres.port") // из переменной окружения (из файла .env)

	return &appConfig, nil
}
