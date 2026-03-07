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
	"strconv"
)

type TransactionHandler struct {
	ts service.TransactionService
}

func NewTransactionHandler(ts service.TransactionService) *TransactionHandler {
	return &TransactionHandler{ts: ts}
}

func (h *TransactionHandler) AddTransaction(c *app.Ctx) errorsx.APIError {
	var input viewmodel.TransactionRequest
	if errs := c.BodyParseValidate(&input); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	tokenID, ok := c.Locals("user_id").(int64)
	if !ok || tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	err := h.ts.TransactionAdd(context.Background(), &input, tokenID)
	if err != nil {
		return errorsx.DatabaseError(err)
	}
	return c.SuccessResponse("", 0, "Transaction added successfully")
}

func (h *TransactionHandler) GetAllTransactionByAsset(c *app.Ctx) errorsx.APIError {
	tokenID, ok := c.Locals("user_id").(int64)
	if !ok || tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	assetType, err := model.IsValidAssetType(c.Params("asset"))
	if err != nil {
		return errorsx.UnauthorizedError(err)
	}

	currency := c.Get("X-Currency", "TRY")
	target, err := model.IsValidAssetType(currency)
	if err != nil {
		return errorsx.ValidationError([]error{
			errors.New("X-Currency is required"),
		})
	}

	resp, targetPrice, err := h.ts.GetAllTransactionByAsset(context.Background(), tokenID, target, assetType)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	if len(resp) == 0 {
		return c.SuccessResponse(resp, len(resp), "No transactions found for the specified asset")
	}

	message := "Transactions retrieved successfully! Hedef Kur (" + resp[0].NowTargetCurrency + "): " +
		strconv.FormatFloat(targetPrice, 'g', 5, 64) + "₺"

	return c.SuccessResponse(resp, len(resp), message)
}

func (h *TransactionHandler) GetAllTransaction(c *app.Ctx) errorsx.APIError {
	tokenID, ok := c.Locals("user_id").(int64)
	if !ok || tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	currency := c.Get("X-Currency", "TRY")
	target, err := model.IsValidAssetType(currency)
	if err != nil {
		return errorsx.ValidationError([]error{
			errors.New("X-Currency is required"),
		})
	}

	resp, targetPrice, err := h.ts.GetAllTransaction(context.Background(), tokenID, target)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	if len(resp) == 0 {
		return c.SuccessResponse(resp, len(resp), "No transactions found")
	}

	message := "Transactions retrieved successfully! Hedef Kur (" + resp[0].NowTargetCurrency + "): " +
		strconv.FormatFloat(targetPrice, 'g', 5, 64) + "₺"

	return c.SuccessResponse(resp, len(resp), message)
}

func (h *TransactionHandler) GetTransactionPDF(c *app.Ctx) errorsx.APIError {
	txID, err := strconv.ParseInt(c.Params("tx_id"), 10, 64)
	if err != nil || txID <= 0 {
		return errorsx.ValidationError([]error{
			errors.New("transaction_id must be a positive integer"),
		})
	}

	tokenID, ok := c.Locals("user_id").(int64)
	if !ok || tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	currency := c.Get("X-Currency", "TRY")
	target, err := model.IsValidAssetType(currency)
	if err != nil {
		return errorsx.ValidationError([]error{
			errors.New("X-Currency is required"),
		})
	}

	pdfVM, err := h.ts.GetTransactionPDF(context.Background(), tokenID, txID, target)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return errorsx.NotFoundError(err)
		}
		return errorsx.DatabaseError(err)
	}

	pdfBytes, err := document.GenerateTransactionPDF(pdfVM)
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

func (h *TransactionHandler) GetAllTransactionExcel(c *app.Ctx) errorsx.APIError {
	tokenID, ok := c.Locals("user_id").(int64)
	if !ok || tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	currency := c.Get("X-Currency", "TRY")
	target, err := model.IsValidAssetType(currency)
	if err != nil {
		return errorsx.ValidationError([]error{
			errors.New("X-Currency is required"),
		})
	}

	resp, targetPrice, err := h.ts.GetAllTransaction(context.Background(), tokenID, target)
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
