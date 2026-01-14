package main

import (
	"context"
	"l3/ImageProcessor/internal/appcfg"
	"l3/ImageProcessor/internal/httpserver"
	"l3/ImageProcessor/internal/worker"
	"l3/ImageProcessor/repo"
	"os"
	"path/filepath"
	"time"

	"github.com/wb-go/wbf/dbpg"
	kafkawbf "github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	cfg, err := appcfg.NewAppConfig()
	if err != nil {
		zlog.Logger.Error().Err(err)
		return
	}

	zlog.InitConsole()
	_ = zlog.SetLevel(cfg.LoggerConfig.LogLevel)
	opts := &dbpg.Options{MaxOpenConns: cfg.PostgresConfig.MaxOpenConns,
		MaxIdleConns:    cfg.PostgresConfig.MaxIdleConns,
		ConnMaxLifetime: cfg.PostgresConfig.ConnMaxLifeTime,
	}
	db, err := dbpg.New(cfg.PostgresConfig.MasterDSN, cfg.PostgresConfig.SlaveDSN, opts)
	if err != nil {
		zlog.Logger.Error().Err(err)
		return
	}
	imagesRepo := repo.NewImagesRepo(db)

	//init kafka
	producer := kafkawbf.NewProducer(cfg.KafkaConfig.Brokers, cfg.KafkaConfig.Topic)
	defer producer.Close()
	consumer := kafkawbf.NewConsumer(cfg.KafkaConfig.Brokers, cfg.KafkaConfig.Topic, cfg.KafkaConfig.GroupID)
	defer consumer.Close()
	strategy := retry.Strategy{Attempts: 3, Delay: 5 * time.Second, Backoff: 2}

	templatesGlob := filepath.Join("..", "web", "*")

	r := httpserver.NewRouter(imagesRepo, producer, strategy, templatesGlob)

	err = os.MkdirAll("./data/original", os.ModePerm)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Filesystem is not created")
		return
	}
	if err := os.MkdirAll("./data/processed", os.ModePerm); err != nil {
		zlog.Logger.Error().Err(err).Msg("Filesystem is not created")
		return
	}
	if err := os.MkdirAll("./data/thumbs", os.ModePerm); err != nil {
		zlog.Logger.Error().Err(err).Msg("Filesystem is not created")
		return
	}

	go worker.ProcessImages(context.Background(), imagesRepo, consumer, strategy)

	if err := r.Run(cfg.ServerConfig.Addr); err != nil {
		zlog.Logger.Error().Msg("Server is down")
		return
	}
}
