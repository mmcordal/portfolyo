package repository

import (
	"context"
	"portfolyo/internal/model"

	"github.com/uptrace/bun"
)

type ReminderRepository interface {
	Create(ctx context.Context, reminde *model.Reminder) error
	GetByReminderID(ctx context.Context, id, userID int64) (*model.Reminder, error)
	Delete(ctx context.Context, reminde *model.Reminder) error
	GetAll(ctx context.Context, userID int64) ([]*model.Reminder, error)
}

type reminderRepository struct {
	db *bun.DB
}

func NewReminderRepository(db *bun.DB) ReminderRepository {
	return &reminderRepository{db: db}
}

func (r *reminderRepository) Create(ctx context.Context, reminde *model.Reminder) error {
	_, err := r.db.NewInsert().Model(reminde).Exec(ctx)
	return err
}

func (r *reminderRepository) GetByReminderID(ctx context.Context, id, userID int64) (*model.Reminder, error) {
	reminde := new(model.Reminder)
	reminde.ID = id

	err := r.db.NewSelect().
		Model(reminde).
		WherePK("id").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return reminde, nil
}

func (r *reminderRepository) Delete(ctx context.Context, reminde *model.Reminder) error {
	_, err := r.db.NewDelete().Model(reminde).WherePK().Exec(ctx)
	return err

}

func (r *reminderRepository) GetAll(ctx context.Context, userID int64) ([]*model.Reminder, error) {
	var reminde []*model.Reminder
	err := r.db.NewSelect().
		Model(&reminde).
		Where("user_id = ?", userID).
		Relation("User").
		Scan(ctx)
	return reminde, err

}
