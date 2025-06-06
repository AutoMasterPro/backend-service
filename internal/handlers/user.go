package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) getProfile(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("UID").(string))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing user id",
		})
	}

	user, err := h.services.UserRoleService.GetById(c.Context(), userID)
	if err != nil {
		h.log.Error().Err(err).Msg("error getting user profile")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": user.ToProfileResponse(),
	})
}

func (h *Handler) getAllClientsWithAppointments(c *fiber.Ctx) error {
	ctx := c.Context()
	clients, err := h.services.UserRoleService.GetAllClients(ctx)
	if err != nil {
		h.log.Error().Err(err).Msg("error getting clients")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error getting clients"})
	}

	result := make([]map[string]interface{}, 0, len(clients))
	for _, client := range clients {
		appointments, err := h.services.AppointmentService.GetByUserId(ctx, client.ID)
		if err != nil {
			h.log.Error().Err(err).Msg("error getting appointments for client")
			appointments = nil
		}
		result = append(result, map[string]interface{}{
			"client":       client,
			"appointments": appointments,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": result,
	})
}
