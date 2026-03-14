package service

import (
	"context"
	"l3/EventBooker/internal/customErrs"
	"l3/EventBooker/internal/models"
	"l3/EventBooker/internal/repository"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type BookingService struct {
	db          *dbpg.DB
	eventRepo   *repository.EventRepository
	bookingRepo *repository.BookingRepository
}

func NewBookingService(db *dbpg.DB, eventRepo *repository.EventRepository, bookingRepo *repository.BookingRepository) *BookingService {
	return &BookingService{
		db:          db,
		eventRepo:   eventRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *BookingService) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	return s.eventRepo.GetAll(ctx)
}

func (s *BookingService) CreateEvent(ctx context.Context, event *models.Event) error {
	return s.eventRepo.Create(ctx, event)
}

func (s *BookingService) Book(ctx context.Context, eventID, username string) (*models.Booking, error) {
	if len(eventID) == 0 {
		return nil, customErrs.ErrInvalidEventID
	}
	if len(username) == 0 {
		return nil, customErrs.ErrInvalidUsername
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	event, err := s.eventRepo.GetByIDForUpdateTx(ctx, tx, eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, customErrs.ErrEventNotFound
	}
	if event.AvailableSeats <= 0 {
		return nil, customErrs.ErrNoAvailableSeats
	}
	//уменьшаем
	err = s.eventRepo.UpdateSeatsTx(ctx, tx, eventID, -1)
	if err != nil {
		return nil, err
	}
	booking := &models.Booking{
		EventID:   eventID,
		Username:  username,
		Status:    "pending",
		CreatedAt: time.Now().UTC(),
		ExpiredAt: ptrTime(time.Now().Add(1 * time.Minute)),
	}
	if !event.PaymentRequired {
		booking.Status = "confirmed"
		booking.ExpiredAt = nil
		booking.ConfirmedAt = ptrTime(time.Now())
	}

	//do book
	err = s.bookingRepo.CreateTx(ctx, tx, booking)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) Confirm(ctx context.Context, bookingID string) error {
	if len(bookingID) == 0 {
		return customErrs.ErrInvalidBookingID
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	booking, err := s.bookingRepo.GetByIDForUpdateTx(ctx, tx, bookingID)
	if err != nil {
		return err
	}
	if booking == nil {
		return customErrs.ErrBookingNotFound
	}
	if booking.Status == "confirmed" {
		return customErrs.ErrBookingAlreadyConfirmed
	}
	if booking.Status == "canceled" {
		return customErrs.ErrBookingCanceled
	}

	if booking.ExpiredAt != nil && booking.ExpiredAt.Before(time.Now().UTC()) {
		return customErrs.ErrBookingExpired
	}

	err = s.bookingRepo.UpdateStatusTx(ctx, tx, bookingID, "confirmed", ptrTime(time.Now().UTC()))
	booking.ExpiredAt = nil

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *BookingService) GetEventWithDetails(ctx context.Context, eventID string) (*models.EventsWithDetails, error) {
	if len(eventID) == 0 {
		return nil, customErrs.ErrInvalidEventID
	}

	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, customErrs.ErrEventNotFound
	}
	bookings, err := s.bookingRepo.GetByEventID(ctx, event.ID)
	if err != nil {
		return nil, err
	}
	var totalBooked int
	for _, b := range bookings {
		if b.Status == "pending" || b.Status == "confirmed" {
			totalBooked++
		}
	}
	details := &models.EventsWithDetails{
		Event:       event,
		Bookings:    bookings,
		FreeSeats:   event.AvailableSeats,
		TotalBooked: totalBooked,
	}
	return details, nil
}

func (s *BookingService) CancelExpiredBookings(ctx context.Context) error {
	expiredBookings, err := s.bookingRepo.GetExpiredBookings(ctx)
	if err != nil {
		return err
	}
	zlog.Logger.Info().Int("expired_count", len(expiredBookings)).Msg("checked expired bookings")
	if len(expiredBookings) == 0 {
		return nil
	}
	for _, b := range expiredBookings {
		zlog.Logger.Info().Str("booking_id", b.ID).Msg("processing expired booking")
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		func() {
			defer tx.Rollback()
			booking, err := s.bookingRepo.GetByIDForUpdateTx(ctx, tx, b.ID)
			if err != nil {
				zlog.Logger.Error().Err(err).Msg("failed to lock booking")
				return
			}
			if booking == nil {
				zlog.Logger.Info().Msg("booking not found")
				return
			}
			if booking.Status != "pending" {
				zlog.Logger.Info().Str("status", booking.Status).Msg("booking is not pending")
				return
			}
			if booking.ExpiredAt == nil || booking.ExpiredAt.After(time.Now().UTC()) {
				zlog.Logger.Info().Msg("booking is not expired yet")
				return
			}

			event, err := s.eventRepo.GetByIDForUpdateTx(ctx, tx, booking.EventID)
			if err != nil {
				zlog.Logger.Error().Err(err).Msg("failed to lock event")
				return
			}
			if event == nil {
				zlog.Logger.Info().Msg("event not found")
				return
			}
			if err = s.eventRepo.UpdateSeatsTx(ctx, tx, booking.EventID, 1); err != nil {
				zlog.Logger.Error().Err(err).Msg("failed to update seats")
				return
			}
			if err = s.bookingRepo.UpdateStatusTx(ctx, tx, booking.ID, "canceled", nil); err != nil {
				zlog.Logger.Error().Err(err).Msg("failed to update booking status")
				return
			}
			if err = tx.Commit(); err != nil {
				zlog.Logger.Error().Err(err).Msg("failed to commit tx")
				return
			}
			zlog.Logger.Info().Str("booking_id", booking.ID).Msg("booking canceled")
		}()
	}
	return nil
}
