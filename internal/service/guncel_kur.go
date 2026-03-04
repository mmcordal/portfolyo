package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"portfolyo/internal/infrastructure/cache"
	"portfolyo/internal/model"
	"portfolyo/internal/repository"
	"portfolyo/internal/viewmodel"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const kurCacheKey = "guncel_kur_data"

type cachePayload struct {
	Data *viewmodel.GuncelKurVM `json:"data"`
}

type KurService interface {
	FetchFromDoviz() (*viewmodel.GuncelKurVM, error)
	FetchFromDovizForTransaction(asset model.AssetType) (*viewmodel.CurrentByAssetType, error)
	FetchAndSaveRates() error
}

type kurService struct {
	cache *cache.RedisClient
	rate  repository.ExchangeRatesRepository
}

func NewKurService(c *cache.RedisClient, r repository.ExchangeRatesRepository) KurService {
	return &kurService{cache: c, rate: r}
}

func parsePrice(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println("price parse error:", err)
		return 0
	}
	return f
}

func (s *kurService) FetchFromDoviz() (*viewmodel.GuncelKurVM, error) {
	cachedStr, err := s.cache.Get(context.Background(), kurCacheKey)
	if err == nil && cachedStr != "" {
		var payload cachePayload
		if json.Unmarshal([]byte(cachedStr), &payload) == nil {
			log.Println("Redis cache kullanıldı")
			return payload.Data, nil
		}
	} else {
		log.Println("Redis'te cache yok veya okunamadı")
	}

	// doviz verileri için 1. istek
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get("https://kur.doviz.com")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kur.doviz response not OK: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	kur := &model.GuncelKur{}
	doc.Find(`td[data-socket-attr="ask"]`).Each(func(i int, s *goquery.Selection) {
		key, exists := s.Attr("data-socket-key")
		if !exists {
			return
		}

		price := parsePrice(s.Text())

		switch key {
		case "USD":
			kur.Dolar = price
		case "EUR":
			kur.Euro = price
		case "GBP":
			kur.Sterlin = price
		case "CHF":
			kur.Frank = price
		}
	})

	// altin verileri için 2. istek
	client = &http.Client{Timeout: 120 * time.Second}
	resp, err = client.Get("https://altin.doviz.com")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("altin.doviz response not OK: %d", resp.StatusCode)
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(`td[data-socket-attr="ask"]`).Each(func(i int, s *goquery.Selection) {
		key, exists := s.Attr("data-socket-key")
		if !exists {
			return
		}

		price := parsePrice(s.Text())

		switch key {
		case "ceyrek-altin":
			kur.CeyrekAltin = price
		case "yarim-altin":
			kur.YarimAltin = price
		case "tam-altin":
			kur.TamAltin = price
		case "cumhuriyet-altini":
			kur.CumhuriyetAltini = price
		case "gram-altin":
			kur.GramAltin24Ayar = price
		case "22-ayar-bilezik":
			kur.GramAltin22Ayar = price
			kur.Bilezik22Ayar = price
		case "18-ayar-altin":
			kur.GramAltin18Ayar = price
		case "14-ayar-altin":
			kur.GramAltin14Ayar = price
		case "gumus":
			kur.Gumus = price
		}
	})

	kur.CreatedAt = time.Now()
	vm := viewmodel.ToGuncelKurVM(kur)

	payload := cachePayload{Data: vm}
	bytes, _ := json.Marshal(payload)
	_ = s.cache.Set(context.Background(), kurCacheKey, string(bytes), 120*time.Second)

	log.Println("Güncel kur verileri başarıyla alındı:", fmt.Sprintf("%+v", kur))

	return vm, nil
}

func (s *kurService) FetchFromDovizForTransaction(asset model.AssetType) (*viewmodel.CurrentByAssetType, error) {

	kur, err := s.fetchFromDovizForTransaction()
	if err != nil {
		return nil, err
	}

	var price float64

	switch asset {
	case model.AssetTypeTurkLirasi:
		price = 1
	case model.AssetTypeDolar:
		price = kur.Dolar
	case model.AssetTypeEuro:
		price = kur.Euro
	case model.AssetTypeSterlin:
		price = kur.Sterlin
	case model.AssetTypeFrank:
		price = kur.Frank
	case model.AssetTypeCeyrekAltin:
		price = kur.CeyrekAltin
	case model.AssetTypeYarimAltin:
		price = kur.YarimAltin
	case model.AssetTypeTamAltin:
		price = kur.TamAltin
	case model.AssetTypeCumhuriyetAltini:
		price = kur.CumhuriyetAltini
	case model.AssetTypeBilezik22Ayar:
		price = kur.Bilezik22Ayar
	case model.AssetTypeGramAltin14Ayar:
		price = kur.GramAltin14Ayar
	case model.AssetTypeGramAltin18Ayar:
		price = kur.GramAltin18Ayar
	case model.AssetTypeGramAltin22Ayar:
		price = kur.GramAltin22Ayar
	case model.AssetTypeGramAltin24Ayar:
		price = kur.GramAltin24Ayar
	case model.AssetTypeGumus:
		price = kur.Gumus
	default:
		return nil, fmt.Errorf("unknown asset type")
	}

	return &viewmodel.CurrentByAssetType{
		Asset: asset,
		Price: price,
	}, nil
}

func (s *kurService) FetchAndSaveRates() error {
	kur, err := s.fetchFromDovizForTransaction()
	if err != nil {
		return err
	}

	rate := &model.ExchangeRate{
		Dolar:            kur.Dolar,
		Sterlin:          kur.Sterlin,
		Euro:             kur.Euro,
		Frank:            kur.Frank,
		CeyrekAltin:      kur.CeyrekAltin,
		YarimAltin:       kur.YarimAltin,
		TamAltin:         kur.TamAltin,
		CumhuriyetAltini: kur.CumhuriyetAltini,
		Bilezik22Ayar:    kur.Bilezik22Ayar,
		GramAltin14Ayar:  kur.GramAltin14Ayar,
		GramAltin18Ayar:  kur.GramAltin18Ayar,
		GramAltin22Ayar:  kur.GramAltin22Ayar,
		GramAltin24Ayar:  kur.GramAltin24Ayar,
		Gumus:            kur.Gumus,
	}

	return s.rate.Create(context.Background(), rate)
}

func (s *kurService) fetchFromDovizForTransaction() (*viewmodel.GuncelKurVM, error) {
	// doviz verileri için 1. istek
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get("https://kur.doviz.com")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kur.doviz response not OK: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	kur := &model.GuncelKur{}
	doc.Find(`td[data-socket-attr="ask"]`).Each(func(i int, s *goquery.Selection) {
		key, exists := s.Attr("data-socket-key")
		if !exists {
			return
		}

		price := parsePrice(s.Text())

		switch key {
		case "USD":
			kur.Dolar = price
		case "EUR":
			kur.Euro = price
		case "GBP":
			kur.Sterlin = price
		case "CHF":
			kur.Frank = price
		}
	})

	// altin verileri için 2. istek
	client = &http.Client{Timeout: 120 * time.Second}
	resp, err = client.Get("https://altin.doviz.com")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("altin.doviz response not OK: %d", resp.StatusCode)
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(`td[data-socket-attr="ask"]`).Each(func(i int, s *goquery.Selection) {
		key, exists := s.Attr("data-socket-key")
		if !exists {
			return
		}

		price := parsePrice(s.Text())

		switch key {
		case "ceyrek-altin":
			kur.CeyrekAltin = price
		case "yarim-altin":
			kur.YarimAltin = price
		case "tam-altin":
			kur.TamAltin = price
		case "cumhuriyet-altini":
			kur.CumhuriyetAltini = price
		case "gram-altin":
			kur.GramAltin24Ayar = price
		case "22-ayar-bilezik":
			kur.GramAltin22Ayar = price
			kur.Bilezik22Ayar = price
		case "18-ayar-altin":
			kur.GramAltin18Ayar = price
		case "14-ayar-altin":
			kur.GramAltin14Ayar = price
		case "gumus":
			kur.Gumus = price
		}
	})

	kur.CreatedAt = time.Now()
	vm := viewmodel.ToGuncelKurVM(kur)

	payload := cachePayload{Data: vm}
	bytes, _ := json.Marshal(payload)
	_ = s.cache.Set(context.Background(), kurCacheKey, string(bytes), 120*time.Second)

	log.Println("Güncel kur verileri başarıyla alındı:", fmt.Sprintf("%+v", kur))

	return vm, nil
}
