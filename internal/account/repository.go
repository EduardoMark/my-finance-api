package account

import (
	"context"
	"database/sql"
	"errors"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, args db.CreateAccountParams) error
	GetAccount(ctx context.Context, id uuid.UUID) (*db.Account, error)
	GetAccountByUserID(ctx context.Context, userID uuid.UUID) ([]*db.Account, error)
	UpdateAccount(ctx context.Context, args db.UpdateAccountParams) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type accountRepository struct {
	db *db.Queries
}

func NewAccountRepo(db *db.Queries) Repository {
	return &accountRepository{db: db}
}

var ErrAccountNotFound = errors.New("account not found")
var ErrNoAccountsFound = errors.New("accounts not found")

func (r *accountRepository) Create(ctx context.Context, args db.CreateAccountParams) error {
	err := r.db.CreateAccount(ctx, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetAccount(ctx context.Context, id uuid.UUID) (*db.Account, error) {
	record, err := r.db.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return record, nil
}

func (r *accountRepository) GetAccountByUserID(ctx context.Context, userID uuid.UUID) ([]*db.Account, error) {
	records, err := r.db.GetAccountsByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return records, nil
}

func (r *accountRepository) UpdateAccount(ctx context.Context, args db.UpdateAccountParams) error {
	if err := r.db.UpdateAccount(ctx, args); err != nil {
		return err
	}
	return nil
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.GetAccount(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.DeleteAccount(ctx, id); err != nil {
		return err
	}
	return nil
}
