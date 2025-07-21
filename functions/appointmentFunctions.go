package functions

import (
	"elevated_backend/structs"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type AppointmentFunc struct{}

func (*AppointmentFunc) GetAppointment(c *fiber.Ctx, conn *pgx.Conn, userClaimsDetails *structs.UserDetails) ([]structs.Appointment, error) {
	var queryString string
	switch userClaimsDetails.Role {
	case "student":
		queryString = `SELECT * FROM appointments WHERE student_id= $1 and cancellation_reason IS NOT NULL`
	case "teacher":
		queryString = `SELECT * FROM appointments WHERE teacher_id= $1 and cancellation_reason IS NOT NULL`
	default:
		return nil, errors.New("invalid permissions")
	}

	rows, err := conn.Query(c.Context(), queryString, userClaimsDetails.Id)

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

func doesAppointmentExist(c *fiber.Ctx, conn *pgx.Conn, appointmentPatcher *structs.AppointmentPatcher) (bool, error) {
	queryString := `SELECT EXISTS
					(SELECT 1 FROM appointments WHERE id = $1 AND cancelled_at IS NOT NULL)`
	var appointmentExists bool
	err := conn.QueryRow(c.Context(), queryString, appointmentPatcher.AppointmentId).Scan(&appointmentExists)

	if appointmentExists {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return false, nil
}

func (*AppointmentFunc) CreateAppointment(c *fiber.Ctx, conn *pgx.Conn, appointment *structs.Appointment, userClaims *CustomClaimStruct) error {
	valid, err := isValidAppointment(c, conn, appointment)

	if err != nil {
		return err
	}

	if !valid {
		return err
	}

	createAppointmentString := `INSERT INTO APPOINTMENTS
					(id, course_id, office_hours_id, teacher_id, student_id, title, description, start_time, end_time,
					location, meeting_url, status, cancellation_reason, reminder_sent, notes, created_at, cancelled_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	appointment.Status = "scheduled"
	_, creationErr := conn.Exec(c.Context(), createAppointmentString,
		fmt.Sprintf("%d", appointment.Id), appointment.CourseId, appointment.OfficeHoursId, appointment.TeacherId, appointment.StudentId, appointment.Title,
		appointment.Description, appointment.StartTime, appointment.EndTime, appointment.Location, appointment.MeetingUrl, appointment.Status,
		appointment.CancellationReason, appointment.ReminderSent, appointment.Notes, appointment.CreatedAt, appointment.CancelledAt)
	return creationErr
}

func BuildAppointmentUpdateString(patchField string, userId int, userRole string) string {
	return fmt.Sprintf(`UPDATE appointments SET %v=$1 WHERE id=$2 and %v_id=%v`, patchField, userRole, userId)
}

func (*AppointmentFunc) PatchAppointment(c *fiber.Ctx, conn *pgx.Conn, appointmentPatcher *structs.AppointmentPatcher, userClaims *CustomClaimStruct) error {
	exists, err := doesAppointmentExist(c, conn, appointmentPatcher)

	if !exists {
		return errors.New("no appointment to edit")
	}

	if err != nil {
		return err
	}

	updateString := BuildAppointmentUpdateString(appointmentPatcher.PatchField, userClaims.Details.Id, userClaims.Details.Role)
	_, patchErr := conn.Exec(c.Context(), updateString, appointmentPatcher.NewContent, appointmentPatcher.AppointmentId)
	return patchErr
}

func (*AppointmentFunc) DeleteAppointment(c *fiber.Ctx, conn *pgx.Conn, appointmentDeleter *structs.AppointmentDeleter, userClaims *CustomClaimStruct) error {
	var deletionString string
	switch userClaims.Details.Role {
	case "student":
		deletionString = `DELETE FROM appointments WHERE id=$1 and student_id=$2`
	case "teacher":
		deletionString = `DELETE FROM appointments WHERE id=$1 and teacher_id=$2`
	default:
		return errors.New("invalid permissions")
	}

	_, deletionErr := conn.Exec(c.Context(), deletionString, appointmentDeleter.AppointmentId, userClaims.Details.Id)
	return deletionErr
}

func (*AppointmentFunc) BuildAppointmentError(reason string) fiber.Map {
	return fiber.Map{
		"status":       "failed",
		"message":      reason,
		"appointments": nil,
	}
}

func hasDuplications(c *fiber.Ctx, conn *pgx.Conn, appointment *structs.Appointment) (bool, error) {
	studentQueryString := `SELECT EXISTS 
					(SELECT 1 FROM appointments WHERE (id=$1) OR (student_id = $2 AND NOT (end_time < $3 OR start_time > $4)))`
	teacherQueryString := `SELECT EXISTS 
					(SELECT 1 FROM appointments WHERE (id=$1) OR (student_id = $2 AND NOT (end_time < $3 OR start_time > $4)))`

	var studentConflictExists bool
	studentScanErr := conn.QueryRow(
		c.Context(), studentQueryString, appointment.Id, appointment.StudentId, appointment.StartTime.UTC(), appointment.EndTime.UTC(),
	).Scan(&studentConflictExists)

	var teacherConflictExists bool
	teacherScanErr := conn.QueryRow(
		c.Context(), teacherQueryString, appointment.Id, appointment.StudentId, appointment.StartTime.UTC(), appointment.EndTime.UTC(),
	).Scan(&teacherConflictExists)

	if studentConflictExists || teacherConflictExists {
		return true, errors.New("new appointment fails uniqueness test")
	}

	if studentScanErr != nil && teacherScanErr != nil {
		return false, errors.New(studentScanErr.Error() + " and " + teacherScanErr.Error())
	}

	if studentScanErr != nil {
		return false, studentScanErr
	}

	if teacherScanErr != nil {
		return false, teacherScanErr
	}

	return false, nil
}

func hasOfficeHoursConflict(c *fiber.Ctx, conn *pgx.Conn, appointment *structs.Appointment) (bool, error) {
	queryString := `SELECT day_of_week, start_time, end_time FROM office_hours WHERE (teacher_id = $1) AND (current_students < max_students)`
	var dayOfWeek int
	var startTime time.Time
	var endTime time.Time

	scanErr := conn.QueryRow(c.Context(), queryString, appointment.TeacherId).Scan(&dayOfWeek, &startTime, &endTime)

	if scanErr != nil {
		if errors.Is(scanErr, pgx.ErrNoRows) {
			return false, nil
		}
		return false, scanErr
	}

	if int(appointment.StartTime.Weekday()) != dayOfWeek {
		return false, errors.New("day does not match office hour day")
	}

	if !(appointment.EndTime.UTC().Before(startTime) || endTime.Before(appointment.StartTime)) {
		return true, nil
	}
	return false, errors.New("times do not match office hour time")
}

func isValidAppointment(c *fiber.Ctx, conn *pgx.Conn, appointment *structs.Appointment) (bool, error) {
	duplicateExists, dupErr := hasDuplications(c, conn, appointment)

	if duplicateExists || dupErr != nil {
		return false, dupErr
	}

	conflictExists, conflictErr := hasOfficeHoursConflict(c, conn, appointment)
	if conflictExists || conflictErr != nil {
		return false, conflictErr
	}

	return true, nil
}
