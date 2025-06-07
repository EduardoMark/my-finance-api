package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/EduardoMark/my-finance-api/internal/db"
	"github.com/EduardoMark/my-finance-api/pkg/hash"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, dto UserCreateRequest) error
	GetUser(ctx context.Context, id string) (*db.User, error)
	GetAllUsers(ctx context.Context) ([]db.User, error)
	Update(ctx context.Context, id string, arg UserUpdateRequest) error
	Delete(ctx context.Context, id string) error
}

type userService struct {
	repo Repository
}

func NewUserService(repo Repository) Service {
	return &userService{repo: repo}
}

var validate = validator.New(validator.WithRequiredStructEnabled())

var ErrUserNotFound = errors.New("user not found")
var ErrNoUsersFound = errors.New("users not found")

func (s *userService) Create(ctx context.Context, dto UserCreateRequest) error {
	if err := validate.Struct(dto); err != nil {
		return errors.New("invalid body all fields required")
	}

	password, err := hash.HashPassword(dto.Password)
	if err != nil {
		return err
	}

	now := time.Now()
	user := db.CreateUserParams{
		Name:      dto.Name,
		Email:     dto.Email,
		Password:  password,
		CreatedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}

	if err := s.repo.Create(ctx, user); err != nil {
		errorMsg := fmt.Errorf("error on create user: %w", err)
		return errors.New(errorMsg.Error())
	}

	return nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*db.User, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("error on parse id to uuid: %w", err)
	}

	record, err := s.repo.GetUser(ctx, idUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrUserNotFound.Error())
		}
		return nil, fmt.Errorf("error on search user: %w", err)
	}

	return record, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]db.User, error) {
	records, err := s.repo.GetAllUser(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrNoUsersFound.Error())
		}
	}

	return records, nil
}

func (s *userService) Update(ctx context.Context, id string, arg UserUpdateRequest) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("error parsing id to UUID: %w", err)
	}

	record, err := s.repo.GetUser(ctx, idUUID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	now := time.Now()
	updatedParams := db.UpdateUserParams{
		ID:        idUUID,
		Name:      record.Name,
		Email:     record.Email,
		Password:  record.Password,
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}

	if arg.Name != "" {
		updatedParams.Name = arg.Name
	}

	if arg.Email != "" {
		updatedParams.Email = arg.Email
	}

	if arg.Password != "" {
		hashPassword, err := hash.HashPassword(arg.Password)
		if err != nil {
			return err
		}
		updatedParams.Password = hashPassword
	}

	if err := s.repo.Update(ctx, updatedParams); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("error on parse id to uuid: %w", err)
	}

	if err := s.repo.Delete(ctx, idUUID); err != nil {
		return fmt.Errorf("error on delete user: %w", err)
	}

	return nil
}
