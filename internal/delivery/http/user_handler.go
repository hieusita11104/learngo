package http

import (
	"time"

	"learngo/internal/common"
	"learngo/internal/domain"
	"learngo/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type UserHandler struct {
	Usecase *usecase.UserUsecase
}

func NewUserHandler(app *fiber.App, uc *usecase.UserUsecase) {
	handler := &UserHandler{Usecase: uc}

	app.Post("/v1/users/register", handler.Create)
	app.Post("/v1/users/login", limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		LimitReached: func(_ *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "You sent too many requests, please try again later.")
		},
	}), handler.Login)
	app.Post("/v1/users/refresh", common.AuthMiddleware, handler.Refresh)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}

	if err := user.Validate(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := h.Usecase.CreateUser(&user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return common.ResponseCreate(c, user.ID)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}

	if err := user.ValidateLogin(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	accessToken, refreshToken, err := h.Usecase.LoginUser(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return common.ResponseSuccess(c, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *UserHandler) Refresh(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}
	// Get user ID from the context
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}
	newAccessToken, err := h.Usecase.RefreshToken(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return common.ResponseSuccess(c, fiber.Map{
		"access_token": newAccessToken,
	})
}
