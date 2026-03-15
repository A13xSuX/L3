package customErrs

import "errors"

var (
	ErrIDNotFound          = errors.New("sales id  not found")
	ErrPriceNotPositive    = errors.New("price must be positive")
	ErrQuantityNotPositive = errors.New("quantity must be positive")
	ErrTitleEmpty          = errors.New("title is empty")
	ErrCategoryEmpty       = errors.New("category is empty")
)
