package user

import (
	"context"

	"github.com/EduardoMark/my-finance-api/internal/store/pgstore/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	Create(ctx context.Context, arg db.CreateUserParams) error
	GetUser(ctx context.Context, id pgtype.UUID) (*db.User, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	GetAllUser(ctx context.Context) ([]db.User, error)
	Update(ctx context.Context, arg db.UpdateUserParams) error
	Delete(ctx context.Context, id pgtype.UUID) error
}

type userRepository struct {
	db *db.Queries
}

func NewUserRepository(db *db.Queries) Repository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, arg db.CreateUserParams) error {
	return r.db.CreateUser(ctx, arg)
}

func (r *userRepository) GetUser(ctx context.Context, id pgtype.UUID) (*db.User, error) {
	user, err := r.db.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAllUser(ctx context.Context) ([]db.User, error) {
	return r.db.GetAllUsers(ctx)
}

func (r *userRepository) Update(ctx context.Context, arg db.UpdateUserParams) error {
	return r.db.UpdateUser(ctx, arg)
}

func (r *userRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.db.DeleteUser(ctx, id)
}
