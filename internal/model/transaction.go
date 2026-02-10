package model

import (
	"time"

	"github.com/uptrace/bun"
)

type TypeAction string

const (
	TypeAdd      TypeAction = "add"
	TypeSubtract TypeAction = "subtract"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`
	CoreModel
	UserID          int64      `json:"user_id"`
	Type            TypeAction `json:"type"`
	Asset           AssetType  `json:"asset"`
	Amount          float64    `json:"amount"`
	Price           float64    `json:"price"`
	TotalPrice      float64    `json:"total_price"` // price * amount
	TransactionDate time.Time  `json:"transaction_date"`
	Description     string     `json:"description,omitempty"`
}
