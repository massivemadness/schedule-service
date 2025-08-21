package tools

import (
	"fmt"
	"time"
)

const (
	HumanTime string = "15:04"
	HumanDate string = "02.01.06"
)

var weekdays = []string{
	"Воскресенье",
	"Понедельник",
	"Вторник",
	"Среда",
	"Четверг",
	"Пятница",
	"Суббота",
}

func FormatRuDate(t time.Time) string {
	return fmt.Sprintf("%s, %s", weekdays[t.Weekday()], t.Format(HumanDate))
}
