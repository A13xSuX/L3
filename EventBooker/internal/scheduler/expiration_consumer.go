package scheduler

import (
	"context"
	"encoding/json"
	"l3/EventBooker/internal/producer"
	"l3/EventBooker/internal/repository"
	"time"

	seg "github.com/segmentio/kafka-go"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
)

type ExpirationConsumer struct {
	consumer    *kafka.Consumer
	db          *dbpg.DB
	bookingRepo *repository.BookingRepository
	eventRepo   *repository.EventRepository
	retryStrat  retry.Strategy
	msgChan     chan seg.Message
}

func NewExpirationConsumer(brokers []string, topic string, groupID string, db *dbpg.DB, eventRepo *repository.EventRepository,
	bookingRepo *repository.BookingRepository) *ExpirationConsumer {
	consumer := kafka.NewConsumer(brokers, topic, groupID)

	msgChan := make(chan seg.Message, 50)

	retryStrat := retry.Strategy{
		Attempts: 3,
		Delay:    100 * time.Millisecond,
		Backoff:  2.0,
	}

	return &ExpirationConsumer{
		consumer:    consumer,
		db:          db,
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
		retryStrat:  retryStrat,
		msgChan:     msgChan,
	}
}

func (c *ExpirationConsumer) Start(ctx context.Context) {
	c.consumer.StartConsuming(ctx, c.msgChan, c.retryStrat)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-c.msgChan:
				if !ok {
					return
				}
				c.processMessage(ctx, msg)
			}
		}
	}()
}

func (c *ExpirationConsumer) processMessage(ctx context.Context, msg seg.Message) {
	var expMsg producer.ExpirationMessage

	if err := json.Unmarshal(msg.Value, &expMsg); err != nil {
		c.consumer.Commit(ctx, msg)
		return
	}

	if err := expMsg.Validate(); err != nil {
		c.consumer.Commit(ctx, msg)
		return
	}

	if time.Now().Before(expMsg.ExpiresAt) {
		return
	}

	if err := c.cancelExpirationBooking(ctx, expMsg); err != nil {
		return
	}
	c.consumer.Commit(ctx, msg)
}

func (c *ExpirationConsumer) cancelExpirationBooking(ctx context.Context, expMsg producer.ExpirationMessage) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	booking, err := c.bookingRepo.GetByIDForUpdateTx(ctx, tx, expMsg.BookingID)
	if err != nil {
		return err
	}
	if booking == nil {
		// Бронь уже не существует - считаем успехом
		return tx.Commit()
	}
	if booking.Status != "pending" {
		// Уже подтверждена или отменена - ничего не делаем
		return tx.Commit()
	}

	if err := c.eventRepo.UpdateSeatsTx(ctx, tx, booking.EventID, 1); err != nil {
		return err
	}

	if err := c.bookingRepo.UpdateStatusTx(ctx, tx, booking.ID, "canceled", nil); err != nil {
		return err
	}
	return tx.Commit()
}

func (c *ExpirationConsumer) Close() error {
	close(c.msgChan)

	if err := c.consumer.Close(); err != nil {
		return err
	}
	return nil
}
