package functions

import (
	"elevated_backend/structs"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type AppointmentFunc struct{}

func (*AppointmentFunc) BuildGetAppointmentString(personType string, requestId int) string {
	appointmentString := fmt.Sprintf("SELECT * FROM appointments WHERE %v_id = %v and cancellation_reason IS NOT NULL", personType, requestId)
	return appointmentString
}

func (*AppointmentFunc) GetAppointment(c *fiber.Ctx, conn *pgx.Conn, queryString string) ([]structs.Appointment, error) {
	rows, err := conn.Query(c.Context(), queryString)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []structs.Appointment
	for rows.Next() {
		var appointment structs.Appointment
		scanErr := rows.Scan(&appointment.Id, &appointment.CourseId, &appointment.OfficeHoursId, &appointment.TeacherId,
			&appointment.StudentId, &appointment.Title, &appointment.Description, &appointment.StartTime, &appointment.EndTime,
			&appointment.Location, &appointment.MeetingUrl, &appointment.Status, &appointment.CancellationReason,
			&appointment.ReminderSent, &appointment.Notes, &appointment.CreatedAt, &appointment.CancelledAt)

		if scanErr != nil {
			return nil, scanErr
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}

func BuildCreateAppointmentString() string {
	creationString := `INSERT INTO APPOINTMENTS
					  (id, course_id, office_hours_id, teacher_id, student_id, title, description, start_time, end_time,
					  location, meeting_url, status, cancellation_reason, reminder_sent, notes, created_at, cancelled_at)
					  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`
	return creationString
}

func BuildAppointmentCheckerString() string {
	checkerString := `SELECT id FROM Appointments WHERE id = $1 or NOT ((start_time <= $2 and end_time <= $3) or ($2 <= start_time and end_time <= $3)
					  or ($2 <= start_time and $3 <= end_time) or (start_time <= $2 and $3 <= end_time)) and cancelled_at = NULL`
	return checkerString
}

func isValidAppointment(c *fiber.Ctx, conn *pgx.Conn, checkerString string, appointment *structs.Appointment) (bool, error) {
	var appointmentId = -1
	err := conn.QueryRow(c.Context(), checkerString, appointment.Id, appointment.StartTime, appointment.EndTime).Scan(&appointmentId)
	if appointmentId == -1 {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return false, nil
}

func checkAppointmentExistence(c *fiber.Ctx, conn *pgx.Conn, checkerString string, appointmentPatcher *structs.AppointmentPatcher) (bool, error) {
	var appointmentId = -1
	err := conn.QueryRow(c.Context(), checkerString, appointmentPatcher.AppointmentId).Scan(&appointmentId)
	if appointmentId == -1 {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return false, nil
}

func (*AppointmentFunc) CreateAppointment(c *fiber.Ctx, conn *pgx.Conn, appointment *structs.Appointment) error {
	createAppointmentCheckerString := BuildAppointmentCheckerString()
	valid, err := isValidAppointment(c, conn, createAppointmentCheckerString, appointment)

	if err != nil {
		return err
	}

	if !valid {
		return errors.New("appointment is not unique")
	}

	createAppointmentString := BuildCreateAppointmentString()
	_, creationErr := conn.Exec(c.Context(), createAppointmentString,
		fmt.Sprintf("%d", appointment.Id), appointment.CourseId, appointment.OfficeHoursId, appointment.TeacherId, appointment.StudentId, appointment.Title,
		appointment.Description, appointment.StartTime, appointment.EndTime, appointment.Location, appointment.MeetingUrl, appointment.Status,
		appointment.CancellationReason, appointment.ReminderSent, appointment.Notes, appointment.CreatedAt, appointment.CancelledAt)
	return creationErr
}

func BuildAppointmentUpdateString(patchField string) string {
	return fmt.Sprintf(`UPDATE appointments SET %v=$1 WHERE id=$2`, patchField)
}

func BuildAppointmentExistenceCheckerString(appointmentId int) string {
	existenceString := fmt.Sprintf(`SELECT id FROM appointments WHERE id=%v`, appointmentId)
	return existenceString
}

func (*AppointmentFunc) PatchAppointment(c *fiber.Ctx, conn *pgx.Conn, appointmentPatcher *structs.AppointmentPatcher) error {
	appointmentCheckerString := BuildAppointmentExistenceCheckerString(appointmentPatcher.AppointmentId)
	exists, err := checkAppointmentExistence(c, conn, appointmentCheckerString, appointmentPatcher)

	if !exists {
		return errors.New("no appointment to edit")
	}

	if err != nil {
		return err
	}

	updateString := BuildAppointmentUpdateString(appointmentPatcher.PatchField)
	_, patchErr := conn.Exec(c.Context(), updateString, appointmentPatcher.NewContent, appointmentPatcher.AppointmentId)
	fmt.Println(patchErr)
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
