package router

import (
	"elevated_backend/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SetGradeRouter(app *fiber.App, conn *pgx.Conn) {
	app.Post("v1/grade/", func(c *fiber.Ctx) error {
		return handler.GradeHandler(c, conn)
	})
}
