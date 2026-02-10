package handler

import (
	"context"
	"errors"
	"portfolyo/internal/document"
	"portfolyo/internal/infrastructure/app"
	"portfolyo/internal/infrastructure/errorsx"
	"portfolyo/internal/model"
	"portfolyo/internal/service"
	"portfolyo/internal/viewmodel"
	"strings"
)

type UserAssetsHandler struct {
	uas service.UserAssetsService
}

func NewUserAssetsHandler(uas service.UserAssetsService) *UserAssetsHandler {
	return &UserAssetsHandler{uas: uas}
}

func (h *UserAssetsHandler) UserAssetAdd(c *app.Ctx) errorsx.APIError {
	var input viewmodel.TransactionRequest
	if errs := c.BodyParseValidate(&input); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	tokenID := c.Locals("user_id").(int64)
	if tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	err := h.uas.UserAssetAdd(context.Background(), &input, tokenID)
	if err != nil {
		return errorsx.DatabaseError(err)
	}
	return c.SuccessResponse("", 0, "User assets added successfully")
}

func (h *UserAssetsHandler) GetUserAssets(c *app.Ctx) errorsx.APIError {
	tokenID := c.Locals("user_id").(int64)
	if tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	currency := c.Get("X-Currency", "TRY")

	resp, _, err := h.uas.GetUserAssets(context.Background(), tokenID, model.AssetType(currency))

	if err != nil {
		return errorsx.DatabaseError(err)
	}

	return c.SuccessResponse(resp, len(resp.Assets)+1, "User assets retrieved successfully")
}

func (h *UserAssetsHandler) GetAllTransaction(c *app.Ctx) errorsx.APIError {
	tokenID := c.Locals("user_id").(int64)
	if tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	currency := c.Get("X-Currency", "TRY")

	resp, _, err := h.uas.GetAllTransaction(context.Background(), tokenID, model.AssetType(currency))
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	return c.SuccessResponse(resp, len(resp), "User transactions retrieved successfully")
}

func (h *UserAssetsHandler) GetUserAssetsPDF(c *app.Ctx) errorsx.APIError {
	tokenID := c.Locals("user_id").(int64)
	if tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	currency := c.Get("X-Currency", "TRY")

	full := getFullname(c)

	// Service çağırıyoruz
	pdfVM, targetPrice, err := h.uas.GenerateUserAssetsPDF(context.Background(), tokenID, model.AssetType(currency), full)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	// PDF bytes üret
	pdfBytes, err := document.GeneratePortfolioPDF(pdfVM, targetPrice)
	if err != nil {
		return errorsx.InternalError(err)
	}

	c.Ctx.Response().Header.Set("Content-Type", "application/pdf")
	c.Ctx.Response().Header.Set("Content-Disposition", "attachment; filename=portfolio.pdf")

	_, err = c.Ctx.Write(pdfBytes)
	if err != nil {
		return errorsx.InternalError(err)
	}

	return nil
}

func (h *UserAssetsHandler) GetTransactionPDF(c *app.Ctx) errorsx.APIError {
	txID := c.QueryInt("transaction_id")
	if txID == 0 {
		return errorsx.ValidationError([]error{errors.New("transaction_id is required")})
	}

	tokenID := c.Locals("user_id").(int64)
	if tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	currency := c.Get("X-Currency", "USD")

	full := getFullname(c)

	// PDF VM alma
	pdfVM, targetPrice, err := h.uas.GetTransactionPDF(context.Background(), tokenID, int64(txID), model.AssetType(currency), full)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	// PDF üret
	pdfBytes, err := document.GenerateTransactionPDF(pdfVM, targetPrice)
	if err != nil {
		return errorsx.InternalError(err)
	}

	c.Ctx.Response().Header.Set("Content-Type", "application/pdf")
	c.Ctx.Response().Header.Set("Content-Disposition", "attachment; filename=transaction.pdf")
	_, err = c.Ctx.Write(pdfBytes)
	if err != nil {
		return errorsx.InternalError(err)
	}

	return nil
}

func (h *UserAssetsHandler) GetAllTransactionExcel(c *app.Ctx) errorsx.APIError {
	tokenID := c.Locals("user_id").(int64)
	if tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	currency := c.Get("X-Currency", "TRY")

	resp, targetPrice, err := h.uas.GetAllTransaction(context.Background(), tokenID, model.AssetType(currency))
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	excelBytes, err := document.GenerateTransactionsExcel(resp, targetPrice, model.AssetType(currency))
	if err != nil {
		return errorsx.InternalError(err)
	}

	c.Ctx.Response().Header.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Ctx.Response().Header.Set("Content-Disposition", "attachment; filename=transactions.xlsx")

	_, err = c.Ctx.Write(excelBytes)
	if err != nil {
		return errorsx.InternalError(err)
	}

	return nil
}

func getFullname(c *app.Ctx) string {
	name, _ := c.Locals("name").(string)
	surname, _ := c.Locals("surname").(string)

	full := name + " " + surname
	full = strings.TrimSpace(full)

	if full == "" {
		return "Kullanıcı"
	}
	return full
}
