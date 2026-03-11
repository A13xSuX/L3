package service

import (
	"context"
	"l3/EventBooker/internal/customErrs"
	"l3/EventBooker/internal/models"
	"l3/EventBooker/internal/repository"
	"time"

	"github.com/wb-go/wbf/dbpg"
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

func (s *BookingService) Book(ctx context.Context, eventID, username string) (*models.Booking, error) {
	if len(eventID) == 0 {
		return nil, customErrs.ErrInvalidEventID
	}
	if len(username) == 0 {
		return nil, customErrs.ErrInvalidUsername
	}
	//TODO я же не просто так перекидывал opts в main
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
	//TODO newSeats will do rename
	err = s.eventRepo.UpdateSeatsTx(ctx, tx, eventID, -1)
	if err != nil {
		return nil, err
	}
	booking := &models.Booking{
		EventID:   eventID,
		Username:  username,
		Status:    "pending",
		CreatedAt: time.Now(),
		ExpiredAt: ptrTime(time.Now().Add(15 * time.Minute)),
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
