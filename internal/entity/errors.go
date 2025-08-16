package entity

import "github.com/pkg/errors"

var (
	ErrNotFound        error = errors.New("not_found")
	ErrNotLinked       error = errors.New("group_not_linked")
	ErrAlreadyLinked   error = errors.New("group_already_linked")
	ErrOtherUserLinked error = errors.New("group_other_user_linked")
	ErrNotAllowed      error = errors.New("timeslot_not_allowed")
	ErrAlreadyBooked   error = errors.New("timeslot_already_booked")
)
