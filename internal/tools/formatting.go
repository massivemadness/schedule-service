package tools

import (
	"fmt"
	"time"
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
	return fmt.Sprintf("%s, %s", weekdays[t.Weekday()], t.Format("02.01.06"))
}
