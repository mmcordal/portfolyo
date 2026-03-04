package repository

import (
	"context"
	"database/sql"
	"errors"
	"portfolyo/internal/model"

	"github.com/uptrace/bun"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserProfile(ctx context.Context, email string) (*model.User, error)
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
	_, err := r.db.NewUpdate().
		Model(user).
		Set("deleted_at = now()").
		WherePK().
		Exec(ctx)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserProfile(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)

	err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Relation("Assets.Transactions", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("created_at DESC")
		}).
		Relation("Reminders", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("created_at DESC")
		}).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) ExistEmail(ctx context.Context, email string) (bool, error) {
	user := &model.User{}

	count, err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
