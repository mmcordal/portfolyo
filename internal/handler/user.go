package handler

import (
	"context"
	"errors"
	"portfolyo/internal/infrastructure/app"
	"portfolyo/internal/infrastructure/errorsx"
	"portfolyo/internal/service"
	"portfolyo/internal/viewmodel"
)

type AuthHandler struct {
	as service.AuthService
}

func NewAuthHandler(as service.AuthService) *AuthHandler {
	return &AuthHandler{as: as}
}

func (h *AuthHandler) Register(c *app.Ctx) errorsx.APIError {
	var input viewmodel.RegisterRequest
	if errs := c.BodyParseValidate(&input); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}
	err := h.as.Register(context.Background(), &input)
	if err != nil {
		return errorsx.DatabaseError(err)
	}
	return c.SuccessResponse("", 0, "User registered successfully!")
}

func (h *AuthHandler) Login(c *app.Ctx) errorsx.APIError {
	var input viewmodel.LoginRequest
	if errs := c.BodyParseValidate(&input); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	resp, err := h.as.Login(context.Background(), input)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	return c.SuccessResponse(resp, 1, "User logged in successfully!")
}

func (h *AuthHandler) GetUserProfile(c *app.Ctx) errorsx.APIError {
	tokenEmail, ok := c.Locals("email").(string)
	if !ok || tokenEmail == "" {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	resp, err := h.as.GetUserProfile(context.Background(), tokenEmail)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	return c.SuccessResponse(resp, 1, "User profile fetched successfully!")
}

func (h *AuthHandler) UpdateUser(c *app.Ctx) errorsx.APIError {
	tokenEmail, ok := c.Locals("email").(string)
	if !ok || tokenEmail == "" {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	var input viewmodel.UpdateRequest
	if errs := c.BodyParseValidate(&input); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	err := h.as.UpdateUser(context.Background(), tokenEmail, &input)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	return c.SuccessResponse("", 0, "User updated successfully!")
}

func (h *AuthHandler) DeleteUser(c *app.Ctx) errorsx.APIError {
	tokenEmail, ok := c.Locals("email").(string)
	if !ok || tokenEmail == "" {
		return errorsx.UnauthorizedError(errors.New("unauthorized"))
	}

	err := h.as.DeleteUser(context.Background(), tokenEmail)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	return c.SuccessResponse("", 0, "User deleted successfully!")
}
