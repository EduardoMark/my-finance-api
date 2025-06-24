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
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	Create(ctx context.Context, userID string, dto AccountCreateRequest) error
	GetAccount(ctx context.Context, id string) (*db.Account, error)
	GetAllAccountsByUserID(ctx context.Context, userID string) ([]db.Account, error)
	UpdateAccount(ctx context.Context, id string, args UpdateAccountReq) error
	UpdateBalance(ctx context.Context, id string, newBalance float64) (*float64, error)
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

func (s *accountService) GetAccount(ctx context.Context, id string) (*db.Account, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	record, err := s.repo.GetAccount(ctx, pgtype.UUID{Bytes: idUUID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return record, nil
}

func (s *accountService) GetAllAccountsByUserID(ctx context.Context, userID string) ([]db.Account, error) {
	idUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	records, err := s.repo.GetAccountByUserID(ctx, pgtype.UUID{Bytes: idUUID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoAccountsFound
		}
		return nil, err
	}

	return records, nil
}

func (s *accountService) UpdateAccount(ctx context.Context, id string, args UpdateAccountReq) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	record, err := s.repo.GetAccount(ctx, pgtype.UUID{Bytes: idUUID, Valid: true})
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

func (s *accountService) UpdateBalance(ctx context.Context, id string, value float64) (*float64, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	acc, err := s.repo.GetAccount(ctx, pgtype.UUID{Bytes: idUUID, Valid: true})
	if err != nil {
		return nil, err
	}

	newBalance := acc.Balance.Float64 + value

	record, err := s.repo.UpdateAccountBalance(ctx, db.UpdateAccountBalanceParams{
		ID:      pgtype.UUID{Bytes: idUUID, Valid: true},
		Balance: pgtype.Float8{Float64: newBalance, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	return &record.Float64, nil
}

func (s *accountService) Delete(ctx context.Context, id string) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, pgtype.UUID{Bytes: idUUID, Valid: true}); err != nil {
		return err
	}

	return nil
}
