package models

import "time"

type Event struct {
	ID              string    `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	Date            time.Time `json:"date" db:"date"`
	TotalSeats      int       `json:"totalSeats" db:"totalSeats"`
	AvailableSeats  int       `json:"availableSeats" db:"availableSeats"`
	Price           float64   `json:"price" db:"price"`
	PaymentRequired bool      `json:"paymentRequired" db:"paymentRequired"`
	CreatedAt       time.Time `json:"createdAt" db:"createdAt"`
}

type Booking struct {
	ID          string     `json:"id" db:"id"`
	EventID     string     `json:"eventID" db:"eventID"`
	Username    string     `json:"username" db:"username"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"createdAt" db:"createdAt"`
	ExpiredAt   time.Time  `json:"expiredAt" db:"expiredAt"`
	ConfirmedAt *time.Time `json:"confirmed_at" db:"confirmed_at"`
}
