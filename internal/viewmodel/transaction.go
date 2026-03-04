package viewmodel

import "portfolyo/internal/model"

type TransactionRequest struct {
	Type            string  `json:"type" validate:"required,oneof=add subtract" labelName:"tip"`
	Asset           string  `json:"asset" validate:"required" labelName:"varlık"`
	Amount          float64 `json:"amount" validate:"required,gt=0" labelName:"miktar"`
	TransactionDate string  `json:"transaction_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00" labelName:"işlem tarihi"`
	Description     string  `json:"description" validate:"omitempty" labelName:"açıklama"`
}

type TransactionVM struct {
	ID              int64            `json:"id"`
	CreatedAt       string           `json:"created_at"`
	AssetID         int64            `json:"asset_id"`
	Type            model.TypeAction `json:"type"`
	Amount          float64          `json:"amount"`
	Price           float64          `json:"price"`
	TotalPrice      float64          `json:"total_price"`
	TransactionDate string           `json:"transaction_date"`
	Description     string           `json:"description"`

	NowTargetCurrency        string  `json:"now_target_currency,omitempty"`
	TargetCurrencyPrice      float64 `json:"target_currency_price,omitempty"`
	TargetCurrencyTotalPrice float64 `json:"target_currency_total_price,omitempty"`

	UserAsset *UserAssetVM `json:"user_asset,omitempty"`
}

func ToTransactionVM(t *model.Transaction) *TransactionVM {
	vm := &TransactionVM{
		ID:              t.ID,
		CreatedAt:       t.CreatedAt.Format("2006-01-02 15:04:05"),
		AssetID:         t.AssetID,
		Type:            t.Type,
		Amount:          t.Amount,
		Price:           t.Price,
		TotalPrice:      t.TotalPrice,
		TransactionDate: t.TransactionDate.Format("2006-01-02 15:04:05"),
		Description:     t.Description,
	}

	if t.UserAsset != nil {
		vm.UserAsset = ToUserAssetVM(t.UserAsset)
	}

	return vm
}

func ToTransactionVMs(ts []*model.Transaction, asset *CurrentByAssetType) []*TransactionVM {
	out := make([]*TransactionVM, 0, len(ts))
	for _, t := range ts {
		vm := ToTransactionVM(t)
		vm.NowTargetCurrency = string(asset.Asset)
		vm.TargetCurrencyPrice = asset.Price
		vm.TargetCurrencyTotalPrice = vm.TargetCurrencyPrice * vm.Amount
		out = append(out, vm)
	}
	return out
}
