package viewmodel

import "portfolyo/internal/model"

type TransactionRequest struct {
	Type            model.TypeAction `json:"type" validate:"required" labelName:"Aksiyon Tipi"`
	Asset           model.AssetType  `json:"asset" validate:"required" labelName:"Varlık"`
	Amount          float64          `json:"amount" validate:"required,gte=0" labelName:"Miktar"`
	TransactionDate string           `json:"transaction_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00" labelName:"İşlem Tarihi"`
	Description     string           `json:"description" validate:"omitempty" labelName:"Açıklama"`
}

type TotalAssetVM struct {
	Assets     []UserAssetVM   `json:"assets"`
	TotalPrice float64         `json:"total_price"`
	Currency   model.AssetType `json:"currency"`
}

type UserAssetVM struct {
	CreatedAt               string          `json:"created_at"`
	UpdatedAt               string          `json:"updated_at"`
	ID                      int64           `json:"id"`
	UserID                  int64           `json:"user_id"`
	Asset                   model.AssetType `json:"asset"`
	Amount                  float64         `json:"amount"`
	Price                   float64         `json:"price"`
	TotalPriceByTargetAsset float64         `json:"total_price_by_asset"`
}

type TransactionVM struct {
	CreatedAt               string           `json:"created_at"`
	ID                      int64            `json:"id"`
	UserID                  int64            `json:"user_id"`
	Type                    model.TypeAction `json:"type"`
	Asset                   model.AssetType  `json:"asset"`
	Amount                  float64          `json:"amount"`
	Price                   float64          `json:"price"`
	TotalPriceByTargetAsset float64          `json:"total_price_by_asset"`
	TransactionDate         string           `json:"transaction_date"`
	Description             string           `json:"description"`
}

func ToTransactionVM(k *model.Transaction) *TransactionVM {
	return &TransactionVM{
		CreatedAt:               k.CreatedAt.Format("2006-01-02 15:04:05"),
		ID:                      k.ID,
		UserID:                  k.UserID,
		Type:                    k.Type,
		Asset:                   k.Asset,
		Amount:                  k.Amount,
		Price:                   k.Price,
		TotalPriceByTargetAsset: k.TotalPrice,
		TransactionDate:         k.TransactionDate.Format("2006-01-02 15:04:05"),
		Description:             k.Description,
	}
}

func ToUserAssetVM(k *model.UserAsset) *UserAssetVM {
	return &UserAssetVM{
		CreatedAt: k.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: k.UpdatedAt.Format("2006-01-02 15:04:05"),
		ID:        k.ID,
		UserID:    k.UserID,
		Asset:     k.Asset,
		Amount:    k.Amount,
	}
}

func ToUserAssetVMs(assets []*model.UserAsset) []*UserAssetVM {
	vms := make([]*UserAssetVM, 0, len(assets))
	for _, k := range assets {
		vms = append(vms, ToUserAssetVM(k))
	}
	return vms
}
