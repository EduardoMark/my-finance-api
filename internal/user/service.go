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
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	Create(ctx context.Context, dto UserCreateRequest) error
	GetUser(ctx context.Context, id string) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
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
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	if err := s.repo.Create(ctx, user); err != nil {
		errorMsg := fmt.Errorf("error on create user: %w", err)
		return errors.New(errorMsg.Error())
	}

	return nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*db.User, error) {
	idUUID := uuid.MustParse(id)
	pgUUID := pgtype.UUID{Bytes: idUUID, Valid: true}

	record, err := s.repo.GetUser(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrUserNotFound.Error())
		}
		return nil, fmt.Errorf("error on search user: %w", err)
	}

	return record, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	record, err := s.repo.GetUserByEmail(ctx, email)
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
	idUUID := uuid.MustParse(id)
	pgUUID := pgtype.UUID{Bytes: idUUID, Valid: true}

	record, err := s.repo.GetUser(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	now := time.Now()
	updatedParams := db.UpdateUserParams{
		ID:        pgtype.UUID{Bytes: idUUID, Valid: true},
		Name:      record.Name,
		Email:     record.Email,
		Password:  record.Password,
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
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
	idUUID := uuid.MustParse(id)
	pgUUID := pgtype.UUID{Bytes: idUUID, Valid: true}

	if err := s.repo.Delete(ctx, pgUUID); err != nil {
		return fmt.Errorf("error on delete user: %w", err)
	}

	return nil
}
