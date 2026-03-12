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

func (s *BookingService) Confirm(ctx context.Context, bookingID string) error {
	if len(bookingID) == 0 {
		return customErrs.ErrInvalidBookingID
	}
	//TODO options
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	//TODO
	booking, err := s.bookingRepo.GetByIDForUpdateTx(ctx, tx, bookingID)

	if booking == nil {
		return customErrs.ErrBookingNotFound
	}
	if booking.Status == "confirmed" {
		return customErrs.ErrBookingAlreadyConfirmed
	}
	if booking.Status == "canceled" {
		return customErrs.ErrBookingCanceled
	}
	//TODO is it correct logic?
	if booking.ExpiredAt != nil && booking.ExpiredAt.Before(time.Now()) {
		return customErrs.ErrBookingExpired
	}
	//TODO
	//event, err := s.eventRepo.GetByID(ctx, booking.EventID)
	//if err != nil {
	//	return err
	//}
	////check on amount of cost(mb)
	//if event == nil {
	//	return customErrs.ErrEventNotFound
	//}
	err = s.bookingRepo.UpdateStatusTx(ctx, tx, bookingID, "confirmed", ptrTime(time.Now()))
	booking.ExpiredAt = nil
	//TODO continue here

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
		//TODO прочекать та ли ошибка
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
