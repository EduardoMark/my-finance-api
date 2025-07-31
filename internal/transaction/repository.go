package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, args db.CreateTransactionParams) error
	GetTransaction(ctx context.Context, id uuid.UUID) (*db.Transaction, error)
	GetAllTransaction(ctx context.Context, userID uuid.UUID) ([]*db.Transaction, error)
	GetAllTransasctionsByAccount(ctx context.Context, accountID uuid.UUID) ([]*db.Transaction, error)
	GetAllTransasctionsByCategory(ctx context.Context, categoryID uuid.UUID) ([]*db.Transaction, error)
	Update(ctx context.Context, args db.UpdateTransactionParams) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type transactionRepository struct {
	db *db.Queries
}

func NewTransactionRepo(db *db.Queries) Repository {
	return &transactionRepository{
		db: db,
	}
}

var ErrTransactionNotFound = errors.New("transaction not found")

func (r *transactionRepository) Create(ctx context.Context, args db.CreateTransactionParams) error {
	if err := r.db.CreateTransaction(ctx, args); err != nil {
		return fmt.Errorf("repository create: %w", err)
	}

	return nil
}

func (r *transactionRepository) GetTransaction(ctx context.Context, id uuid.UUID) (*db.Transaction, error) {
	record, err := r.db.GetTrasaction(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTransactionNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("repository getTransaction: %w", err)
	}

	return record, nil
}

func (r *transactionRepository) GetAllTransaction(ctx context.Context, userID uuid.UUID) ([]*db.Transaction, error) {
	records, err := r.db.GetAllTransactions(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTransactionNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("repository getTransaction: %w", err)
	}

	return records, nil
}

func (r *transactionRepository) GetAllTransasctionsByAccount(ctx context.Context, accountID uuid.UUID) ([]*db.Transaction, error) {
	records, err := r.db.GetAllTransactionsByAccount(ctx, accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTransactionNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("repository getTransaction: %w", err)
	}

	return records, nil
}

func (r *transactionRepository) GetAllTransasctionsByCategory(ctx context.Context, categoryID uuid.UUID) ([]*db.Transaction, error) {
	records, err := r.db.GetAllTransactionsByCategory(ctx, categoryID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTransactionNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("repository getTransaction: %w", err)
	}

	return records, nil
}

func (r *transactionRepository) Update(ctx context.Context, args db.UpdateTransactionParams) error {
	_, err := r.GetTransaction(ctx, args.ID)
	if err != nil {
		return err
	}

	if err := r.db.UpdateTransaction(ctx, args); err != nil {
		return fmt.Errorf("repository update: %w", err)
	}

	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.GetTransaction(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.DeleteTransaction(ctx, id); err != nil {
		return fmt.Errorf("repository delete: %w", err)
	}

	return nil
}
