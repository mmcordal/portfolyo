package service

import (
	"context"
	"errors"
	"fmt"
	"portfolyo/internal/model"
	"portfolyo/internal/repository"
	"portfolyo/internal/viewmodel"
	"time"
)

type TransactionService interface {
	TransactionAdd(ctx context.Context, vm *viewmodel.TransactionRequest, userID int64) error
	GetAllTransactionByAsset(ctx context.Context, userID int64, target model.AssetType, assetType model.AssetType) ([]*viewmodel.TransactionVM, float64, error)
	GetAllTransaction(ctx context.Context, userID int64, target model.AssetType) ([]*viewmodel.TransactionVM, float64, error)
	GetTransactionPDF(ctx context.Context, userID int64, txID int64, target model.AssetType) (*viewmodel.TransactionReceiptPDFVM, error)
}

type transactionService struct {
	tr  repository.TransactionRepository
	ks  KurService
	uar repository.UserAssetsRepository
}

func NewTransactionService(tr repository.TransactionRepository, ks KurService, uar repository.UserAssetsRepository) TransactionService {
	return &transactionService{tr: tr, ks: ks, uar: uar}
}

func (s *transactionService) TransactionAdd(ctx context.Context, vm *viewmodel.TransactionRequest, userID int64) error {
	if vm == nil {
		return errors.New("input is required")
	}

	typeAction, err := model.IsValidActionType(vm.Type)
	if err != nil {
		return err
	}
	assetType, err := model.IsValidAssetType(vm.Asset)
	if err != nil {
		return err
	}

	price, err := s.ks.FetchFromDovizForTransaction(assetType)
	if err != nil {
		return err
	}

	date, err := parseTransactionDate(vm.TransactionDate)
	if err != nil {
		return err
	}

	return s.tr.WithTx(ctx, func(uar repository.UserAssetsRepository, tr repository.TransactionRepository) error {

		asset, err := uar.FindOrCreateByUserAndAsset(ctx, userID, assetType)
		if err != nil {
			return err
		}

		if typeAction == model.TypeAdd {
			asset.Amount += vm.Amount
		} else {
			asset.Amount -= vm.Amount
		}

		if asset.Amount < 0.0000 {
			return errors.New("Varlık miktarı negatif olamaz")
		}

		asset.UpdatedAt = time.Now()
		if err := uar.Update(ctx, asset); err != nil {
			return err
		}

		transaction := &model.Transaction{
			AssetID:         asset.ID,
			UserAsset:       asset,
			Type:            typeAction,
			Amount:          vm.Amount,
			Price:           price.Price,
			TotalPrice:      price.Price * vm.Amount,
			TransactionDate: date,
			Description:     vm.Description,
		}

		return tr.Create(ctx, transaction)
	})
}

func (s *transactionService) GetAllTransactionByAsset(ctx context.Context, userID int64, target model.AssetType, assetType model.AssetType) ([]*viewmodel.TransactionVM, float64, error) {

	asset, err := s.uar.GetUserAssetWithTransactionByAsset(ctx, userID, assetType)
	if err != nil {
		return nil, 0, err
	}

	t, err := s.tr.GetAllTransactionByAsset(ctx, asset.ID)
	if err != nil {
		return nil, 0, err
	}

	price, err := s.ks.FetchFromDovizForTransaction(asset.Asset)
	if err != nil {
		return nil, 0, err
	}

	kur, err := s.ks.FetchFromDoviz()
	if err != nil {
		return nil, 0, err
	}

	allTVM, targetPrice, err := allTByAsset(t, kur, price, target)
	if err != nil {
		return nil, 0, err
	}

	return allTVM, targetPrice, nil
}

func (s *transactionService) GetAllTransaction(ctx context.Context, userID int64, target model.AssetType) ([]*viewmodel.TransactionVM, float64, error) {

	assets, err := s.uar.GetUserAssets(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	t, err := s.tr.GetAllTransaction(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	kur, err := s.ks.FetchFromDoviz()
	if err != nil {
		return nil, 0, err
	}

	allTVM, targetPrice, err := allT(t, kur, assets, target)
	if err != nil {
		return nil, 0, err
	}

	return allTVM, targetPrice, nil
}

var ErrTransactionNotFoundOrUnauthorized = errors.New("transaction not found or unauthorized")

func (s *transactionService) GetTransactionPDF(ctx context.Context, userID int64, txID int64, target model.AssetType) (*viewmodel.TransactionReceiptPDFVM, error) {
	tx, err := s.tr.GetTransactionByIDAndUserID(ctx, txID, userID)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, fmt.Errorf("%w: transaction not found", ErrNotFound)
	}
	if tx.UserAsset == nil || tx.UserAsset.User == nil {
		return nil, fmt.Errorf("%w: transaction relation not found", ErrNotFound)
	}

	kur, err := s.ks.FetchFromDoviz()
	if err != nil {
		return nil, err
	}

	prices := priceMapTRY(kur)
	targetPrice, ok := prices[target]
	if !ok {
		return nil, errors.New("unsupported asset type")
	}

	price, ok := prices[tx.UserAsset.Asset]
	if !ok {
		price = 1
	}

	pdfVM := &viewmodel.TransactionReceiptPDFVM{
		CreatedAt:       tx.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       tx.UpdatedAt.Format("2006-01-02 15:04:05"),
		NameAndSurname:  tx.UserAsset.User.String(),
		AssetName:       string(tx.UserAsset.Asset),
		Type:            string(tx.Type),
		Amount:          tx.Amount,
		UnitPrice:       tx.Price,
		TotalPrice:      tx.TotalPrice,
		TransactionDate: tx.TransactionDate.Format("2006-01-02 15:04:05"),
		Description:     tx.Description,
		TargetPrice:     targetPrice,

		BaseCurrency:           string(target),
		BaseCurrencyPrice:      price / targetPrice,
		BaseCurrencyTotalPrice: tx.Amount * (price / targetPrice),
	}

	return pdfVM, nil
}

var ErrNotFound = errors.New("record not found")
