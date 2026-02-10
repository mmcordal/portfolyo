package document

import (
	"fmt"
	"portfolyo/internal/model"
	"portfolyo/internal/viewmodel"

	"github.com/xuri/excelize/v2"
)

func GenerateTransactionsExcelOnceki(list []*viewmodel.TransactionVM, base model.AssetType) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Transactions"
	f.SetSheetName("Sheet1", sheet)

	// Başlıklar
	headers := []string{
		"ID",
		"Tarih",
		"Varlık",
		"İşlem Tipi",
		"Miktar",
		"Birim Fiyat",
		fmt.Sprintf("Toplam (%s)", base),
		"Açıklama",
	}

	for i, h := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet, cell, h)
	}

	// Satırlar
	for i, tx := range list {
		row := i + 2

		info := assetInfo(string(tx.Asset))

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), tx.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), tx.TransactionDate)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), info.Label)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), string(tx.Type))
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("%.4f %s", tx.Amount, info.Unit))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("%.4f", tx.Price))
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("%.4f", tx.TotalPriceByTargetAsset))
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), tx.Description)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GenerateTransactionsExcel(list []*viewmodel.TransactionVM, targetPrice float64, base model.AssetType) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Transactions"
	f.SetSheetName("Sheet1", sheet)

	baseInfo := assetInfo(string(base))

	// ===== STYLES =====
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	addStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#008000"},
	})

	subStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#FF0000"},
	})

	moneyStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt: 2,
	})

	// ===== RAPOR HEADER =====
	f.SetCellValue(sheet, "A1", "KULLANICI İŞLEM RAPORU")
	f.MergeCell(sheet, "A1", "H1")
	f.SetCellStyle(sheet, "A1", "A1", headerStyle)

	f.SetCellValue(sheet, "A2", "Rapor Para Birimi:")
	f.SetCellValue(sheet, "B2", baseInfo.Label)
	f.SetCellValue(sheet, "C2", fmt.Sprintf("%.4f₺", targetPrice))

	f.SetCellValue(sheet, "A3", "Toplam İşlem:")
	f.SetCellValue(sheet, "B3", len(list))

	headers := []string{
		"ID", "Tarih", "Varlık", "İşlem Tipi",
		"Miktar", "Birim Fiyat",
		fmt.Sprintf("Toplam (%s)", baseInfo.Unit),
		"Açıklama",
	}

	for i, h := range headers {
		cell := fmt.Sprintf("%c5", 'A'+i)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	if err := f.AutoFilter(sheet, "A5:H5", nil); err != nil {
		return nil, err
	}

	f.SetColWidth(sheet, "A", "A", 8)
	f.SetColWidth(sheet, "B", "B", 18)
	f.SetColWidth(sheet, "C", "C", 22)
	f.SetColWidth(sheet, "D", "D", 14)
	f.SetColWidth(sheet, "E", "F", 16)
	f.SetColWidth(sheet, "G", "G", 18)
	f.SetColWidth(sheet, "H", "H", 30)

	totalBirim := 0.0
	totalSum := 0.0

	for i, tx := range list {
		row := i + 6
		info := assetInfo(string(tx.Asset))

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), tx.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), tx.TransactionDate)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), info.Label)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), typeDegistir(string(tx.Type)))

		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("%.4f %s", tx.Amount, info.Unit))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), tx.Price)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), tx.TotalPriceByTargetAsset)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), tx.Description)

		f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("G%d", row), moneyStyle)

		if tx.Type == model.TypeAdd {
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), addStyle)
			totalSum += tx.TotalPriceByTargetAsset
		} else {
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), subStyle)
			totalSum -= tx.TotalPriceByTargetAsset
		}

		totalBirim += tx.TotalPriceByTargetAsset
	}

	sumRow := len(list) + 7
	f.SetCellValue(sheet, fmt.Sprintf("F%d", sumRow), "İŞLEM YAPILAN TOPLAM BİRİM")
	f.SetCellValue(sheet, fmt.Sprintf("G%d", sumRow), totalBirim)
	f.SetCellStyle(sheet, fmt.Sprintf("F%d", sumRow), fmt.Sprintf("G%d", sumRow), headerStyle)

	totalRow := len(list) + 8
	f.SetCellValue(sheet, fmt.Sprintf("F%d", totalRow), "GENEL TOPLAM")
	f.SetCellValue(sheet, fmt.Sprintf("G%d", totalRow), totalSum)
	f.SetCellStyle(sheet, fmt.Sprintf("F%d", totalRow), fmt.Sprintf("G%d", totalRow), headerStyle)

	if err := f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		YSplit:      5,
		TopLeftCell: "A6",
		ActivePane:  "bottomLeft",
	}); err != nil {
		return nil, err
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
