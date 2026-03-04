package viewmodel

import "portfolyo/internal/model"

type UserAssetVM struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	User         *UserVM          `json:"user,omitempty"`
	Transactions []*TransactionVM `json:"transactions,omitempty"`

	UserID                  int64   `json:"user_id"`
	Asset                   string  `json:"asset"`
	Amount                  float64 `json:"amount"`
	Price                   float64 `json:"price,omitempty"`
	TargetCurrency          string  `json:"target_currency,omitempty"`
	TotalPriceByTargetAsset float64 `json:"total_price_by_asset,omitempty"`
}

type TotalAssetVM struct {
	Assets     []*UserAssetVM `json:"assets"`
	TotalPrice float64        `json:"total_price,omitempty"`
	Currency   string         `json:"currency,omitempty"`
	User       *UserVM        `json:"user,omitempty"`
	FullName   string         `json:"full_name"`
}

func ToUserAssetVM(a *model.UserAsset) *UserAssetVM {
	if a == nil {
		return nil
	}
	vm := &UserAssetVM{
		ID:        a.ID,
		UserID:    a.UserID,
		Asset:     string(a.Asset),
		Amount:    a.Amount,
		CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: a.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if a.User != nil {
		vm.User = ToUserVM(a.User)
	}
	return vm
}

func ToUserAssetVMs(assets []*model.UserAsset, c *CurrentByAssetType) []*UserAssetVM {
	out := make([]*UserAssetVM, 0, len(assets))
	for _, a := range assets {
		vm := ToUserAssetVM(a)
		vm.Transactions = ToTransactionVMs(a.Transactions, c)
		out = append(out, vm)
	}
	return out
}
