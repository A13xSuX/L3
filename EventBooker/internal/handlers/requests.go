package handlers

import "time"

type CreateEventRequest struct {
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description"`
	Date            time.Time `json:"date" binding:"required"`
	TotalSeats      int       `json:"totalSeats" binding:"required,gt=0"`
	Price           float64   `json:"price"`
	PaymentRequired bool      `json:"paymentRequired"`
}

type BookRequest struct {
	Username string `json:"username"`
}
