package handler

import (
	"context"
	"errors"
	"portfolyo/internal/document"
	"portfolyo/internal/infrastructure/app"
	"portfolyo/internal/infrastructure/errorsx"
	"portfolyo/internal/model"
	"portfolyo/internal/service"
	"strconv"
)

type UserAssetsHandler struct {
	uas service.UserAssetsService
}

func NewUserAssetsHandler(uas service.UserAssetsService) *UserAssetsHandler {
	return &UserAssetsHandler{uas: uas}
}

func (h *UserAssetsHandler) GetUserAssets(c *app.Ctx) errorsx.APIError {
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

	resp, targetPrice, err := h.uas.GetUserAssets(context.Background(), tokenID, target)

	if err != nil {
		return errorsx.DatabaseError(err)
	}

	message := "User assets retrieved successfully! Hedef Kur (" + resp.Currency + "): " +
		strconv.FormatFloat(targetPrice, 'g', 5, 64) + "₺"

	return c.SuccessResponse(resp, len(resp.Assets), message)
}

func (h *UserAssetsHandler) GetUserAsset(c *app.Ctx) errorsx.APIError {
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

	asset := c.Params("asset")

	assetType, err := model.IsValidAssetType(asset)
	if err != nil {
		return errorsx.UnauthorizedError(err)
	}

	resp, targetPrice, err := h.uas.GetUserAsset(context.Background(), tokenID, target, assetType)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return errorsx.NotFoundError(err)
		}
		return errorsx.DatabaseError(err)
	}

	message := "User asset retrieved successfully! Hedef Kur  (" + resp.TargetCurrency + "): " +
		strconv.FormatFloat(targetPrice, 'g', 5, 64) + "₺"

	return c.SuccessResponse(resp, 1, message)
}

func (h *UserAssetsHandler) GetUserAssetsPDF(c *app.Ctx) errorsx.APIError {
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

	pdfVM, targetPrice, err := h.uas.GenerateUserAssetsPDF(context.Background(), tokenID, target)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

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
