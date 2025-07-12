package account

import (
	"context"
	"time"

	"github.com/EduardoMark/my-finance-api/internal/validator"
)

type AccountCreateRequest struct {
	Name    string   `json:"name" validate:"required"`
	Type    string   `json:"type"`
	Balance *float64 `json:"balance"`
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

func (r *AccountCreateRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(r.Name), "name", "this field cannot be empty")
	eval.CheckField(validator.NotBlank(r.Type), "type", "this field cannot be empty")

	if r.Balance != nil {
		eval.CheckField(validator.CheckBalance(*r.Balance), "balance", "this field must be bigger than or equal to 0")
	}

	return eval
}

type AccountUpdateAccountReq struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (r *AccountUpdateAccountReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(
		validator.NotBlank(r.Name) ||
			validator.NotBlank(r.Type),
		"fields", "at least one field must be sent to update",
	)

	return eval
}
