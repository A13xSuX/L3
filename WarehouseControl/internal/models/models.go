package models

import (
	"encoding/json"
	"time"
)

type Item struct {
	ID        int64
	Title     string
	Sku       string
	Quantity  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}

type Audit struct {
	ID                int64
	ItemID            int64
	Action            string
	OldData           json.RawMessage
	NewData           json.RawMessage
	ChangedByUserID   *int64
	ChangedByUsername string
	ChangedByRole     string
	ChangedAt         time.Time
}
