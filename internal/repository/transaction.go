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
	GetTransaction(ctx context.Context, userID int64) ([]*model.Transaction, error)
	GetTransactionByID(ctx context.Context, userID, txID int64) (*model.Transaction, error)
}

type transactionRepository struct {
	db bun.IDB
}

//  ROUTER A EKLEMEYİ UNUTMA!!!!

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

func (r *transactionRepository) GetTransaction(ctx context.Context, userID int64) ([]*model.Transaction, error) {
	var transaction []*model.Transaction
	err := r.db.NewSelect().
		Model(&transaction).
		Where("user_id = ?", userID).
		Scan(ctx)
	return transaction, err
}

func (r *transactionRepository) GetTransactionByID(ctx context.Context, userID, txID int64) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.db.NewSelect().
		Model(&tx).
		Where("user_id = ?", userID).
		Where("id = ?", txID).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}
