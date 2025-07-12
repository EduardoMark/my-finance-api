package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, arg db.CreateUserParams) error
	GetUser(ctx context.Context, id uuid.UUID) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	GetAllUser(ctx context.Context) ([]*db.User, error)
	Update(ctx context.Context, arg db.UpdateUserParams) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *db.Queries
}

func NewUserRepository(db *db.Queries) Repository {
	return &userRepository{
		db: db,
	}
}

var ErrDuplicatedCredential = errors.New("credential already exist")
var ErrUserNotFound = errors.New("user not found")
var ErrNoUsersFound = errors.New("no users found")

func (r *userRepository) Create(ctx context.Context, arg db.CreateUserParams) error {
	var pgErr *pgconn.PgError

	err := r.db.CreateUser(ctx, arg)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicatedCredential
		}
		return err
	}
	return nil
}

func (r *userRepository) GetUser(ctx context.Context, id uuid.UUID) (*db.User, error) {
	user, err := r.db.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetAllUser(ctx context.Context) ([]*db.User, error) {
	users, err := r.db.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNoUsersFound
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, arg db.UpdateUserParams) error {
	var pgErr *pgconn.PgError

	_, err := r.GetUser(ctx, arg.ID)
	if err != nil {
		return err
	}

	err = r.db.UpdateUser(ctx, arg)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicatedCredential
		}
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.GetUser(ctx, id)
	if err != nil {
		return err
	}

	return r.db.DeleteUser(ctx, id)
}
