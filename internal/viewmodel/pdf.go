package viewmodel

import (
	"portfolyo/internal/model"
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
	AssetName  model.AssetType
	Amount     float64
	UnitPrice  float64
	TotalPrice float64
}

type TransactionReceiptPDFVM struct {
	NameAndSurname  string
	BaseCurrency    string
	AssetName       string
	Type            string
	Amount          float64
	UnitPrice       float64
	TotalPrice      float64
	TransactionDate time.Time
	Description     string
}
