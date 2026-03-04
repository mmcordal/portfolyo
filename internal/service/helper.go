package service

import (
	"errors"
	"fmt"
	"log"
	"portfolyo/internal/model"
	"portfolyo/internal/viewmodel"
	"time"
)

func totalAsset(assets []*model.UserAsset, kur *viewmodel.GuncelKurVM, target model.AssetType) (*viewmodel.TotalAssetVM, float64, error) {
	totalAsset := new(viewmodel.TotalAssetVM)
	totalAsset.TotalPrice = 0.0
	vms := make([]*viewmodel.UserAssetVM, 0, len(assets))

	prices := priceMapTRY(kur)
	targetPrice, ok := prices[target]
	if !ok {
		return nil, 0, errors.New("unsupported asset type")
	}
	if len(assets) == 0 {
		return nil, 0, nil
	}

	if assets[0].ID > 0 {
		totalAsset.User = viewmodel.ToUserVM(assets[0].User)
		totalAsset.FullName = assets[0].User.String()
	}

	c := new(viewmodel.CurrentByAssetType)
	c.Asset = target

	for _, k := range assets {
		price, ok := prices[k.Asset]
		if !ok {
			continue
		}
		vm := viewmodel.ToUserAssetVM(k)
		vm.Price = price / targetPrice
		vm.TotalPriceByTargetAsset = vm.Price * vm.Amount
		c.Price = vm.Price
		vm.Transactions = viewmodel.ToTransactionVMs(k.Transactions, c)

		totalAsset.TotalPrice += vm.TotalPriceByTargetAsset

		vms = append(vms, vm)
	}
	totalAsset.Assets = vms
	totalAsset.Currency = string(target)
	return totalAsset, targetPrice, nil
}

func convertFromTRY(totalTRY float64, target model.AssetType, kur *viewmodel.GuncelKurVM) (float64, error) {
	prices := priceMapTRY(kur)

	priceTRY, ok := prices[target]
	if !ok {
		return 0, errors.New("unsupported asset type")
	}

	return totalTRY / priceTRY, nil
}

func priceMapTRY(kur *viewmodel.GuncelKurVM) map[model.AssetType]float64 {
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

func assetMap(assets []*model.UserAsset) map[model.AssetType]*model.UserAsset {
	m := make(map[model.AssetType]*model.UserAsset)
	for _, a := range assets {
		m[a.Asset] = a
	}
	return m
}

func allTByAsset(t []*model.Transaction, kur *viewmodel.GuncelKurVM, price *viewmodel.CurrentByAssetType, target model.AssetType) ([]*viewmodel.TransactionVM, float64, error) {
	vms := make([]*viewmodel.TransactionVM, 0, len(t))

	prices := priceMapTRY(kur)
	targetPrice, ok := prices[target]
	if !ok {
		return nil, 0, errors.New("unsupported asset type")
	}

	for _, k := range t {
		vm := viewmodel.ToTransactionVM(k)

		vm.NowTargetCurrency = string(target)
		vm.TargetCurrencyPrice = price.Price / targetPrice
		vm.TargetCurrencyTotalPrice = vm.TargetCurrencyPrice * vm.Amount

		vms = append(vms, vm)
	}
	return vms, targetPrice, nil
}

func allT(t []*model.Transaction, kur *viewmodel.GuncelKurVM, assets []*model.UserAsset, target model.AssetType) ([]*viewmodel.TransactionVM, float64, error) {
	vms := make([]*viewmodel.TransactionVM, 0, len(t))

	assetMap := assetMap(assets)

	prices := priceMapTRY(kur)
	targetPrice, ok := prices[target]
	if !ok {
		return nil, 0, errors.New("unsupported asset type")
	}

	for _, k := range t {
		price, ok := prices[k.UserAsset.Asset]
		if !ok {
			log.Println("an unsupported curreny")
			continue
		}

		vm := viewmodel.ToTransactionVM(k)
		vm.NowTargetCurrency = string(target)
		vm.TargetCurrencyPrice = price / targetPrice
		vm.TargetCurrencyTotalPrice = vm.TargetCurrencyPrice * vm.Amount

		asset, ok := assetMap[k.UserAsset.Asset]
		if !ok {
			log.Println("an unsupported asset")
			continue
		}
		vm.UserAsset = viewmodel.ToUserAssetVM(asset)

		vms = append(vms, vm)
	}
	return vms, targetPrice, nil
}

func parseTransactionDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Now(), nil // kullanıcı göndermezse şuan
	}

	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("geçersiz tarih formatı, RFC3339 olmalı")
	}
	return t, nil
}
