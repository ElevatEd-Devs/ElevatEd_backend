package router

import (
	"elevated_backend/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SetAuthRouter(app *fiber.App, conn *pgx.Conn) {
	app.Post("/users/", func(c *fiber.Ctx) error {
		return handler.SignupHandler(c, conn)
	})

	app.Post("/login/", func(c *fiber.Ctx) error {
		return handler.LoginHandler(c, conn)
	})

	app.Post("/jwtToken/", func(c *fiber.Ctx) error {
		return handler.JWTHandler(c, conn)
	})
}
