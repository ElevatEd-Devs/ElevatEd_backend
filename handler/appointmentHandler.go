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
	errorMap, userClaims, initialErr := initialCheckerG(c, &appointmentFunc)

	if initialErr != nil {
		return c.JSON(errorMap)
	}

	appointments, getErr := appointmentFunc.GetAppointment(c, conn, &userClaims.Details)

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
	errorMap, userClaims, initialErr := initialCheckerPPD(c, &appointment, &appointmentFunc)

	if initialErr != nil {
		return c.JSON(errorMap)
	}

	creationErr := appointmentFunc.CreateAppointment(c, conn, &appointment, &userClaims)

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
	errorMap, userClaims, initialErr := initialCheckerPPD(c, &appointmentPatcher, &appointmentFunc)

	if initialErr != nil {
		return c.JSON(errorMap)
	}

	patchErr := appointmentFunc.PatchAppointment(c, conn, &appointmentPatcher, &userClaims)

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
	errorMap, userClaims, initialErr := initialCheckerPPD(c, &appointmentDeleter, &appointmentFunc)

	if initialErr != nil {
		return c.JSON(errorMap)
	}

	deletionErr := appointmentFunc.DeleteAppointment(c, conn, &appointmentDeleter, &userClaims)

	if deletionErr != nil {
		return c.JSON(appointmentFunc.BuildAppointmentError("could not delete appointment"))
	}

	return c.JSON(fiber.Map{
		"status":  "successful",
		"message": "appointment was deleted",
	})
}

func initialCheckerG(c *fiber.Ctx, appointmentFunc *functions.AppointmentFunc) (fiber.Map, functions.CustomClaimStruct, error) {
	var authFunc functions.AuthFunc
	jwt := authFunc.ExtractJWTFromHeader(c)
	verified, userClaims, verificationErr := authFunc.VerifyJWT(jwt)

	if verificationErr != nil {
		return appointmentFunc.BuildAppointmentError(verificationErr.Error()), userClaims, verificationErr
	}

	if !verified {
		const errorString = "identity could not be verified"
		return appointmentFunc.BuildAppointmentError(errorString), userClaims, errors.New(errorString)
	}

	return nil, userClaims, nil
}

func initialCheckerPPD(c *fiber.Ctx, appointmentStruct any, appointmentFunc *functions.AppointmentFunc) (fiber.Map, functions.CustomClaimStruct, error) {
	err := c.BodyParser(&appointmentStruct)
	if err != nil {
		return appointmentFunc.BuildAppointmentError(err.Error()), functions.CustomClaimStruct{}, err
	}

	var authFunc functions.AuthFunc
	jwt := authFunc.ExtractJWTFromHeader(c)
	verified, userClaims, verificationErr := authFunc.VerifyJWT(jwt)

	if verificationErr != nil {
		return appointmentFunc.BuildAppointmentError(verificationErr.Error()), userClaims, verificationErr
	}

	if !verified {
		const errorString = "identity could not be verified"
		return appointmentFunc.BuildAppointmentError(errorString), userClaims, errors.New(errorString)
	}

	return nil, userClaims, nil
}
