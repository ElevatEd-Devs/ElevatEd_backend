package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SetRouter(conn *pgx.Conn) {
	app := fiber.New()
	SetGresTestRouter(app, conn)
	SetAuthRouter(app, conn)
	SetEventsRouter(app, conn)
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
