package models

import (
	"encoding/json"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}
type CreateItemRequest struct {
	Title    string `json:"title"`
	Sku      string `json:"sku"`
	Quantity int64  `json:"quantity"`
}
type UpdateItemRequest struct {
	Title    string `json:"title"`
	Sku      string `json:"sku"`
	Quantity int64  `json:"quantity"`
}

type ItemResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Sku       string    `json:"sku"`
	Quantity  int64     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HistoryItemResponse struct {
	ID                int64           `json:"id"`
	ItemID            int64           `json:"item_id"`
	Action            string          `json:"action"`
	OldData           json.RawMessage `json:"old_data"`
	NewData           json.RawMessage `json:"new_data"`
	ChangedByUsername string          `json:"changed_by_username"`
	ChangedByRole     string          `json:"changed_by_role"`
	ChangedAt         time.Time       `json:"changed_at"`
}
