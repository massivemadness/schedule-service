package model

import (
	"database/sql"
)

type Form struct {
	ID           int64
	InstructorID int64
	Date         sql.NullString
	Timeslots    sql.NullString
}
