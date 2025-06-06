package handlers

import (
	"backend-service/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) createAppointment(c *fiber.Ctx) error {
	var input entity.AppointmentCreate
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

	appointmentID, err := h.services.AppointmentService.Create(c.Context(), userID, &input)
	if err != nil {
		h.log.Error().Err(err).Msg("error creating appointment")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "ok",
		"details": fiber.Map{
			"id": appointmentID,
		},
	})
}

func (h *Handler) getAppointments(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("UID").(string))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing user id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing user id",
		})
	}

	appointments, err := h.services.AppointmentService.GetByUserId(c.Context(), userID)
	if err != nil {
		h.log.Error().Err(err).Msg("error getting appointments")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": appointments,
	})
}

func (h *Handler) getAppointment(c *fiber.Ctx) error {
	appointmentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing appointment id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing appointment id",
		})
	}

	appointment, err := h.services.AppointmentService.GetById(c.Context(), appointmentID)
	if err != nil {
		h.log.Error().Err(err).Msg("error getting appointment")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	//
	//// Check if the appointment belongs to the requesting user
	//userID, err := uuid.Parse(c.Locals("UID").(string))
	//if err != nil {
	//	h.log.Error().Err(err).Msg("error parsing user id")
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	//		"message": "error parsing user id",
	//	})
	//}
	//
	//if appointment.UserID != userID {
	//	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	//		"message": "forbidden",
	//	})
	//}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
		"details": appointment,
	})
}

func (h *Handler) updateAppointment(c *fiber.Ctx) error {
	appointmentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing appointment id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing appointment id",
		})
	}

	// Check if the appointment exists and belongs to the user
	//appointment, err := h.services.AppointmentService.GetById(c.Context(), appointmentID)
	//if err != nil {
	//	h.log.Error().Err(err).Msg("error getting appointment")
	//	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	//		"message": err.Error(),
	//	})
	//}

	//userID, err := uuid.Parse(c.Locals("UID").(string))
	//if err != nil {
	//	h.log.Error().Err(err).Msg("error parsing user id")
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	//		"message": "error parsing user id",
	//	})
	//}

	//if appointment.UserID != userID {
	//	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	//		"message": "forbidden",
	//	})
	//}

	var input entity.AppointmentUpdate
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

	if err := h.services.AppointmentService.Update(c.Context(), appointmentID, &input); err != nil {
		h.log.Error().Err(err).Msg("error updating appointment")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}

func (h *Handler) cancelAppointment(c *fiber.Ctx) error {
	appointmentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.log.Error().Err(err).Msg("error parsing appointment id")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing appointment id",
		})
	}

	// Check if the appointment exists and belongs to the user
	appointment, err := h.services.AppointmentService.GetById(c.Context(), appointmentID)
	if err != nil {
		h.log.Error().Err(err).Msg("error getting appointment")
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

	if appointment.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	if err := h.services.AppointmentService.Cancel(c.Context(), appointmentID); err != nil {
		h.log.Error().Err(err).Msg("error cancelling appointment")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}
