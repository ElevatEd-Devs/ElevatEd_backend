package functions

import (
	"elevated_backend/structs"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type EventFunc struct{}

func (*EventFunc) GetEvents(c *fiber.Ctx, conn *pgx.Conn, userClaims *CustomClaimStruct) ([]structs.Event, error) {
	queryString := `SELECT * FROM events WHERE course_id IN
					(SELECT course_id FROM course_memberships WHERE user_id=$1 and role=$2 AND completed_at is NULL AND dropped_at is NULL)`
	rows, err := conn.Query(c.Context(), queryString, userClaims.Details.Id, userClaims.Details.Role)

	if err != nil {
		return nil, err
	}

	var events []structs.Event
	for rows.Next() {
		var event structs.Event
		scanErr := rows.Scan(&event.Id, &event.CourseId, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.Location, &event.EventType, &event.CreatedBy, &event.CreatedAt)
		if scanErr != nil {
			return nil, scanErr
		}
		events = append(events, event)
	}

	return events, nil
}

func hasDuplications(c *fiber.Ctx, conn *pgx.Conn, event *structs.Event) (bool, error) {
	queryString := `SELECT EXISTS 
					(SELECT 1 FROM events WHERE (id=$1) OR (course_id = $2 AND NOT (end_time < $3 OR start_time > $4)))`

	var eventConflictExists bool
	eventScanErr := conn.QueryRow(
		c.Context(), queryString, event.Id, event.CourseId, event.StartTime.UTC(), event.EndTime.UTC(),
	).Scan(&eventConflictExists)

	if eventScanErr != nil {
		return false, eventScanErr
	}

	if eventConflictExists {
		return true, errors.New("new event fails uniqueness test")
	}

	return false, nil
}

func (*EventFunc) CreateEvent(c *fiber.Ctx, conn *pgx.Conn, event *structs.Event) error {
	eventDuplicated, getErr := hasDuplications(c, conn, event)

	if getErr != nil || eventDuplicated {
		return getErr
	}

	creationString := `INSERT INTO events
					   (id, course_id, title, description, start_time, end_time, location, event_type, created_by, created_at)
					   VALUES
					   ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, creationErr := conn.Exec(c.Context(), creationString, event.Id, event.CourseId, event.Title, event.Description, event.StartTime.UTC(), event.EndTime.UTC(), event.Location, event.EventType, event.CreatedBy, event.CreatedAt.UTC())

	return creationErr
}

func (*EventFunc) PatchEvent(c *fiber.Ctx, conn *pgx.Conn, eventPatcher *structs.EventPatcher) error {
	updateString := fmt.Sprintf(`UPDATE events SET %v = %v WHERE id = %v`, eventPatcher.PatchField, eventPatcher.NewContent, eventPatcher.Id)
	_, patchErr := conn.Exec(c.Context(), updateString)
	return patchErr
}

func (*EventFunc) DeleteEvent(c *fiber.Ctx, conn *pgx.Conn, eventDeleter *structs.EventDeleter) error {
	deleteString := `DELETE FROM events WHERE id = $1`
	tag, deleteErr := conn.Exec(c.Context(), deleteString, eventDeleter.Id)

	rowsAffected := tag.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("item with id = %v cannot be deleted because it does not exist", eventDeleter.Id)
	}
	return deleteErr
}

func (*EventFunc) BuildErrorString(reason string) fiber.Map {
	return fiber.Map{
		"status":  "failed",
		"message": reason,
	}
}
