package customErrs

import "errors"

var (
	ErrBookingCanceled         = errors.New("book has canceled")
	ErrInvalidBookingID        = errors.New("invalid bookingID")
	ErrInvalidEventID          = errors.New("invalid eventID")
	ErrInvalidUsername         = errors.New("invalid username")
	ErrEventNotFound           = errors.New("event not found")
	ErrNoAvailableSeats        = errors.New("no available seats")
	ErrBookingNotFound         = errors.New("booking not found")
	ErrBookingExpired          = errors.New("booking has expired")
	ErrBookingAlreadyConfirmed = errors.New("booking already confirmed")
	ErrInvalidPayment          = errors.New("invalid payment amount")
	//TODO errNotEnoughPayment
)
