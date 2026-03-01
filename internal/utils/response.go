package utils

import "github.com/gofiber/fiber/v2"

type Response[T any] struct {
	Data T `json:"response"`
}

func Success[T any](c *fiber.Ctx, status int, data T) error {
	return c.Status(status).JSON(Response[T]{
		Data: data,
	})
}
