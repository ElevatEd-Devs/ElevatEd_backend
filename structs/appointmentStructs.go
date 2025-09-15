package structs

import "time"

type Appointment struct {
	Id                 int       `json:"id"`
	CourseId           int       `json:"course_id"`
	OfficeHoursId      int       `json:"office_hours_id"`
	TeacherId          int       `json:"teacher_id"`
	StudentId          int       `json:"student_id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	Location           string    `json:"location"`
	MeetingUrl         string    `json:"meeting_url"`
	Status             string    `json:"Status"`
	CancellationReason string    `json:"cancellation_reason"`
	ReminderSent       bool      `json:"reminder_sent"`
	Notes              string    `json:"notes"`
	CreatedAt          time.Time `json:"created_at"`
	CancelledAt        time.Time `json:"cancelled_at"`
}

type AppointmentsRequester struct {
	RequesterId   int    `json:"requester_id"`
	RequesterType string `json:"requester_type"`
}

type AppointmentPatcher struct {
	AppointmentId int    `json:"appointment_id"`
	PatchField    string `json:"patch_field"`
	NewContent    any    `json:"new_content"`
}

type AppointmentDeleter struct {
	AppointmentId int `json:"appointment_id"`
}
