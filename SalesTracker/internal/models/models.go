package models

import (
	"l3/SalesTracker/internal/customErrs"
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

type AnalyticsResponse struct {
	Sum    float64
	Avg    float64
	Count  int
	Median float64
	P90    float64
}

type AnalyticsFilter struct {
	From     time.Time
	To       time.Time
	Category string
	Title    string
}

func (s *Sale) Validation() error {
	if s.Price <= 0 {
		return customErrs.ErrPriceNotPositive
	}
	if s.Quantity <= 0 {
		return customErrs.ErrQuantityNotPositive
	}
	if s.Title == "" {
		return customErrs.ErrTitleEmpty
	}
	if s.Category == "" {
		return customErrs.ErrCategoryEmpty
	}
	return nil
}
