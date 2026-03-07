package service

import (
	"context"
	"fmt"
	"portfolyo/internal/model"
	"portfolyo/internal/repository"
	"portfolyo/internal/viewmodel"
	"time"
)

type UserAssetsService interface {
	GetUserAssets(ctx context.Context, userID int64, target model.AssetType) (*viewmodel.TotalAssetVM, float64, error)
	GetUserAsset(ctx context.Context, userID int64, target model.AssetType, assetType model.AssetType) (*viewmodel.UserAssetVM, float64, error)
	GenerateUserAssetsPDF(ctx context.Context, userID int64, currency model.AssetType) (*viewmodel.PortfolioReportPDFVM, float64, error)
}

type userAssetsService struct {
	uar repository.UserAssetsRepository
	ks  KurService
}

func NewUserAssetsService(uar repository.UserAssetsRepository, ks KurService) UserAssetsService {
	return &userAssetsService{uar: uar, ks: ks}
}

func (s *userAssetsService) GetUserAssets(ctx context.Context, userID int64, target model.AssetType) (*viewmodel.TotalAssetVM, float64, error) {

	assets, err := s.uar.GetUserAssets(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	kur, err := s.ks.FetchFromDoviz()
	if err != nil {
		return nil, 0, err
	}

	totalVM, targetPrice, err := totalAsset(assets, kur, target)
	if err != nil {
		return nil, 0, err
	}

	return totalVM, targetPrice, nil
}

func (s *userAssetsService) GetUserAsset(ctx context.Context, userID int64, target model.AssetType, assetType model.AssetType) (*viewmodel.UserAssetVM, float64, error) {
	asset, err := s.uar.GetUserAssetWithTransactionByAsset(ctx, userID, assetType)
	if err != nil {
		return nil, 0, err
	}
	if asset == nil {
		return nil, 0, fmt.Errorf("%w: user asset not found", ErrNotFound)
	}

	kur, err := s.ks.FetchFromDoviz()
	if err != nil {
		return nil, 0, err
	}

	prices := priceMapTRY(kur)
	targetPrice, ok := prices[target]
	if !ok {
		return nil, 0, err
	}

	price, ok := prices[asset.Asset]
	if !ok {
		return nil, 0, err
	}

	c := new(viewmodel.CurrentByAssetType)
	c.Asset = target

	vm := viewmodel.ToUserAssetVM(asset)
	vm.TargetCurrency = string(target)
	vm.Price = price / targetPrice
	vm.TotalPriceByTargetAsset = vm.Price * vm.Amount
	c.Price = vm.Price
	vm.Transactions = viewmodel.ToTransactionVMs(asset.Transactions, c)

	return vm, targetPrice, nil
}

func (s *userAssetsService) GenerateUserAssetsPDF(ctx context.Context, userID int64, currency model.AssetType) (*viewmodel.PortfolioReportPDFVM, float64, error) {
	resp, targetPrice, err := s.GetUserAssets(ctx, userID, currency)
	if err != nil {
		return nil, 0, err
	}

	pdfVM := &viewmodel.PortfolioReportPDFVM{
		NameAndSurname: resp.FullName,
		ReportDate:     time.Now(),
		BaseCurrency:   string(resp.Currency),
		Rows:           make([]viewmodel.PortfolioRowPDFVM, len(resp.Assets)),
		TotalValue:     resp.TotalPrice,
	}

	for i, asset := range resp.Assets {
		pdfVM.Rows[i] = viewmodel.PortfolioRowPDFVM{
			AssetName:  asset.Asset,
			Amount:     asset.Amount,
			UnitPrice:  asset.Price,
			TotalPrice: asset.TotalPriceByTargetAsset,
		}
	}

	return pdfVM, targetPrice, nil
}
