package structs

import "time"

type ClientAnswer struct {
	QuestionId int    `json:"question_id"`
	ChoiceText string `json:"choice_text"`
}

type Asessment struct {
	AssessmentId int            `json:"assessment_id"`
	Choices      []ClientAnswer `json:"choices"`
}

type GradeMeta struct {
	AssessmentId    int       `json:"assessment_id"`
	Score           int       `json:"score"`
	PercentageScore int       `json:"percentage_score"`
	Feedback        string    `json:"feedback"`
	GradedAt        time.Time `json:"graded_at"`
	GradedBy        string    `json:"graded_by"`
}
