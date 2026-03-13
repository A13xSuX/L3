package producer

import (
	"l3/EventBooker/internal/customErrs"
	"time"
)

type ExpirationMessage struct {
	BookingID string    `json:"booking_id"`
	EventID   string    `json:"event_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (m ExpirationMessage) Validate() error {
	if m.BookingID == "" {
		return customErrs.ErrEmptyBookingID
	}
	if m.EventID == "" {
		return customErrs.ErrEmptyEventID
	}
	if m.ExpiresAt.IsZero() {
		return customErrs.ErrBookingExpired
	}
	return nil
}
