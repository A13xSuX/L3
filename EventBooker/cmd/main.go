package main

import (
	"context"
	"fmt"
	"l3/EventBooker/internal/appcfg"
	"l3/EventBooker/internal/handlers"
	"l3/EventBooker/internal/repository"
	"l3/EventBooker/internal/scheduler"
	"l3/EventBooker/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
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

	//init
	eventRepo := repository.NewEventRepository(db)
	bookingRepo := repository.NewBookingRepo(db)

	bookingService := service.NewBookingService(db, eventRepo, bookingRepo)

	expirationWorker := scheduler.NewExpirationWorker(bookingService, 5*time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go expirationWorker.Start(ctx)

	eventHandler := handlers.NewEventHandler(bookingService)

	//replace release
	router := ginext.New("debug")

	router.GET("/events", eventHandler.GetAllEvents)
	router.POST("/events", eventHandler.CreateEvent)
	router.POST("/events/:id/book", eventHandler.Book)
	router.POST("/events/:id/confirm", eventHandler.Confirm)
	router.GET("/events/:id", eventHandler.GetEventWithDetails)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		zlog.Logger.Info().Msg("Shutting down")
		cancel()

		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()

	if err = router.Run(":8080"); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
