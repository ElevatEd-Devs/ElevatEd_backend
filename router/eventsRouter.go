package router

import (
	"elevated_backend/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SetEventsRouter(app *fiber.App, conn *pgx.Conn) {
	app.Get("/v1/events", func(c *fiber.Ctx) error {
		return handler.GetEventsHandler(c, conn)
	})
}
