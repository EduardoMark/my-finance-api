package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/EduardoMark/my-finance-api/pkg/converter"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	Create(ctx context.Context, userID string, dto AccountCreateRequest) error
}

type accountService struct {
	repo Repository
}

func NewAccountService(repo Repository) Service {
	return &accountService{repo: repo}
}

var validate = validator.New(validator.WithRequiredStructEnabled())

var ErrAccountNotFound = errors.New("account not found")
var ErrNoAccountsFound = errors.New("accounts not found")

func (s *accountService) Create(ctx context.Context, userID string, dto AccountCreateRequest) error {
	if err := validate.Struct(dto); err != nil {
		return errors.New("invalid body the field name is required")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("cannot convert userID to UUID")
	}

	params := db.CreateAccountParams{
		UserID: pgtype.UUID{Bytes: userUUID, Valid: true},
		Name:   dto.Name,
	}

	if dto.Type != "" {
		params.Type = dto.Type
	}

	if dto.Balance >= converter.Float64(params.Balance) {
		params.Balance = converter.ToFloat8(dto.Balance)
	}

	if err := s.repo.Create(ctx, params); err != nil {
		return fmt.Errorf("error on creating account: %w", err)
	}

	return nil
}
