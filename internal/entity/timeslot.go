package entity

import "database/sql"

type TimeSlot struct {
	ID         int64
	ScheduleID int64
	Time       string
	UserID     sql.NullInt64
	UserName   sql.NullString
}
