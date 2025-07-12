package category

import (
	"context"
	"time"

	"github.com/EduardoMark/my-finance-api/internal/validator"
)

type CreateCategoryReq struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CategoryRes struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *CreateCategoryReq) Valid(context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(r.Name), "name", "this field cannot must be empty")
	eval.CheckField(validator.MinChars(r.Name, 3), "name", "this field need have min 3 chars")

	eval.CheckField(validator.NotBlank(r.Type), "type", "this field cannot must be empty")
	eval.CheckField(validator.TransactionType(r.Type), "type", "the type must be income or expense")

	return eval
}

type UpdateCategoryReq struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (r *UpdateCategoryReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(
		validator.NotBlank(r.Name) ||
			validator.NotBlank(r.Type),
		"fields", "at least one field must be sent to update",
	)

	if validator.NotBlank(r.Type) {
		eval.CheckField(validator.TransactionType(r.Type), "type", "the type must be income or expense")
	}

	return eval
}
