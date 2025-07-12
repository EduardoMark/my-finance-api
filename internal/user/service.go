package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/EduardoMark/my-finance-api/internal/validator"
	"github.com/EduardoMark/my-finance-api/pkg/hash"
	"github.com/EduardoMark/my-finance-api/pkg/token"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, dto UserCreateRequest) error
	GetUser(ctx context.Context, id string) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	GetAllUsers(ctx context.Context) ([]*db.User, error)
	Update(ctx context.Context, id string, arg UserUpdateRequest) error
	Delete(ctx context.Context, id string) error
	Login(ctx context.Context, tm *token.TokenManager, dto UserLoginRequest) (string, error)
}

type userService struct {
	repo Repository
}

func NewUserService(repo Repository) Service {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, dto UserCreateRequest) error {
	password, err := hash.HashPassword(dto.Password)
	if err != nil {
		return err
	}

	user := db.CreateUserParams{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: password,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, ErrDuplicatedCredential) {
			return ErrDuplicatedCredential
		}

		return fmt.Errorf("service create: %w", err)
	}

	return nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*db.User, error) {
	idUUID := uuid.MustParse(id)

	record, err := s.repo.GetUser(ctx, idUUID)
	if err != nil {
		return nil, fmt.Errorf("error on search user: %w", err)
	}

	return record, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	record, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error on search user: %w", err)
	}

	return record, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*db.User, error) {
	records, err := s.repo.GetAllUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("error on get all users: %w", err)
	}

	return records, nil
}

func (s *userService) Update(ctx context.Context, id string, arg UserUpdateRequest) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	record, err := s.repo.GetUser(ctx, idUUID)
	if err != nil {
		return err
	}

	updatedParams := db.UpdateUserParams{
		ID:       idUUID,
		Name:     record.Name,
		Email:    record.Email,
		Password: record.Password,
	}

	if validator.NotBlank(arg.Name) {
		updatedParams.Name = arg.Name
	}

	if validator.NotBlank(arg.Email) {
		updatedParams.Email = arg.Email
	}

	if validator.NotBlank(arg.Password) {
		hashPassword, err := hash.HashPassword(arg.Password)
		if err != nil {
			return err
		}
		updatedParams.Password = hashPassword
	}

	if err := s.repo.Update(ctx, updatedParams); err != nil {
		return err
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	idUUID := uuid.MustParse(id)

	if err := s.repo.Delete(ctx, idUUID); err != nil {
		return fmt.Errorf("error on delete user: %w", err)
	}

	return nil
}

func (s *userService) Login(ctx context.Context, tm *token.TokenManager, dto UserLoginRequest) (string, error) {
	record, err := s.repo.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("error on search user: %w", err)
	}

	if err := hash.ComparePassword(dto.Password, record.Password); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := tm.GenerateToken(record.ID.String(), record.Name)
	if err != nil {
		return "", err
	}

	return token, err
}
