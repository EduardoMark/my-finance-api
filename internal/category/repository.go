package category

import (
	"context"
	"database/sql"
	"errors"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, arg db.CreateCategoryParams) error
	GetCategory(ctx context.Context, id uuid.UUID) (*db.Category, error)
	GetAllCategoriesByUserId(ctx context.Context, userID uuid.UUID) ([]*db.Category, error)
	Update(ctx context.Context, arg db.UpdateCategoryParams) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type categoryRepository struct {
	db *db.Queries
}

func NewCategoryRepository(db *db.Queries) Repository {
	return &categoryRepository{db: db}
}

var ErrCategoryNotFound = errors.New("category not found")
var ErrCategoriesNotFound = errors.New("categories not found")

func (r *categoryRepository) Create(ctx context.Context, arg db.CreateCategoryParams) error {
	err := r.db.CreateCategory(ctx, arg)

	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) GetCategory(ctx context.Context, id uuid.UUID) (*db.Category, error) {
	record, err := r.db.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return record, nil
}

func (r *categoryRepository) GetAllCategoriesByUserId(ctx context.Context, userID uuid.UUID) ([]*db.Category, error) {
	records, err := r.db.GetAllCategoriesByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoriesNotFound
		}
		return nil, err
	}
	return records, nil
}

func (r *categoryRepository) Update(ctx context.Context, arg db.UpdateCategoryParams) error {
	if err := r.db.UpdateCategory(ctx, arg); err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return err
	}

	if err := r.db.DeleteCategory(ctx, id); err != nil {
		return err
	}
	return nil
}
