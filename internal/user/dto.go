package user

import "github.com/jackc/pgx/v5/pgtype"

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        pgtype.UUID        `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}
