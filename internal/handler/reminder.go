package handler

import (
	"context"
	"errors"
	"portfolyo/internal/infrastructure/app"
	"portfolyo/internal/infrastructure/errorsx"
	"portfolyo/internal/service"
	"portfolyo/internal/viewmodel"
)

type ReminderHandler struct {
	rs service.ReminderService
}

func NewReminderHandler(rs service.ReminderService) *ReminderHandler {
	return &ReminderHandler{rs: rs}
}

func (h *ReminderHandler) Create(c *app.Ctx) errorsx.APIError {
	tokenID := c.Locals("user_id").(int64)
	if tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}
	var input viewmodel.ReminderRequest
	if errs := c.BodyParseValidate(&input); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	err := h.rs.Create(context.Background(), &input, tokenID)
	if err != nil {
		return errorsx.DatabaseError(err)
	}
	return c.SuccessResponse("", 0, "Reminder registered successfully!")
}

func (h *ReminderHandler) Delete(c *app.Ctx) errorsx.APIError {
	remindeID := int64(c.QueryInt("id"))
	if remindeID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	err := h.rs.Delete(context.Background(), remindeID)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	return c.SuccessResponse("", 0, "Reminder deleted successfully!")
}

func (h *ReminderHandler) GetAll(c *app.Ctx) errorsx.APIError {
	tokenID := c.Locals("user_id").(int64)
	if tokenID == 0 {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	resp, err := h.rs.GetAll(context.Background(), tokenID)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	return c.SuccessResponse(resp, len(resp), "Reminder fetched successfully!")
}
