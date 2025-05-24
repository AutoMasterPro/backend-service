package handlers

import (
	"backend-service/internal/services"
	"backend-service/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/rs/zerolog"
	"time"
)

type Handler struct {
	log        zerolog.Logger
	services   *services.Service
	jwtService *jwt.Service
}

func NewHandler(log zerolog.Logger, services *services.Service, jwtService *jwt.Service) *Handler {
	return &Handler{
		log:        log,
		services:   services,
		jwtService: jwtService,
	}
}

func (h *Handler) InitRoutes(port string) {
	app := fiber.New(fiber.Config{
		DisableDefaultContentType: true,
		CaseSensitive:             false,
	})

	app.Use(cors.New())

	// 3 requests per 10 seconds max
	app.Use(limiter.New(limiter.Config{
		Expiration: 1 * time.Second,
		Max:        10,
	}))

	api := app.Group("/tss/api/v1")
	{
		api.Get("/", func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok"})
		})

		api.Get("/metrics", monitor.New())

		// auth
		auth := api.Group("/auth")
		{
			auth.Use(limiter.New(limiter.Config{
				Expiration: 1 * time.Second,
				Max:        1,
			}))

			auth.Post("/register", h.register)
			auth.Post("/login", h.login)
		}

		userProfile := api.Group("/profile")
		{
			userProfile.Get("/", h.middlewareAuth, h.getProfile)
		}

		serv := api.Group("/services")
		{
			serv.Use(h.middlewareAuth)

			serv.Get("/", h.getServices)
			serv.Post("/", h.createService)
			//serv.Get("/:id", h.getOne)
			serv.Put("/:id", h.updateService)
			serv.Delete("/:id", h.deleteService)
		}

		vehicles := api.Group("/vehicles")
		{
			vehicles.Use(h.middlewareAuth)

			vehicles.Post("/", h.createVehicle)
			vehicles.Get("/", h.getVehicles)
			vehicles.Get("/:id", h.getVehicle)
			vehicles.Put("/:id", h.updateVehicle)
			vehicles.Delete("/:id", h.deleteVehicle)
		}

		appointments := api.Group("/appointments")
		{
			appointments.Use(h.middlewareAuth)

			appointments.Post("/", h.createAppointment)
			appointments.Get("/", h.getAppointments)
			appointments.Get("/:id", h.getAppointment)
			appointments.Put("/:id", h.updateAppointment)
			appointments.Post("/:id/cancel", h.cancelAppointment)
		}
	}

	h.log.Info().Msg("Starting server on port " + port)
	err := app.Listen(":" + port)
	if err != nil {
		h.log.Error().Msg(err.Error())
	}
}
