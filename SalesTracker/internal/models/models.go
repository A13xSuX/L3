package models

import (
	"errors"
	"time"
)

type Sale struct {
	ID       string
	Title    string
	Category string
	Price    float64
	Quantity int
	SaleDate time.Time
}

func (s *Sale) Validation() error {
	if s.Price <= 0 {
		return errors.New("price must be positive")
	}
	if s.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if s.Title == "" {
		return errors.New("title is empty")
	}
	if s.Category == "" {
		return errors.New("category is empty")
	}
	return nil
}
