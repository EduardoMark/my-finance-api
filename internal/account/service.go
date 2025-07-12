package account

import (
	"context"
	"fmt"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/EduardoMark/my-finance-api/pkg/converter"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID string, dto AccountCreateRequest) error
	GetAccount(ctx context.Context, id string) (*db.Account, error)
	GetAllAccountsByUserID(ctx context.Context, userID string) ([]*db.Account, error)
	UpdateAccount(ctx context.Context, id string, args AccountUpdateAccountReq) error
	Delete(ctx context.Context, id string) error
}

type accountService struct {
	repo Repository
}

func NewAccountService(repo Repository) Service {
	return &accountService{repo: repo}
}

func (s *accountService) Create(ctx context.Context, userID string, dto AccountCreateRequest) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	params := db.CreateAccountParams{
		UserID: userUUID,
		Name:   dto.Name,
	}

	if dto.Type != "" {
		params.Type = dto.Type
	}

	// Define balance como 0 por padrão se não informado
	if dto.Balance != nil {
		params.Balance = converter.ToFloat8(*dto.Balance)
	} else {
		params.Balance = converter.ToFloat8(0)
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

	record, err := s.repo.GetAccount(ctx, idUUID)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *accountService) GetAllAccountsByUserID(ctx context.Context, userID string) ([]*db.Account, error) {
	idUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	records, err := s.repo.GetAccountByUserID(ctx, idUUID)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *accountService) UpdateAccount(ctx context.Context, id string, args AccountUpdateAccountReq) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	record, err := s.repo.GetAccount(ctx, idUUID)
	if err != nil {
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
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, idUUID); err != nil {
		return err
	}

	return nil
}
