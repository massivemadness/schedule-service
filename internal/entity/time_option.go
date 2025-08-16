package entity

import "time"

type TimeOption struct {
	ID       string
	Time     time.Time
	Selected bool
}
