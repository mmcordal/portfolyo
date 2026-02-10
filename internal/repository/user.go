package repository

import (
	"context"
	"portfolyo/internal/model"

	"github.com/uptrace/bun"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	ExistEmail(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.db.NewUpdate().Model(user).Where("id = ?", user.ID).Exec(ctx)
	return err
}

func (r *userRepository) Delete(ctx context.Context, user *model.User) error {
	_, err := r.db.NewDelete().Model(user).WherePK().Exec(ctx)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) ExistEmail(ctx context.Context, email string) (bool, error) {
	var user *model.User

	count, err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
