package handlers

import (
	"backend-service/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) register(c *fiber.Ctx) error {
	var user entity.UserRegister
	// Парсим тело запроса
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing request body",
		})
	}
	// Проверяем тело запроса
	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Регистрируем
	userId, err := h.services.AuthService.Register(c.Context(), user)
	if err != nil {
		h.log.Error().Err(err).Msg("error registering user")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Создаем токены
	accessToken, _, err := h.jwtService.GenerateTokenPair(userId.String())
	if err != nil {
		h.log.Error().Err(err).Msg("error generating tokens")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": fiber.Map{
			"access_token": accessToken,
		},
	})
}

func (h *Handler) login(c *fiber.Ctx) error {
	var user entity.UserLogin
	// Парсим тело запроса
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing request body",
		})
	}
	// Проверяем тело запроса
	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	//
	userId, err := h.services.AuthService.Login(c.Context(), user)
	if err != nil {
		h.log.Error().Err(err).Msg("error logging in")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Создаем токены
	accessToken, _, err := h.jwtService.GenerateTokenPair(userId.ID.String())
	if err != nil {
		h.log.Error().Err(err).Msg("error generating tokens")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error generating tokens",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": fiber.Map{
			"access_token": accessToken,
		},
	})
}
