package entity

import (
	"database/sql"
	"time"
)

type Schedule struct {
	ID           int64
	InstructorID int64
	GroupID      int64
	MessageID    sql.NullInt64
	Date         time.Time
	Timeslots    []TimeSlot
}
