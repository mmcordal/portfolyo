package repository

import (
	"context"
	"database/sql"
	"errors"
	"portfolyo/internal/model"

	"github.com/uptrace/bun"
)

type UserAssetsRepository interface {
	Update(ctx context.Context, ua *model.UserAsset) error
	GetUserAssets(ctx context.Context, userID int64) ([]*model.UserAsset, error)
	GetUserAssetWithTransactionByAsset(ctx context.Context, userID int64, target model.AssetType) (*model.UserAsset, error)
	FindOrCreateByUserAndAsset(ctx context.Context, userID int64, asset model.AssetType) (*model.UserAsset, error)
}

type userAssetsRepository struct {
	db bun.IDB
}

func NewUserAssetsRepository(db bun.IDB) UserAssetsRepository {
	return &userAssetsRepository{db: db}
}

func (r *userAssetsRepository) Update(ctx context.Context, ua *model.UserAsset) error {
	_, err := r.db.NewUpdate().Model(ua).Where("id = ?", ua.ID).Exec(ctx)
	return err
}

func (r *userAssetsRepository) GetUserAssets(ctx context.Context, userID int64) ([]*model.UserAsset, error) {
	var userAssets []*model.UserAsset
	err := r.db.NewSelect().
		Model(&userAssets).
		Where("user_id = ?", userID).
		Relation("User").
		Relation("Transactions", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("created_at DESC")
		}).
		Scan(ctx)
	return userAssets, err
}

func (r *userAssetsRepository) GetUserAssetWithTransactionByAsset(ctx context.Context, userID int64, target model.AssetType) (*model.UserAsset, error) {
	ua := new(model.UserAsset)
	err := r.db.NewSelect().
		Model(ua).
		Where("user_id = ?", userID).
		Where("asset = ?", target).
		Relation("User").
		Relation("Transactions", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("created_at DESC")
		}).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ua, err
}

func (r *userAssetsRepository) FindOrCreateByUserAndAsset(ctx context.Context, userID int64, asset model.AssetType) (*model.UserAsset, error) {

	ua := &model.UserAsset{
		UserID: userID,
		Asset:  asset,
		Amount: 0,
	}

	_, err := r.db.NewInsert().
		Model(ua).
		On("CONFLICT (user_id, asset) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	err = r.db.NewSelect().
		Model(ua).
		Where("user_id = ?", userID).
		Where("asset = ?", asset).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return ua, nil
}
