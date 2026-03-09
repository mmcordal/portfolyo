package repository

import (
	"context"
	"database/sql"
	"errors"
	"portfolyo/internal/model"

	"github.com/uptrace/bun"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *model.Transaction) error
	WithTx(ctx context.Context, fn func(r UserAssetsRepository, tr TransactionRepository) error) error
	GetAllTransactionByAsset(ctx context.Context, userID int64) ([]*model.Transaction, error)
	GetAllTransaction(ctx context.Context, userID int64) ([]*model.Transaction, error)
	GetTransactionByID(ctx context.Context, txID int64) (*model.Transaction, error)
	GetTransactionByIDAndUserID(ctx context.Context, txID int64, userID int64) (*model.Transaction, error)
}

type transactionRepository struct {
	db bun.IDB
}

func NewTransactionRepository(db bun.IDB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, ta *model.Transaction) error {
	_, err := r.db.NewInsert().Model(ta).Exec(ctx)
	return err
}

func (r *transactionRepository) WithTx(ctx context.Context, fn func(UserAssetsRepository, TransactionRepository) error) error {

	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

		uarTx := &userAssetsRepository{db: tx}
		trTx := &transactionRepository{db: tx}

		return fn(uarTx, trTx)
	})
}

func (r *transactionRepository) GetAllTransactionByAsset(ctx context.Context, assetID int64) ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	err := r.db.NewSelect().
		Model(&transactions).
		Where("t.asset_id = ?", assetID).
		Relation("UserAsset.User").
		Order("t.transaction_date DESC").
		Scan(ctx)

	return transactions, err
}

func (r *transactionRepository) GetAllTransaction(ctx context.Context, userID int64) ([]*model.Transaction, error) {
	transactions := []*model.Transaction{}

	err := r.db.NewSelect().
		Model(&transactions).
		Relation("UserAsset", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("user_id = ?", userID)
		}).
		Relation("UserAsset.User").
		Order("t.transaction_date DESC").
		Scan(ctx)

	return transactions, err
}

func (r *transactionRepository) GetTransactionByID(ctx context.Context, txID int64) (*model.Transaction, error) {
	tx := new(model.Transaction)
	tx.ID = txID
	err := r.db.NewSelect().
		Model(tx).
		WherePK("id").
		Relation("UserAsset.User").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return tx, nil
}

func (r *transactionRepository) GetTransactionByIDAndUserID(ctx context.Context, txID int64, userID int64) (*model.Transaction, error) {
	tx := new(model.Transaction)
	err := r.db.NewSelect().
		Model(tx).
		Where("t.id = ?", txID).
		Relation("UserAsset", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("user_id = ?", userID)
		}).
		Relation("UserAsset.User").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return tx, nil
}
