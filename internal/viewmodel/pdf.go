package viewmodel

import (
	"time"
)

type PortfolioReportPDFVM struct {
	NameAndSurname string
	ReportDate     time.Time
	BaseCurrency   string
	Rows           []PortfolioRowPDFVM
	TotalValue     float64
}

type PortfolioRowPDFVM struct {
	AssetName  string
	Amount     float64
	UnitPrice  float64
	TotalPrice float64
}

type TransactionReceiptPDFVM struct {
	CreatedAt       string
	UpdatedAt       string
	NameAndSurname  string
	AssetName       string
	Type            string
	Amount          float64
	UnitPrice       float64
	TotalPrice      float64
	TransactionDate string
	Description     string
	TargetPrice     float64

	BaseCurrency           string
	BaseCurrencyPrice      float64
	BaseCurrencyTotalPrice float64
}
