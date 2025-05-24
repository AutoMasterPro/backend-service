package handlers

import (
	"backend-service/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) createVehicle(c *fiber.Ctx) error {
	var input entity.VehicleCreate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing request body",
		})
	}

	if err := input.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userID, err := uuid.Parse(c.Locals("UID").(string))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing user id",
		})
	}

	vehicleID, err := h.services.VehicleService.Create(c.Context(), userID, &input)
	if err != nil {
		h.log.Error().Err(err).Msg("error creating vehicle")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "ok",
		"details": fiber.Map{
			"id": vehicleID,
		},
	})
}

func (h *Handler) getVehicles(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("UID").(string))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing user id",
		})
	}

	vehicles, err := h.services.VehicleService.GetByUserId(c.Context(), userID)
	if err != nil {
		h.log.Error().Err(err).Msg("error getting vehicles")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": vehicles,
	})
}

func (h *Handler) getVehicle(c *fiber.Ctx) error {
	vehicleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing vehicle id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing vehicle id",
		})
	}

	vehicle, err := h.services.VehicleService.GetById(c.Context(), vehicleID)
	if err != nil {
		h.log.Error().Err(err).Msg("error getting vehicle")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Check if the vehicle belongs to the requesting user
	userID, err := uuid.Parse(c.Locals("UID").(string))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing user id",
		})
	}

	if vehicle.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": vehicle,
	})
}

func (h *Handler) updateVehicle(c *fiber.Ctx) error {
	vehicleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing vehicle id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing vehicle id",
		})
	}

	// Check if the vehicle exists and belongs to the user
	vehicle, err := h.services.VehicleService.GetById(c.Context(), vehicleID)
	if err != nil {
		h.log.Error().Err(err).Msg("error getting vehicle")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userID, err := uuid.Parse(c.Locals("UID").(string))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing user id",
		})
	}

	if vehicle.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	var input entity.VehicleUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing request body",
		})
	}

	if err := input.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.services.VehicleService.Update(c.Context(), vehicleID, &input); err != nil {
		h.log.Error().Err(err).Msg("error updating vehicle")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}

func (h *Handler) deleteVehicle(c *fiber.Ctx) error {
	vehicleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing vehicle id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing vehicle id",
		})
	}

	// Check if the vehicle exists and belongs to the user
	vehicle, err := h.services.VehicleService.GetById(c.Context(), vehicleID)
	if err != nil {
		h.log.Error().Err(err).Msg("error getting vehicle")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userID, err := uuid.Parse(c.Locals("UID").(string))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing user id",
		})
	}

	if vehicle.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	if err := h.services.VehicleService.Delete(c.Context(), vehicleID); err != nil {
		h.log.Error().Err(err).Msg("error deleting vehicle")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}
