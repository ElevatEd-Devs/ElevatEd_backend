package router

import (
	"elevated_backend/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SetAppointmentRouter(app *fiber.App, conn *pgx.Conn) {
	app.Get("v1/appointments", func(c *fiber.Ctx) error {
		return handler.GetAppointmentHandler(c, conn)
	})

	app.Post("v1/appointments", func(c *fiber.Ctx) error {
		return handler.CreateAppointmentHandler(c, conn)
	})

	app.Patch("v1/appointments", func(c *fiber.Ctx) error {
		return handler.UpdateAppointmentHandler(c, conn)
	})

	app.Delete("v1/appointments", func(c *fiber.Ctx) error {
		return handler.DeleteAppointmentHandler(c, conn)
	})

}
