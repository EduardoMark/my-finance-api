package account

import (
	"context"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	Create(ctx context.Context, args db.CreateAccountParams) error
	GetAccount(ctx context.Context, id pgtype.UUID) (*db.Account, error)
	GetAccountByUserID(ctx context.Context, userID pgtype.UUID) ([]db.Account, error)
	UpdateAccount(ctx context.Context, args db.UpdateAccountParams) error
	UpdateAccountBalance(ctx context.Context, args db.UpdateAccountBalanceParams) (pgtype.Float8, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}

type accountRepository struct {
	db *db.Queries
}

func NewAccountRepo(db *db.Queries) Repository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, args db.CreateAccountParams) error {
	_, err := r.db.CreateAccount(ctx, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetAccount(ctx context.Context, id pgtype.UUID) (*db.Account, error) {
	record, err := r.db.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *accountRepository) GetAccountByUserID(ctx context.Context, userID pgtype.UUID) ([]db.Account, error) {
	records, err := r.db.GetAccountsByUserId(ctx, userID)
	if err != nil {
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

func (r *accountRepository) UpdateAccountBalance(ctx context.Context, args db.UpdateAccountBalanceParams) (pgtype.Float8, error) {
	balance, err := r.db.UpdateAccountBalance(ctx, args)
	if err != nil {
		return pgtype.Float8{}, err
	}
	return balance, nil
}

func (r *accountRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	if err := r.db.DeleteAccount(ctx, id); err != nil {
		return err
	}
	return nil
}
