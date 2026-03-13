package scheduler

import (
	"context"
	"l3/EventBooker/internal/service"
	"time"

	"github.com/wb-go/wbf/zlog"
)

type ExpirationWorker struct {
	bookingService *service.BookingService
	interval       time.Duration
}

func NewExpirationWorker(bookingService *service.BookingService, interval time.Duration) *ExpirationWorker {
	return &ExpirationWorker{
		bookingService: bookingService,
		interval:       interval,
	}
}

func (w *ExpirationWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			zlog.Logger.Info().Msg("expiration worker stopped")
			return
		case <-ticker.C:
			zlog.Logger.Info().Msg("expiration worker tick")

			if err := w.bookingService.CancelExpiredBookings(ctx); err != nil {
				zlog.Logger.Error().Err(err).Msg("failed to cancel expired bookings")
			}
		}
	}
}
