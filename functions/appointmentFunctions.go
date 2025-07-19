package functions

import (
	"elevated_backend/structs"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type AppointmentFunc struct{}

func (*AppointmentFunc) BuildGetAppointmentString(personType string, appointmentId int) string {
	appointmentString := fmt.Sprintf("SELECT * FROM APPOINTMENTS WHERE %v_id = %v and cancellation_reason != NULL", personType, appointmentId)
	return appointmentString
}

func (*AppointmentFunc) GetAppointment(c *fiber.Ctx, conn *pgx.Conn, queryString string) ([]structs.Appointment, error) {
	rows, err := conn.Query(c.Context(), queryString)

	if err != nil {
		return nil, err
	}

	var appointments []structs.Appointment
	for rows.Next() {
		var appointment structs.Appointment
		scanErr := rows.Scan(appointment.Id, appointment.CourseId, appointment.OfficeHoursId, appointment.TeacherId,
			appointment.StudentId, appointment.Title, appointment.Description, appointment.StartTime, appointment.EndTime,
			appointment.Location, appointment.MeetingUrl, appointment.Status, appointment.CancellationReason,
			appointment.ReminderSent, appointment.Notes, appointment.CreatedAt, appointment.CancelledAt)

		if scanErr != nil {
			return nil, scanErr
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func BuildCreateAppointmentString(appointment *structs.Appointment) string {
	creationString := fmt.Sprintf(`INSERT INTO APPOINTMENTS
					  (id, course_id, office_hours_id, teacher_id, student_id, title, description, start_time, end_time,
					  location, meeting_url, status, cancellation_reason, reminder_sent, notes, created_at, cancelled_at)
					  VALUES (%v, %v, %v, %v, %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', %v, '%v', '%v', '%v')`,
		appointment.Id, appointment.CourseId, appointment.OfficeHoursId, appointment.TeacherId, appointment.StudentId, appointment.Title,
		appointment.Description, appointment.StartTime, appointment.EndTime, appointment.Location, appointment.MeetingUrl, appointment.Status,
		appointment.CancellationReason, appointment.ReminderSent, appointment.Notes, appointment.CreatedAt, appointment.CancelledAt)
	return creationString
}

func BuildAppointmentCheckerString(appointment *structs.Appointment) string {
	checkerString := fmt.Sprintf(`SELECT id FROM Appointments WHERE id = %v or NOT ((start_time <= '%v' and end_time <= '%v') or ('%v' <= start_time and end_time <= '%v')
					  or ('%v' <= start_time and '%v' <= end_time) or (start_time <= '%v' and '%v' <= end_time)) and cancelled_at = NULL`,
		appointment.Id, appointment.StartTime, appointment.EndTime, appointment.StartTime, appointment.EndTime,
		appointment.StartTime, appointment.EndTime, appointment.StartTime, appointment.EndTime)
	return checkerString
}

func isValidAppointment(c *fiber.Ctx, conn *pgx.Conn, checkerString string) error {
	var appointmentId int
	err := conn.QueryRow(c.Context(), checkerString).Scan(&appointmentId)

	if err != nil {
		return err
	}

	return nil
}

func (*AppointmentFunc) CreateAppointment(c *fiber.Ctx, conn *pgx.Conn, appointment *structs.Appointment) error {
	createAppointmentCheckerString := BuildAppointmentCheckerString(appointment)
	err := isValidAppointment(c, conn, createAppointmentCheckerString)

	if err != nil {
		return err
	}

	_, creationErr := conn.Exec(c.Context(), BuildCreateAppointmentString(appointment))

	return creationErr
}

func BuildAppointmentUpdateString(appointmentPatcher *structs.AppointmentPatcher) string {
	switch appointmentPatcher.NewContent.(type) {
	case int:
		return fmt.Sprintf(`UPDATE appointments SET %v=%v WHERE id=%v`, appointmentPatcher.PatchField, appointmentPatcher.NewContent, appointmentPatcher.AppointmentId)
	case bool:
		return fmt.Sprintf(`UPDATE appointments SET %v=%v WHERE id=%v`, appointmentPatcher.PatchField, appointmentPatcher.NewContent, appointmentPatcher.AppointmentId)
	default:
		return fmt.Sprintf(`UPDATE appointments SET %v='%v' WHERE id=%v`, appointmentPatcher.PatchField, appointmentPatcher.NewContent, appointmentPatcher.AppointmentId)
	}
}

func BuildAppointmentExistenceCheckerString(appointmentId int) string {
	existenceString := fmt.Sprintf(`SELECT id FROM appointments WHERE id=%v`, appointmentId)
	return existenceString
}

func (*AppointmentFunc) PatchAppointment(c *fiber.Ctx, conn *pgx.Conn, appointmentPatcher *structs.AppointmentPatcher) error {
	appointmentCheckerString := BuildAppointmentExistenceCheckerString(appointmentPatcher.AppointmentId)
	err := isValidAppointment(c, conn, appointmentCheckerString)

	if err != nil {
		return err
	}

	_, patchErr := conn.Exec(c.Context(), BuildAppointmentUpdateString(appointmentPatcher))

	return patchErr
}

func BuildAppointmentDeletionString(appointmentId int) string {
	deletionString := fmt.Sprintf(`DELETE FROM appointments WHERE id=%v`, appointmentId)
	return deletionString
}

func (*AppointmentFunc) DeleteAppointment(c *fiber.Ctx, conn *pgx.Conn, appointmentDeleter *structs.AppointmentDeleter) error {
	deletionString := BuildAppointmentDeletionString(appointmentDeleter.AppointmentId)
	_, deletionErr := conn.Exec(c.Context(), deletionString)
	return deletionErr
}

func (*AppointmentFunc) BuildAppointmentError(reason string) fiber.Map {
	return fiber.Map{
		"status":       "failed",
		"message":      reason,
		"appointments": "",
	}
}
