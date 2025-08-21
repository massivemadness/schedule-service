package entity

type Form struct {
	ID           int64
	InstructorID int64
	Date         string
	Timeslots    []string
}
