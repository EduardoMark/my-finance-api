package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/EduardoMark/my-finance-api/internal/validator"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID string, arg *CreateCategoryReq) error
	GetCategory(ctx context.Context, id string) (*db.Category, error)
	GetAllCategoriesByUserId(ctx context.Context, userId string) ([]*db.Category, error)
	Update(ctx context.Context, id string, ags UpdateCategoryReq) error
	Delete(ctx context.Context, id string) error
}

type categoryService struct {
	repo Repository
}

func NewCategoryService(repo Repository) Service {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, userID string, req *CreateCategoryReq) error {
	userUUID := uuid.MustParse(userID)

	arg := db.CreateCategoryParams{
		Name:   req.Name,
		Type:   db.TransactionType(req.Type),
		UserID: userUUID,
	}

	if err := s.repo.Create(ctx, arg); err != nil {
		return fmt.Errorf("service.create: %v", err)
	}

	return nil
}

func (s *categoryService) GetCategory(ctx context.Context, id string) (*db.Category, error) {
	idUUID := uuid.MustParse(id)

	record, err := s.repo.GetCategory(ctx, idUUID)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("service get category: %w", err)
	}

	return record, nil
}

func (s *categoryService) GetAllCategoriesByUserId(ctx context.Context, userId string) ([]*db.Category, error) {
	userUUID := uuid.MustParse(userId)

	records, err := s.repo.GetAllCategoriesByUserId(ctx, userUUID)
	if err != nil {
		if errors.Is(err, ErrCategoriesNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("service get all categories: %w", err)
	}

	return records, nil
}

func (s *categoryService) Update(ctx context.Context, id string, ags UpdateCategoryReq) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	record, err := s.repo.GetCategory(ctx, idUUID)
	if err != nil {
		return err
	}

	updateParams := db.UpdateCategoryParams{
		ID:   record.ID,
		Name: record.Name,
		Type: record.Type,
	}

	if validator.NotBlank(ags.Name) {
		updateParams.Name = ags.Name
	}

	if validator.NotBlank(ags.Type) {
		updateParams.Type = db.TransactionType(ags.Type)
	}

	if err := s.repo.Update(ctx, updateParams); err != nil {
		return err
	}

	return nil
}

func (s *categoryService) Delete(ctx context.Context, id string) error {
	idUUID := uuid.MustParse(id)

	if err := s.repo.Delete(ctx, idUUID); err != nil {
		return err
	}

	return nil
}
