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

type UserAssetsService interface {
	UserAssetAdd(ctx context.Context, vm *viewmodel.TransactionRequest, userID int64) error
	GetUserAssets(ctx context.Context, userID int64, target model.AssetType) (*viewmodel.TotalAssetVM, float64, error)
	GetAllTransaction(ctx context.Context, userID int64, target model.AssetType) ([]*viewmodel.TransactionVM, float64, error)
	GetTransactionPDF(ctx context.Context, userID int64, txID int64, target model.AssetType, fullname string) (*viewmodel.TransactionReceiptPDFVM, float64, error)
	GenerateUserAssetsPDF(ctx context.Context, userID int64, currency model.AssetType, fullname string) (*viewmodel.PortfolioReportPDFVM, float64, error)
}

type userAssetsService struct {
	uar repository.UserAssetsRepository
	ks  KurService
	tr  repository.TransactionRepository
}

func NewUserAssetsService(uar repository.UserAssetsRepository, ks KurService, tr repository.TransactionRepository) UserAssetsService {
	return &userAssetsService{uar: uar, ks: ks, tr: tr}
}

func (s *userAssetsService) UserAssetAdd(ctx context.Context, vm *viewmodel.TransactionRequest, userID int64) error {
	if vm == nil {
		return errors.New("input is required")
	}

	return s.tr.WithTx(ctx, func(uar repository.UserAssetsRepository, tr repository.TransactionRepository) error {

		date, err := parseTransactionDate(vm.TransactionDate)
		if err != nil {
			return err
		}

		price, err := s.ks.FetchFromDovizTransaction(vm.Asset)
		if err != nil {
			return err
		}

		asset, err := uar.GetByUserIDAndAssetType(ctx, userID, vm.Asset)
		if err != nil {
			return err
		}

		if asset == nil {
			if vm.Type == model.TypeSubtract {
				return errors.New("ilk miktar negatif olamaz")
			}

			newAsset := &model.UserAsset{
				UserID: userID,
				Asset:  vm.Asset,
				Amount: vm.Amount,
			}
			if err := uar.UserAssetAdd(ctx, newAsset); err != nil {
				return err
			}
		} else {
			if vm.Type == model.TypeAdd {
				asset.Amount += vm.Amount
			} else {
				asset.Amount -= vm.Amount
			}

			if asset.Amount < 0.0000 {
				return errors.New("asset miktarı negatif olamaz")
			}

			if err := uar.Update(ctx, asset); err != nil {
				return err
			}
		}

		transaction := &model.Transaction{
			UserID:          userID,
			Type:            vm.Type,
			Asset:           vm.Asset,
			Amount:          vm.Amount,
			Price:           price.Price,
			TotalPrice:      price.Price * vm.Amount,
			TransactionDate: date,
			Description:     vm.Description,
		}

		return tr.Create(ctx, transaction)
	})
}

func (s *userAssetsService) GetAllTransaction(ctx context.Context, userID int64, target model.AssetType) ([]*viewmodel.TransactionVM, float64, error) {
	if userID == 0 {
		return nil, 0, errors.New("tokenID is required!")
	}

	t, err := s.tr.GetTransaction(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	kur, err := s.ks.FetchFromDoviz()
	if err != nil {
		return nil, 0, err
	}

	allTVM, targetPrice, err := s.allT(t, kur, target)
	if err != nil {
		return nil, 0, err
	}

	return allTVM, targetPrice, nil
}

func (s *userAssetsService) GetUserAssets(ctx context.Context, userID int64, target model.AssetType) (*viewmodel.TotalAssetVM, float64, error) {
	if userID == 0 {
		return nil, 0, errors.New("tokenID is required!")
	}

	assets, err := s.uar.GetUserAssets(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	kur, err := s.ks.FetchFromDoviz()
	if err != nil {
		return nil, 0, err
	}

	totalVM, targetPrice, err := s.totalAsset(assets, kur, target)
	if err != nil {
		return nil, 0, err
	}

	/*

		converted, err := s.convertFromTRY(totalVM.TotalPrice, target, kur)
		if err != nil {
			return nil, err
		}

	*/

	totalVM.Currency = target
	return totalVM, targetPrice, nil
}

func (s *userAssetsService) GetTransactionPDF(ctx context.Context, userID int64, txID int64, target model.AssetType, fullname string) (*viewmodel.TransactionReceiptPDFVM, float64, error) {
	// 1. Transaction verisini çek
	tx, err := s.tr.GetTransactionByID(ctx, userID, txID)
	if err != nil {
		return nil, 0, err
	}

	// 2. Kur verisini al
	kur, err := s.ks.FetchFromDoviz()
	if err != nil {
		return nil, 0, err
	}

	// 3. Fiyat hesaplamaları
	prices := s.priceMapTRY(kur)
	targetPrice, ok := prices[target]
	if !ok {
		return nil, 0, errors.New("unsupported asset type")
	}

	price, ok := prices[tx.Asset]
	if !ok {
		price = 1
	}

	total := tx.Amount * (price / targetPrice)

	// 4. PDF VM oluştur
	pdfVM := &viewmodel.TransactionReceiptPDFVM{
		NameAndSurname:  fullname,
		BaseCurrency:    string(target),
		AssetName:       string(tx.Asset),
		Type:            string(tx.Type),
		Amount:          tx.Amount,
		UnitPrice:       price / targetPrice,
		TotalPrice:      total,
		TransactionDate: tx.TransactionDate,
		Description:     tx.Description,
	}

	return pdfVM, targetPrice, nil
}

func (s *userAssetsService) GenerateUserAssetsPDF(ctx context.Context, userID int64, currency model.AssetType, fullname string) (*viewmodel.PortfolioReportPDFVM, float64, error) {
	resp, targetPrice, err := s.GetUserAssets(ctx, userID, currency)
	if err != nil {
		return nil, 0, err
	}

	pdfVM := &viewmodel.PortfolioReportPDFVM{
		NameAndSurname: fullname,
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

func (s *userAssetsService) totalAsset(assets []*model.UserAsset, kur *viewmodel.GuncelKurVM, target model.AssetType) (*viewmodel.TotalAssetVM, float64, error) {
	totalAsset := &viewmodel.TotalAssetVM{}
	totalAsset.TotalPrice = 0.0
	vms := make([]viewmodel.UserAssetVM, 0, len(assets))

	prices := s.priceMapTRY(kur)
	targetPrice, ok := prices[target]
	if !ok {
		return nil, 0, errors.New("unsupported asset type")
	}

	for _, k := range assets {
		price, ok := prices[k.Asset]
		if !ok {
			continue
		}

		vm := viewmodel.ToUserAssetVM(k)
		vm.Price = price / targetPrice
		vm.TotalPriceByTargetAsset = vm.Price * vm.Amount

		totalAsset.TotalPrice += vm.TotalPriceByTargetAsset
		vms = append(vms, *vm)
	}
	totalAsset.Assets = vms
	return totalAsset, targetPrice, nil
}

func (s *userAssetsService) convertFromTRY(totalTRY float64, target model.AssetType, kur *viewmodel.GuncelKurVM) (float64, error) {
	prices := s.priceMapTRY(kur)

	priceTRY, ok := prices[target]
	if !ok {
		return 0, errors.New("unsupported asset type")
	}

	return totalTRY / priceTRY, nil
}

func (s *userAssetsService) priceMapTRY(kur *viewmodel.GuncelKurVM) map[model.AssetType]float64 {
	return map[model.AssetType]float64{
		model.AssetTypeTurkLirasi:       1,
		model.AssetTypeDolar:            kur.Dolar,
		model.AssetTypeEuro:             kur.Euro,
		model.AssetTypeSterlin:          kur.Sterlin,
		model.AssetTypeFrank:            kur.Frank,
		model.AssetTypeCeyrekAltin:      kur.CeyrekAltin,
		model.AssetTypeYarimAltin:       kur.YarimAltin,
		model.AssetTypeTamAltin:         kur.TamAltin,
		model.AssetTypeCumhuriyetAltini: kur.CumhuriyetAltini,
		model.AssetTypeBilezik22Ayar:    kur.Bilezik22Ayar,
		model.AssetTypeGramAltin14Ayar:  kur.GramAltin14Ayar,
		model.AssetTypeGramAltin18Ayar:  kur.GramAltin18Ayar,
		model.AssetTypeGramAltin22Ayar:  kur.GramAltin22Ayar,
		model.AssetTypeGramAltin24Ayar:  kur.GramAltin24Ayar,
		model.AssetTypeGumus:            kur.Gumus,
	}
}

func (s *userAssetsService) allT(t []*model.Transaction, kur *viewmodel.GuncelKurVM, target model.AssetType) ([]*viewmodel.TransactionVM, float64, error) {
	vms := make([]*viewmodel.TransactionVM, 0, len(t))

	prices := s.priceMapTRY(kur)
	targetPrice, ok := prices[target]
	if !ok {
		return nil, 0, errors.New("unsupported asset type")
	}

	for _, k := range t {
		price, ok := prices[k.Asset]
		if !ok {
			continue
		}

		vm := viewmodel.ToTransactionVM(k)
		vm.Price = price / targetPrice
		vm.TotalPriceByTargetAsset = vm.Price * vm.Amount

		vms = append(vms, vm)
	}
	return vms, targetPrice, nil

}

func parseTransactionDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Now(), nil // kullanıcı göndermezse şu an
	}

	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("geçersiz tarih formatı, RFC3339 olmalı")
	}

	return t, nil
}
