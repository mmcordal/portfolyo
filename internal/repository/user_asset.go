package repository

import (
	"context"
	"database/sql"
	"errors"
	"portfolyo/internal/model"

	"github.com/uptrace/bun"
)

type UserAssetsRepository interface {
	GetByUserIDAndAssetType(ctx context.Context, userID int64, assetType model.AssetType) (*model.UserAsset, error)
	UserAssetAdd(ctx context.Context, userAsset *model.UserAsset) error
	Update(ctx context.Context, userAsset *model.UserAsset) error
	GetUserAssets(ctx context.Context, userID int64) ([]*model.UserAsset, error)
}

type userAssetsRepository struct {
	db bun.IDB
}

func NewUserAssetsRepository(db bun.IDB) UserAssetsRepository {
	return &userAssetsRepository{db: db}
}

func (r *userAssetsRepository) GetByUserIDAndAssetType(ctx context.Context, userID int64, assetType model.AssetType) (*model.UserAsset, error) {
	var ua model.UserAsset

	err := r.db.NewSelect().
		Model(&ua).
		Where("user_id = ?", userID).
		Where("asset = ?", assetType).
		Limit(1).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // kayıt yok / hata değil
		}
		return nil, err
	}

	return &ua, nil
}

func (r *userAssetsRepository) UserAssetAdd(ctx context.Context, userAsset *model.UserAsset) error {
	_, err := r.db.NewInsert().Model(userAsset).Exec(ctx)
	return err
}

func (r *userAssetsRepository) Update(ctx context.Context, userAsset *model.UserAsset) error {
	_, err := r.db.NewUpdate().Model(userAsset).Where("id = ?", userAsset.ID).Exec(ctx)
	return err
}

func (r *userAssetsRepository) GetUserAssets(ctx context.Context, userID int64) ([]*model.UserAsset, error) {
	var userAssets []*model.UserAsset
	err := r.db.NewSelect().
		Model(&userAssets).
		Where("user_id = ?", userID).
		Scan(ctx)
	return userAssets, err
}
