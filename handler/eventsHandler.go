package handler

import (
	"elevated_backend/functions"
	"elevated_backend/structs"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func GetEventsHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var eventFunc functions.EventFunc
	var authFunc functions.AuthFunc
	jwt := authFunc.ExtractJWTFromHeader(c)
	valid, userClaims, err := authFunc.VerifyJWT(jwt)

	if err != nil {
		return c.JSON(eventFunc.BuildErrorString(err.Error()))
	}

	if !valid {
		return c.JSON(eventFunc.BuildErrorString("invalid permissions"))
	}

	events, getErr := eventFunc.GetEvents(c, conn, &userClaims)

	if getErr != nil {
		return c.JSON(eventFunc.BuildErrorString(getErr.Error()))
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "events were gotten",
		"events":  events,
	})
}

func PostEventsHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var eventFunc functions.EventFunc
	var authFunc functions.AuthFunc
	jwt := authFunc.ExtractJWTFromHeader(c)
	valid, userClaims, err := authFunc.VerifyJWT(jwt)

	if err != nil {
		return c.JSON(eventFunc.BuildErrorString(err.Error()))
	}

	if !valid {
		return c.JSON(eventFunc.BuildErrorString("identity could not be verified"))
	}

	var event structs.Event
	parseErr := c.BodyParser(&event)

	if parseErr != nil {
		return c.JSON(eventFunc.BuildErrorString("malformed request"))
	}

	if userClaims.Details.Role != "teacher" {
		return c.JSON(eventFunc.BuildErrorString("invalid permissions"))
	}

	createErr := eventFunc.CreateEvent(c, conn, &event)

	if createErr != nil {
		return c.JSON(eventFunc.BuildErrorString(createErr.Error()))
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "event was saved",
	})
}

// func PatchEventHandler(c *fiber.Ctx, conn *pgx.Conn) error {
// 	var eventFunc functions.EventFunc
// 	var authFunc functions.AuthFunc
// 	jwt := authFunc.ExtractJWTFromHeader(c)
// 	valid, userClaims, err := authFunc.VerifyJWT(jwt)

// 	if err != nil {
// 		return c.JSON(eventFunc.BuildErrorString(err.Error()))
// 	}

// 	if !valid {
// 		return c.JSON(eventFunc.BuildErrorString("identity could not be verified"))
// 	}

// 	var eventPatcher structs.EventPatcher
// 	parseErr := c.BodyParser(&eventPatcher)

// 	if parseErr != nil {
// 		return c.JSON(eventFunc.BuildErrorString("malformed request"))
// 	}

// 	if userClaims.Details.Role != "teacher" {
// 		return c.JSON(eventFunc.BuildErrorString("invalid permissions"))
// 	}

// 	patchErr := eventFunc.PatchEvent(c, conn, &eventPatcher)

// 	if patchErr != nil {
// 		return c.JSON(eventFunc.BuildErrorString(patchErr.Error()))
// 	}

// 	return c.JSON(fiber.Map{
// 		"status":  "success",
// 		"message": "event was edited",
// 	})
// }

func DeleteEventHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var eventFunc functions.EventFunc
	var authFunc functions.AuthFunc
	jwt := authFunc.ExtractJWTFromHeader(c)
	valid, userClaims, err := authFunc.VerifyJWT(jwt)

	if err != nil {
		return c.JSON(eventFunc.BuildErrorString(err.Error()))
	}

	if !valid {
		return c.JSON(eventFunc.BuildErrorString("identity could not be verified"))
	}

	var eventDeleter structs.EventDeleter
	parseErr := c.BodyParser(&eventDeleter)

	if parseErr != nil {
		return c.JSON(eventFunc.BuildErrorString("malformed request"))
	}

	if userClaims.Details.Role != "teacher" {
		return c.JSON(eventFunc.BuildErrorString("invalid permissions"))
	}

	deleteErr := eventFunc.DeleteEvent(c, conn, &eventDeleter)

	if deleteErr != nil {
		return c.JSON(eventFunc.BuildErrorString(deleteErr.Error()))
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "event was deleted",
	})
}
