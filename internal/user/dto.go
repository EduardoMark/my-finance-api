package user

import (
	"context"

	"github.com/EduardoMark/my-finance-api/internal/validator"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r *UserCreateRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(r.Name), "name", "this field cannot be empty")
	eval.CheckField(validator.MinChars(r.Name, 3), "name", "this field need have min 3 chars")

	eval.CheckField(validator.NotBlank(r.Email), "email", "this field cannot be empty")
	eval.CheckField(validator.Matches(r.Email, validator.EmailRX), "email", "email need as valid email")

	eval.CheckField(validator.NotBlank(r.Password), "password", "this field cannot be empty")
	eval.CheckField(validator.MinChars(r.Password, 8), "password", "this field need have min 8 chars")

	return eval
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

func (r *UserLoginRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(r.Email), "email", "this field cannot be empty")
	eval.CheckField(validator.Matches(r.Email, validator.EmailRX), "email", "this field need as valid email")

	eval.CheckField(validator.NotBlank(r.Password), "password", "this field cannot be empty")
	eval.CheckField(validator.MinChars(r.Password, 8), "password", "this field need have min 8 chars")

	return eval
}

type UserUpdateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func (r *UserUpdateRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(
		validator.NotBlank(r.Name) ||
			validator.NotBlank(r.Email) ||
			validator.NotBlank(r.Password),
		"fields", "at least one field must be sent to update",
	)

	if validator.NotBlank(r.Name) {
		eval.CheckField(validator.MinChars(r.Name, 3), "name", "this field need min 3 chars")
	}

	if validator.NotBlank(r.Email) {
		eval.CheckField(validator.Matches(r.Email, validator.EmailRX), "email", "this field need as valid email")
	}

	if validator.NotBlank(r.Password) {
		eval.CheckField(validator.MinChars(r.Password, 8), "password", "this field need have min 8 chars")
	}

	return eval
}
