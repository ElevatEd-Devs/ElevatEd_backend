package functions

import (
	"elevated_backend/structs"
	"errors"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type GradeFunc struct{}

func (*GradeFunc) GradeAssessment(clientChoices []structs.ClientAnswer, correctChoices map[int]string) (int, error) {
	score := 0

	for _, choice := range clientChoices {
		val, contained := correctChoices[choice.QuestionId]

		if !contained {
			return 0, errors.New("extra question detected")
		}
		if val != choice.ChoiceText {
			continue
		}
		score += 1
	}

	return score, nil
}

func (*GradeFunc) RegisterGrade(c *fiber.Ctx, conn *pgx.Conn, userId int, gradeMeta *structs.GradeMeta, assessmentId int) error {
	writeString := `INSERT INTO assessment_submissions
	(id, assessment_id, user_id, score, percentage_score, feedback, graded_at, graded_by)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := conn.Exec(c.Context(), writeString, rand.Intn(1000), assessmentId, userId, gradeMeta.Score, gradeMeta.PercentageScore, gradeMeta.Feedback, gradeMeta.GradedAt, gradeMeta.GradedBy)
	return err
	// return nil
}

func (*GradeFunc) GetAnswersForAssessment(c *fiber.Ctx, conn *pgx.Conn, assessmentId int) (map[int]string, error) {
	queryString := `SELECT question_id, choice_text question_choices WHERE id = $1 and is_correct = TRUE`
	rows, err := conn.Query(c.Context(), queryString, assessmentId)

	if err != nil {
		return nil, err
	}

	choices := make(map[int]string)
	for rows.Next() {
		var choice structs.ClientAnswer
		scanErr := rows.Scan(&choice.QuestionId, &choice.ChoiceText)
		choices[choice.QuestionId] = choice.ChoiceText

		if scanErr != nil {
			return nil, scanErr
		}
	}

	return choices, nil

	// choices := make(map[int]string)
	// choices[1] = "ienwi920dkdK"
	// choices[2] = "ienwi920dkdK"

	// return choices, nil
}
