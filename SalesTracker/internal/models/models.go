package models

import "time"

type Sale struct {
	ID       string
	Title    string
	Category int
	Price    float64
	Quantity int
	SaleDate time.Time
}
