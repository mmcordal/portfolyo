package document

import (
	"bytes"
	"fmt"
	"portfolyo/internal/model"
	"portfolyo/internal/viewmodel"

	"github.com/jung-kurt/gofpdf"
)

type assetMeta struct {
	Label string
	Unit  string
}

func initPDF() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("dejavu", "", "assets/fonts/DejaVuSans.ttf")
	pdf.AddUTF8Font("dejavu", "B", "assets/fonts/DejaVuSans-Bold.ttf")
	pdf.SetFont("dejavu", "", 12)
	pdf.AddPage()
	return pdf
}

func GeneratePortfolioPDF(vm *viewmodel.PortfolioReportPDFVM, targetPrice float64) ([]byte, error) {
	pdf := initPDF()

	pdf.SetFont("dejavu", "B", 16)
	pdf.Cell(0, 10, "Genel Portföy Dekontu")
	pdf.Ln(12)

	pdf.SetFont("dejavu", "", 12)
	pdf.Cell(0, 7, fmt.Sprintf("Kullanıcı: %s", vm.NameAndSurname))
	pdf.Ln(6)
	pdf.Cell(0, 7, fmt.Sprintf("Rapor Tarihi: %s", vm.ReportDate.Format("02.01.2006 15:04")))
	pdf.Ln(6)

	base := assetInfo(vm.BaseCurrency)

	if base.Label == "Türk Lirası" {
		pdf.Cell(0, 7, fmt.Sprintf("Rapor Para Birimi: %s", base.Label))
	} else {
		pdf.Cell(0, 7, fmt.Sprintf("Rapor Para Birimi: %s = %.4f₺", base.Label, targetPrice))
	}
	pdf.Ln(10)

	wAsset, wAmount, wUnit, wUnitPrice, wTotal := 55.0, 35.0, 30.0, 35.0, 35.0

	pdf.SetFont("dejavu", "B", 11)
	pdf.CellFormat(wAsset, 8, "Varlık", "1", 0, "C", false, 0, "")
	pdf.CellFormat(wAmount, 8, "Miktar", "1", 0, "C", false, 0, "")
	pdf.CellFormat(wUnit, 8, "Birim", "1", 0, "C", false, 0, "")
	pdf.CellFormat(wUnitPrice, 8, "Birim Değer", "1", 0, "C", false, 0, "")
	pdf.CellFormat(wTotal, 8, fmt.Sprintf("Toplam (%s)", base.Unit), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("dejavu", "", 11)
	for _, row := range vm.Rows {
		info := assetInfo(string(row.AssetName))

		pdf.CellFormat(wAsset, 8, info.Label, "1", 0, "", false, 0, "")
		pdf.CellFormat(wAmount, 8, fmt.Sprintf("%.4f", row.Amount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(wUnit, 8, info.Unit, "1", 0, "C", false, 0, "")
		pdf.CellFormat(wUnitPrice, 8, fmt.Sprintf("%.4f", row.UnitPrice), "1", 0, "C", false, 0, "")
		pdf.CellFormat(wTotal, 8, fmt.Sprintf("%.4f", row.TotalPrice), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(6)
	pdf.SetFont("dejavu", "B", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Toplam Portföy Değeri: %.4f %s", vm.TotalValue, base.Unit))

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

func GenerateTransactionPDF(vm *viewmodel.TransactionReceiptPDFVM, targetPrice float64) ([]byte, error) {
	pdf := initPDF()

	pdf.SetFont("dejavu", "B", 16)
	pdf.Cell(0, 10, "İşlem Dekontu")
	pdf.Ln(12)

	pdf.SetFont("dejavu", "", 12)
	pdf.Cell(0, 7, fmt.Sprintf("Kullanıcı: %s", vm.NameAndSurname))
	pdf.Ln(6)
	pdf.Cell(0, 7, fmt.Sprintf("İşlem Tarihi: %s", vm.TransactionDate.Format("02.01.2006 15:04")))
	pdf.Ln(6)
	base := assetInfo(vm.BaseCurrency)

	pdf.SetFont("dejavu", "", 12)
	if base.Label == "Türk Lirası" {
		pdf.Cell(0, 7, fmt.Sprintf("Para Birimi: %s", base.Label))
	} else {
		pdf.Cell(0, 7, fmt.Sprintf("Para Birimi: %s = %.4f₺", base.Label, targetPrice))
	}
	pdf.Ln(10)

	info := assetInfo(string(vm.AssetName))

	pdf.SetFont("dejavu", "B", 11)
	pdf.CellFormat(60, 8, "Alan", "1", 0, "C", false, 0, "")
	pdf.CellFormat(130, 8, "Değer", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("dejavu", "", 11)

	row := func(label, val string) {
		pdf.CellFormat(60, 8, label, "1", 0, "", false, 0, "")
		pdf.CellFormat(130, 8, val, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	row("Varlık", info.Label)
	row("Tip", typeDegistir(vm.Type))
	row("Miktar", fmt.Sprintf("%.4f %s", vm.Amount, info.Unit))
	row("Birim Fiyat", fmt.Sprintf("%.4f", vm.UnitPrice))
	row("Toplam", fmt.Sprintf("%.4f %s", vm.TotalPrice, base.Unit))

	if vm.Description != "" {
		pdf.Ln(6)
		pdf.MultiCell(0, 7, "Açıklama: "+vm.Description, "", "", false)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

func assetInfo(a string) assetMeta {
	switch model.AssetType(a) {
	case model.AssetTypeTurkLirasi:
		return assetMeta{"Türk Lirası", "TL"}
	case model.AssetTypeDolar:
		return assetMeta{"Amerikan Doları", "USD"}
	case model.AssetTypeEuro:
		return assetMeta{"Euro", "EUR"}
	case model.AssetTypeSterlin:
		return assetMeta{"İngiliz Sterlini", "GBP"}
	case model.AssetTypeFrank:
		return assetMeta{"İsviçre Frangı", "CHF"}
	case model.AssetTypeGramAltin24Ayar:
		return assetMeta{"Gram Altın (24 Ayar)", "gr"}
	case model.AssetTypeGramAltin22Ayar:
		return assetMeta{"22 Ayar Altın", "gr"}
	case model.AssetTypeGramAltin18Ayar:
		return assetMeta{"18 Ayar Altın", "gr"}
	case model.AssetTypeGramAltin14Ayar:
		return assetMeta{"14 Ayar Altın", "gr"}
	case model.AssetTypeBilezik22Ayar:
		return assetMeta{"22 Ayar Bilezik", "gr"}
	case model.AssetTypeCeyrekAltin:
		return assetMeta{"Çeyrek Altın", "adet"}
	case model.AssetTypeYarimAltin:
		return assetMeta{"Yarım Altın", "adet"}
	case model.AssetTypeTamAltin:
		return assetMeta{"Tam Altın", "adet"}
	case model.AssetTypeCumhuriyetAltini:
		return assetMeta{"Cumhuriyet Altını", "adet"}
	case model.AssetTypeGumus:
		return assetMeta{"Gümüş", "gr"}
	default:
		return assetMeta{a, ""}
	}
}

func typeDegistir(a string) string {
	switch model.TypeAction(a) {
	case model.TypeSubtract:
		return "Çıkarma"
	case model.TypeAdd:
		return "Ekleme"
	default:
		return ""
	}
}
