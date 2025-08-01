package structs

import "time"

type Event struct {
	Id          int       `json:"id"`
	CourseId    int       `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location"`
	EventType   string    `json:"event_type"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type EventPatcher struct {
	Id         int    `json:"event_id"`
	PatchField string `json:"patch_field"`
	NewContent any    `json:"new_content"`
}

type EventDeleter struct {
	Id int `json:"event_id"`
}
