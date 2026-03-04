package model

import (
	"errors"
	"strings"
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

	AssetID         int64      `bun:",notnull"`
	Type            TypeAction `bun:",type:varchar(20),notnull"`
	Amount          float64    `bun:",notnull,check:amount > 0"`
	Price           float64    `bun:",notnull"`
	TotalPrice      float64    `bun:",notnull"`
	TransactionDate time.Time  `bun:",notnull"`
	Description     string     `bun:",nullzero"`

	UserAsset *UserAsset `bun:"rel:belongs-to,join:asset_id=id"`
}

var validAssetTypes = map[AssetType]struct{}{
	AssetTypeTurkLirasi:       {},
	AssetTypeDolar:            {},
	AssetTypeSterlin:          {},
	AssetTypeEuro:             {},
	AssetTypeFrank:            {},
	AssetTypeCeyrekAltin:      {},
	AssetTypeYarimAltin:       {},
	AssetTypeTamAltin:         {},
	AssetTypeCumhuriyetAltini: {},
	AssetTypeBilezik22Ayar:    {},
	AssetTypeGramAltin14Ayar:  {},
	AssetTypeGramAltin18Ayar:  {},
	AssetTypeGramAltin22Ayar:  {},
	AssetTypeGramAltin24Ayar:  {},
	AssetTypeGumus:            {},
}

func IsValidAssetType(currency string) (AssetType, error) {
	normalized := AssetType(strings.ToLower(currency))
	if _, ok := validAssetTypes[normalized]; !ok {
		return "", errors.New("invalid asset type")
	}
	return normalized, nil
}

var validActionTypes = map[TypeAction]struct{}{
	TypeAdd:      {},
	TypeSubtract: {},
}

func IsValidActionType(action string) (TypeAction, error) {
	normalized := TypeAction(strings.ToLower(action))
	if _, ok := validActionTypes[normalized]; !ok {
		return "", errors.New("invalid asset type")
	}
	return normalized, nil
}
