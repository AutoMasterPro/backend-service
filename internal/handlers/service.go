package handlers

import (
	"backend-service/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) createService(c *fiber.Ctx) error {
	var service entity.Service
	// Парсим тело запроса
	if err := c.BodyParser(&service); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing request body",
		})
	}
	// Проверяем тело запроса
	if err := service.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	uidRaw := c.Locals("UID").(string)
	userID, err := uuid.Parse(uidRaw)
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing user id",
		})
	}
	// Проверяем что admin
	isAdmin, err := h.services.UserRoleService.IsAdmin(c.Context(), userID)
	if err != nil {
		h.log.Error().Err(err).Msg("error checking admin")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
	// Создаем услугу
	serviceId, err := h.services.ServiceService.Create(c.Context(), &service)
	if err != nil {
		h.log.Error().Err(err).Msg("error creating service")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Возвращаем сервис
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": fiber.Map{
			"id": serviceId,
		},
	})
}

func (h *Handler) getServices(c *fiber.Ctx) error {
	// Создаем услугу
	services, err := h.services.ServiceService.GetAll(c.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("error getting service")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Возвращаем сервис
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": services,
	})
}

func (h *Handler) updateService(c *fiber.Ctx) error {
	serviceIdRaw := c.Params("id")
	serviceId, err := uuid.Parse(serviceIdRaw)
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing service id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing service id",
		})
	}

	var service entity.Service
	// Парсим тело запроса
	if err := c.BodyParser(&service); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing request body",
		})
	}
	// Проверяем тело запроса
	if err := service.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	uidRaw := c.Locals("UID").(string)
	userID, err := uuid.Parse(uidRaw)
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing user id",
		})
	}
	// Проверяем что admin
	isAdmin, err := h.services.UserRoleService.IsAdmin(c.Context(), userID)
	if err != nil {
		h.log.Error().Err(err).Msg("error checking admin")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
	service.ID = serviceId
	// Редактируем услугу
	_, err = h.services.ServiceService.Update(c.Context(), &service)
	if err != nil {
		h.log.Error().Err(err).Msg("error creating service")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Возвращаем сервис
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}

func (h *Handler) deleteService(c *fiber.Ctx) error {
	serviceIdRaw := c.Params("id")
	serviceId, err := uuid.Parse(serviceIdRaw)
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing service id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing service id",
		})
	}

	uidRaw := c.Locals("UID").(string)
	userID, err := uuid.Parse(uidRaw)
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing user id",
		})
	}
	// Проверяем что admin
	isAdmin, err := h.services.UserRoleService.IsAdmin(c.Context(), userID)
	if err != nil {
		h.log.Error().Err(err).Msg("error checking admin")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
	// Удаляем услугу
	err = h.services.ServiceService.Delete(c.Context(), serviceId)
	if err != nil {
		h.log.Error().Err(err).Msg("error deleting service")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Возвращаем сервис
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}
