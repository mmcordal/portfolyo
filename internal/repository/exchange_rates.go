package repository

import (
	"context"
	"portfolyo/internal/model"

	"github.com/uptrace/bun"
)

type ExchangeRatesRepository interface {
	Create(ctx context.Context, rate *model.ExchangeRate) error
}

type exchangeRatesRepository struct {
	db *bun.DB
}

func NewExchangeRatesRepository(db *bun.DB) ExchangeRatesRepository {
	return &exchangeRatesRepository{db: db}
}

func (r *exchangeRatesRepository) Create(ctx context.Context, rate *model.ExchangeRate) error {
	_, err := r.db.NewInsert().Model(rate).Exec(ctx)
	return err
}
