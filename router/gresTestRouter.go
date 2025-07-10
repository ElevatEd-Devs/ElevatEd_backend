package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SetGresTestRouter(app *fiber.App, conn *pgx.Conn) {
	app.Get("/pingServer", func(c *fiber.Ctx) error {
		err := conn.Ping(c.Context())

		if err != nil {
			return c.JSON(fiber.Map{
				"message": "no connection",
			})
		}

		return c.JSON(fiber.Map{
			"message": "connection found",
		})
	})
}
