package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/EduardoMark/my-finance-api/pkg/converter"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID string, dto AccountCreateRequest) error
	GetAccount(ctx context.Context, id string) (*db.Account, error)
	GetAllAccountsByUserID(ctx context.Context, userID string) ([]*db.Account, error)
	UpdateAccount(ctx context.Context, id string, args UpdateAccountReq) error
	Delete(ctx context.Context, id string) error
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

	userUUID := uuid.MustParse(userID)

	params := db.CreateAccountParams{
		UserID: userUUID,
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

func (s *accountService) GetAccount(ctx context.Context, id string) (*db.Account, error) {
	idUUID := uuid.MustParse(id)

	record, err := s.repo.GetAccount(ctx, idUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return record, nil
}

func (s *accountService) GetAllAccountsByUserID(ctx context.Context, userID string) ([]*db.Account, error) {
	idUUID := uuid.MustParse(userID)

	records, err := s.repo.GetAccountByUserID(ctx, idUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoAccountsFound
		}
		return nil, err
	}

	return records, nil
}

func (s *accountService) UpdateAccount(ctx context.Context, id string, args UpdateAccountReq) error {
	idUUID := uuid.MustParse(id)

	record, err := s.repo.GetAccount(ctx, idUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows); err != nil {
			return ErrAccountNotFound
		}
		return err
	}

	updateParams := db.UpdateAccountParams{
		ID:   record.ID,
		Name: record.Name,
		Type: record.Type,
	}

	if args.Name != "" {
		updateParams.Name = args.Name
	}

	if args.Type != "" {
		updateParams.Type = args.Type
	}

	if err := s.repo.UpdateAccount(ctx, updateParams); err != nil {
		return err
	}

	return nil
}

func (s *accountService) Delete(ctx context.Context, id string) error {
	idUUID := uuid.MustParse(id)

	if err := s.repo.Delete(ctx, idUUID); err != nil {
		return err
	}

	return nil
}
