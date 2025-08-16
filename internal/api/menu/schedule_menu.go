package menu

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api/consts"
	"github.com/massivemadness/schedule-service/internal/entity"
	"github.com/massivemadness/schedule-service/internal/tools"
)

func NewScheduleMenuMessage(schedule *entity.Schedule) tgbotapi.MessageConfig {
	text := fmt.Sprintf("📅 %s:\n", tools.FormatRuDate(schedule.Date))
	for _, timeslot := range schedule.Timeslots {
		var student string
		if timeslot.UserName.Valid {
			student = timeslot.UserName.String
		} else {
			student = "________________"
		}
		text += fmt.Sprintf("%s %s\n", timeslot.Time, student)
	}
	text += "\nЕсть свободные места, выберите время чтобы записаться ⬇️"
	msg := tgbotapi.NewMessage(schedule.GroupID, text)
	msg.ReplyMarkup = buildScheduleMenuKeyboard(schedule.Timeslots)
	return msg
}

func EditScheduleMenuMessage(schedule *entity.Schedule) tgbotapi.EditMessageTextConfig {
	text := fmt.Sprintf("📅 %s:\n", tools.FormatRuDate(schedule.Date))
	freeSlots := false
	for _, timeslot := range schedule.Timeslots {
		var student string
		if timeslot.UserName.Valid {
			student = timeslot.UserName.String
		} else {
			freeSlots = true
			student = "________________"
		}
		text += fmt.Sprintf("%s %s\n", timeslot.Time, student)
	}
	if freeSlots {
		text += "\nЕсть свободные места, выберите время чтобы записаться ⬇️"
	}
	edit := tgbotapi.NewEditMessageTextAndMarkup(
		schedule.GroupID,
		int(schedule.MessageID.Int64),
		text,
		buildScheduleMenuKeyboard(schedule.Timeslots),
	)
	return edit
}

func buildScheduleMenuKeyboard(timeslots []entity.TimeSlot) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	for _, slot := range timeslots {
		if slot.UserID.Valid {
			// Слот уже занят
			continue
		}
		text := slot.Time
		callback := consts.Book + fmt.Sprintf(":%d:%d", slot.ScheduleID, slot.ID)
		btn := tgbotapi.NewInlineKeyboardButtonData(text, callback)

		row = append(row, btn)

		if len(row) == 3 { // 3 кнопки в ряд
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}

	// Если свободных слотов нет
	if len(rows) == 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет свободных мест", "noop"),
		))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
