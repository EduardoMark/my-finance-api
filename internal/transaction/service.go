package transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	Create(ctx context.Context, userID string, dto TransactionCreateRequest) error
	GetTransaction(ctx context.Context, id string) (*TransactionResponse, error)
	GetAllTransactions(ctx context.Context, userID string, filters *TransactionFilters) ([]*TransactionResponse, error)
	UpdateTransaction(ctx context.Context, id string, dto TransactionUpdateRequest) error
	DeleteTransaction(ctx context.Context, id string) error
}

type transactionService struct {
	repo Repository
}

func NewTransactionService(repo Repository) Service {
	return &transactionService{
		repo: repo,
	}
}

func (s *transactionService) Create(ctx context.Context, userID string, dto TransactionCreateRequest) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	accountUUID, err := uuid.Parse(dto.AccountID)
	if err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	categoryUUID, err := uuid.Parse(dto.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}

	date, err := time.Parse("2006-01-02", dto.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	params := db.CreateTransactionParams{
		Description: dto.Description,
		Amount:      dto.Amount,
		Date:        pgtype.Date{Time: date, Valid: true},
		Type:        db.TransactionType(dto.Type),
		AccountID:   accountUUID,
		CategoryID:  categoryUUID,
		UserID:      userUUID,
	}

	if err := s.repo.Create(ctx, params); err != nil {
		return fmt.Errorf("service create transaction: %w", err)
	}

	return nil
}

func (s *transactionService) GetTransaction(ctx context.Context, id string) (*TransactionResponse, error) {
	transactionUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	transaction, err := s.repo.GetTransaction(ctx, transactionUUID)
	if err != nil {
		return nil, fmt.Errorf("service get transaction: %w", err)
	}

	response := TransactionToResponse(transaction)
	return &response, nil
}

func (s *transactionService) GetAllTransactions(ctx context.Context, userID string, filters *TransactionFilters) ([]*TransactionResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var transactions []*db.Transaction

	if filters != nil && filters.AccountID != nil {
		accountUUID, parseErr := uuid.Parse(*filters.AccountID)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid account ID: %w", parseErr)
		}
		transactions, err = s.repo.GetAllTransasctionsByAccount(ctx, accountUUID)
		if err != nil {
			return nil, fmt.Errorf("service get all transactions: %w", err)
		}
		return s.processTransactions(transactions, filters), nil
	}

	if filters != nil && filters.CategoryID != nil {
		categoryUUID, parseErr := uuid.Parse(*filters.CategoryID)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid category ID: %w", parseErr)
		}
		transactions, err = s.repo.GetAllTransasctionsByCategory(ctx, categoryUUID)
		if err != nil {
			return nil, fmt.Errorf("service get all transactions: %w", err)
		}
		return s.processTransactions(transactions, filters), nil
	}

	transactions, err = s.repo.GetAllTransaction(ctx, userUUID)

	if err != nil {
		return nil, fmt.Errorf("service get all transactions: %w", err)
	}

	return s.processTransactions(transactions, filters), nil
}

func (s *transactionService) processTransactions(transactions []*db.Transaction, filters *TransactionFilters) []*TransactionResponse {
	filteredTransactions := s.applyFilters(transactions, filters)

	var responses []*TransactionResponse
	for _, transaction := range filteredTransactions {
		response := TransactionToResponse(transaction)
		responses = append(responses, &response)
	}

	return responses
}

func (s *transactionService) UpdateTransaction(ctx context.Context, id string, dto TransactionUpdateRequest) error {
	transactionUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %w", err)
	}

	existing, err := s.repo.GetTransaction(ctx, transactionUUID)
	if err != nil {
		return fmt.Errorf("service update transaction: %w", err)
	}

	params := db.UpdateTransactionParams{
		ID:          transactionUUID,
		Description: existing.Description,
		Amount:      existing.Amount,
		Date:        existing.Date,
		Type:        existing.Type,
	}

	if dto.Description != nil {
		params.Description = *dto.Description
	}

	if dto.Amount != nil {
		params.Amount = *dto.Amount
	}

	if dto.Date != nil {
		date, err := time.Parse("2006-01-02", *dto.Date)
		if err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}
		params.Date = pgtype.Date{Time: date, Valid: true}
	}

	if dto.Type != nil {
		params.Type = db.TransactionType(*dto.Type)
	}

	if err := s.repo.Update(ctx, params); err != nil {
		return fmt.Errorf("service update transaction: %w", err)
	}

	return nil
}

func (s *transactionService) DeleteTransaction(ctx context.Context, id string) error {
	transactionUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %w", err)
	}

	if err := s.repo.Delete(ctx, transactionUUID); err != nil {
		return fmt.Errorf("service delete transaction: %w", err)
	}

	return nil
}

func (s *transactionService) applyFilters(transactions []*db.Transaction, filters *TransactionFilters) []*db.Transaction {
	if filters == nil {
		return transactions
	}

	var filtered []*db.Transaction

	for _, transaction := range transactions {
		if filters.Type != nil && string(transaction.Type) != *filters.Type {
			continue
		}

		if filters.StartDate != nil {
			startDate, err := time.Parse("2006-01-02", *filters.StartDate)
			if err == nil && transaction.Date.Valid && transaction.Date.Time.Before(startDate) {
				continue
			}
		}

		if filters.EndDate != nil {
			endDate, err := time.Parse("2006-01-02", *filters.EndDate)
			if err == nil && transaction.Date.Valid && transaction.Date.Time.After(endDate) {
				continue
			}
		}

		filtered = append(filtered, transaction)
	}

	return filtered
}
