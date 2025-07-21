package handler

import (
	"elevated_backend/functions"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func GetEventsHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var authFunc functions.AuthFunc
	jwt := authFunc.ExtractJWTFromHeader(c)
	valid, err := authFunc.VerifyJWT(jwt)

	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	if !valid {
		return c.JSON(fiber.Map{
			"status":  "failed",
			"message": "session expired",
		})
	}

}
