package handler

import (
	"elevated_backend/functions"
	"elevated_backend/structs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func GradeHandler(c *fiber.Ctx, conn *pgx.Conn) error {
	var authFunc functions.AuthFunc
	jwt := authFunc.ExtractJWTFromHeader(c)
	valid, userClaims, err := authFunc.VerifyJWT(jwt)

	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "failed",
			"message": "error encountered during verification",
		})
	}

	if !valid {
		return c.JSON(fiber.Map{
			"status":  "failed",
			"message": "could not verify identity",
		})
	}

	var assessment structs.Asessment
	parseErr := c.BodyParser(&assessment)
	if parseErr != nil {
		return parseErr
	}

	var gradeFunc functions.GradeFunc
	correctChoices, getErr := gradeFunc.GetAnswersForAssessment(c, conn, assessment.AssessmentId)
	if getErr != nil {
		return getErr
	}
	score, gradeErr := gradeFunc.GradeAssessment(assessment.Choices, correctChoices)

	if gradeErr != nil {
		return gradeErr
	}

	var gradeMeta structs.GradeMeta
	gradeMeta.AssessmentId = assessment.AssessmentId
	gradeMeta.Feedback = "Autograded"
	gradeMeta.GradedAt = time.Now().Local().UTC()
	gradeMeta.GradedBy = "Autograder"
	gradeMeta.PercentageScore = score / len(correctChoices) * 100
	gradeMeta.Score = score

	registrationErr := gradeFunc.RegisterGrade(c, conn, userClaims.Details.Id, &gradeMeta, assessment.AssessmentId)

	if registrationErr != nil {
		return c.JSON(fiber.Map{
			"status":  "failed",
			"message": registrationErr.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":         "success",
		"message":        "assessment has been graded",
		"grade_metadata": gradeMeta,
	})
}
