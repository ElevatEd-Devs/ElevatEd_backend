package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SetRouter(conn *pgx.Conn) {
	app := fiber.New()
	SetGresTestRouter(app, conn)
	app.Listen(":3000")
}
