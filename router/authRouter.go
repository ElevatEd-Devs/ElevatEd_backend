package router

import (
	"elevated_backend/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SetAuthRouter(app *fiber.App, conn *pgx.Conn) {
	app.Post("v1/users/", func(c *fiber.Ctx) error {
		return handler.SignupHandler(c, conn)
	})

	app.Post("v1/login/", func(c *fiber.Ctx) error {
		return handler.LoginHandler(c, conn)
	})

	app.Post("v1/jwtToken/", func(c *fiber.Ctx) error {
		return handler.JWTHandler(c, conn)
	})
}
