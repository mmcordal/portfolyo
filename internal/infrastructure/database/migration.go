package database

import (
	"context"
	"portfolyo/internal/model"

	"github.com/uptrace/bun"
)

func AutoMigration(db *bun.DB) {
	ctx := context.Background()

	models := []interface{}{
		(*model.User)(nil),
		(*model.UserAsset)(nil),
		(*model.Transaction)(nil),
		(*model.Reminder)(nil),
	}

	for _, m := range models {
		_, err := db.NewCreateTable().
			Model(m).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			panic(err)
		}
	}
}
