package transaction

import (
	"context"
	"time"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/EduardoMark/my-finance-api/internal/validator"
)

type TransactionCreateRequest struct {
	Description string  `json:"description" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
	Date        string  `json:"date" validate:"required"` // formato: YYYY-MM-DD
	Type        string  `json:"type" validate:"required"` // "income" ou "expense"
	AccountID   string  `json:"account_id" validate:"required"`
	CategoryID  string  `json:"category_id" validate:"required"`
}

func (r *TransactionCreateRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(r.Description), "description", "this field cannot be empty")
	eval.CheckField(r.Amount > 0, "amount", "this field must be greater than 0")
	eval.CheckField(validator.NotBlank(r.Date), "date", "this field cannot be empty")
	eval.CheckField(r.Type == "income" || r.Type == "expense", "type", "this field must be 'income' or 'expense'")
	eval.CheckField(validator.NotBlank(r.AccountID), "account_id", "this field cannot be empty")
	eval.CheckField(validator.NotBlank(r.CategoryID), "category_id", "this field cannot be empty")

	return eval
}

type TransactionUpdateRequest struct {
	Description *string  `json:"description,omitempty"`
	Amount      *float64 `json:"amount,omitempty"`
	Date        *string  `json:"date,omitempty"` // formato: YYYY-MM-DD
	Type        *string  `json:"type,omitempty"` // "income" ou "expense"
	AccountID   *string  `json:"account_id,omitempty"`
	CategoryID  *string  `json:"category_id,omitempty"`
}

func (r *TransactionUpdateRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	hasAtLeastOneField := r.Description != nil || r.Amount != nil || r.Date != nil ||
		r.Type != nil || r.AccountID != nil || r.CategoryID != nil

	eval.CheckField(hasAtLeastOneField, "fields", "at least one field must be sent to update")

	if r.Amount != nil {
		eval.CheckField(*r.Amount > 0, "amount", "this field must be greater than 0")
	}

	if r.Type != nil {
		eval.CheckField(*r.Type == "income" || *r.Type == "expense", "type", "this field must be 'income' or 'expense'")
	}

	return eval
}

type TransactionResponse struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        string    `json:"date"`
	Type        string    `json:"type"`
	AccountID   string    `json:"account_id"`
	CategoryID  string    `json:"category_id"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TransactionFilters struct {
	AccountID  *string `json:"account_id,omitempty"`
	CategoryID *string `json:"category_id,omitempty"`
	Type       *string `json:"type,omitempty"`
	StartDate  *string `json:"start_date,omitempty"`
	EndDate    *string `json:"end_date,omitempty"`
}

func TransactionToResponse(t *db.Transaction) TransactionResponse {
	var dateStr string
	if t.Date.Valid {
		dateStr = t.Date.Time.Format("2006-01-02")
	}

	return TransactionResponse{
		ID:          t.ID.String(),
		Description: t.Description,
		Amount:      t.Amount,
		Date:        dateStr,
		Type:        string(t.Type),
		AccountID:   t.AccountID.String(),
		CategoryID:  t.CategoryID.String(),
		UserID:      t.UserID.String(),
		CreatedAt:   t.CreatedAt.Time,
		UpdatedAt:   t.UpdatedAt.Time,
	}
}
