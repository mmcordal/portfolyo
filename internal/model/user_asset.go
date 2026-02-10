package model

import "github.com/uptrace/bun"

type AssetType string

const (
	AssetTypeTurkLirasi       AssetType = "TRY"
	AssetTypeDolar            AssetType = "USD"
	AssetTypeSterlin          AssetType = "GBP"
	AssetTypeEuro             AssetType = "EUR"
	AssetTypeFrank            AssetType = "CHF"
	AssetTypeCeyrekAltin      AssetType = "ceyrek-altin"
	AssetTypeYarimAltin       AssetType = "yarim-altin"
	AssetTypeTamAltin         AssetType = "tam-altin"
	AssetTypeCumhuriyetAltini AssetType = "cumhuriyet-altini"
	AssetTypeBilezik22Ayar    AssetType = "22-ayar-bilezik"
	AssetTypeGramAltin14Ayar  AssetType = "14-ayar-altin"
	AssetTypeGramAltin18Ayar  AssetType = "18-ayar-altin"
	AssetTypeGramAltin22Ayar  AssetType = "22-ayar-altin" // bilezikle aynı
	AssetTypeGramAltin24Ayar  AssetType = "gram-altin"
	AssetTypeGumus            AssetType = "gumus"
)

type UserAsset struct {
	bun.BaseModel `bun:"user_asset,alias:ua"`
	CoreModel
	UserID int64     `json:"user_id" bun:"type:bigint,notnull"`
	Asset  AssetType `json:"asset" bun:"type:varchar(255),notnull"`
	Amount float64   `json:"amount"`
}
