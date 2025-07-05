package common

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type PaginationResponse struct {
	Data   any   `json:"data"`
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

func ResponseCreate(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(data)
}

func RespondCreated(c *fiber.Ctx, id uint) error {
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}
func ResponseSuccess(c *fiber.Ctx, data any) error {
	if data == nil {
		return c.Status(fiber.StatusNoContent).JSON(data)
	}
	return c.Status(fiber.StatusOK).JSON(data)
}

func ResponseContent(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}

func RespondError(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{
		"error":   http.StatusText(code),
		"message": msg,
	})
}

func ResponseNoContent(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}
