package handler

import (
	"elevated_backend/functions"
	"elevated_backend/structs"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func GetAppointmentHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var appointmentFunc functions.AppointmentFunc
	var appointmentsRequester structs.AppointmentsRequester
	errorMap, initialErr := initialChecker(c, &appointmentsRequester, &appointmentFunc)

	if initialErr != nil {
		return c.JSON(errorMap)
	}

	appointmentString := appointmentFunc.BuildGetAppointmentString(appointmentsRequester.RequesterType, appointmentsRequester.RequesterId)
	appointments, getErr := appointmentFunc.GetAppointment(c, conn, appointmentString)

	if getErr != nil {
		return c.JSON(appointmentFunc.BuildAppointmentError(getErr.Error()))
	}

	return c.JSON(fiber.Map{
		"status":       "successful",
		"message":      "appointments were gotten",
		"appointments": appointments,
	})
}

func CreateAppointmentHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var appointmentFunc functions.AppointmentFunc
	var appointment structs.Appointment
	errorMap, initialErr := initialChecker(c, &appointment, &appointmentFunc)

	if initialErr != nil {
		return c.JSON(errorMap)
	}

	creationErr := appointmentFunc.CreateAppointment(c, conn, &appointment)

	if creationErr != nil {
		return c.JSON(appointmentFunc.BuildAppointmentError(creationErr.Error()))
	}

	return c.JSON(fiber.Map{
		"status":  "successful",
		"message": "appointment was saved",
	})
}

func UpdateAppointmentHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var appointmentFunc functions.AppointmentFunc
	var appointmentPatcher structs.AppointmentPatcher
	errorMap, initialErr := initialChecker(c, &appointmentPatcher, &appointmentFunc)

	if initialErr != nil {
		return c.JSON(errorMap)
	}

	patchErr := appointmentFunc.PatchAppointment(c, conn, &appointmentPatcher)

	if patchErr != nil {
		return c.JSON(appointmentFunc.BuildAppointmentError(patchErr.Error()))
	}

	return c.JSON(fiber.Map{
		"status":  "successful",
		"message": "appointment was updated",
	})
}

func DeleteAppointmentHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var appointmentFunc functions.AppointmentFunc
	var appointmentDeleter structs.AppointmentDeleter
	errorMap, initialErr := initialChecker(c, &appointmentDeleter, &appointmentFunc)

	if initialErr != nil {
		return c.JSON(errorMap)
	}

	deletionErr := appointmentFunc.DeleteAppointment(c, conn, &appointmentDeleter)

	if deletionErr != nil {
		return c.JSON(appointmentFunc.BuildAppointmentError("could not delete appointment"))
	}

	return c.JSON(fiber.Map{
		"status":  "successful",
		"message": "appointment was deleted",
	})
}

func initialChecker(c *fiber.Ctx, appointmentStruct any, appointmentFunc *functions.AppointmentFunc) (fiber.Map, error) {
	err := c.BodyParser(&appointmentStruct)
	if err != nil {
		return appointmentFunc.BuildAppointmentError(err.Error()), err
	}

	var authFunc functions.AuthFunc
	jwt := authFunc.ExtractJWTFromHeader(c)
	verified, verificationErr := authFunc.VerifyJWT(jwt)

	if verificationErr != nil {
		return appointmentFunc.BuildAppointmentError(verificationErr.Error()), verificationErr
	}

	if !verified {
		const errorString = "identity could not be verified"
		return appointmentFunc.BuildAppointmentError(errorString), errors.New(errorString)
	}

	return nil, nil
}
