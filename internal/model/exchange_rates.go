package model

import (
	"github.com/uptrace/bun"
)

type ExchangeRate struct {
	bun.BaseModel `bun:"table:exchange_rates,alias:er"`
	CoreModel

	Dolar            float64 `bun:",notnull"`
	Sterlin          float64 `bun:",notnull"`
	Euro             float64 `bun:",notnull"`
	Frank            float64 `bun:",notnull"`
	CeyrekAltin      float64 `bun:",notnull"`
	YarimAltin       float64 `bun:",notnull"`
	TamAltin         float64 `bun:",notnull"`
	CumhuriyetAltini float64 `bun:",notnull"`
	Bilezik22Ayar    float64 `bun:",notnull"`
	GramAltin14Ayar  float64 `bun:",notnull"`
	GramAltin18Ayar  float64 `bun:",notnull"`
	GramAltin22Ayar  float64 `bun:",notnull"`
	GramAltin24Ayar  float64 `bun:",notnull"`
	Gumus            float64 `bun:",notnull"`
}
