package account

import (
	"time"
)

type AccountCreateRequest struct {
	Name    string  `json:"name" validate:"required"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
}

type AccountResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateAccountBalanceReq struct {
	Balance float64 `json:"balance" validate:"required"`
}

type UpdateAccountReq struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type UpdateAccountBalanceRes struct {
	Balance float64 `json:"balance"`
}
