package model

import "github.com/uptrace/bun"

type AssetType string

const (
	AssetTypeTurkLirasi       AssetType = "try"
	AssetTypeDolar            AssetType = "usd"
	AssetTypeSterlin          AssetType = "gbp"
	AssetTypeEuro             AssetType = "eur"
	AssetTypeFrank            AssetType = "chf"
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
	bun.BaseModel `bun:"table:user_assets,alias:ua"`
	CoreModel

	UserID       int64          `bun:",notnull,unique:idx_user_asset"`
	User         *User          `bun:"rel:belongs-to,join:user_id=id"`
	Asset        AssetType      `bun:",type:varchar(50),notnull,unique:idx_user_asset"`
	Amount       float64        `bun:",notnull"`
	Transactions []*Transaction `bun:"rel:has-many,join:id=asset_id"`
}
