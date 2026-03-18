package customErrs

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrNotFoundID          = errors.New("id not found")
	ErrTitleEmpty          = errors.New("title is empty")
	ErrSkuEmpty            = errors.New("sku is empty")
	ErrQuantityNotPositive = errors.New("quantity cannot be negative")
)
