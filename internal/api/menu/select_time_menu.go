package menu

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api/consts"
	"github.com/massivemadness/schedule-service/internal/entity"
	"github.com/massivemadness/schedule-service/internal/tools"
)

func EditSelectTimeMenuMessage(
	chatID int64,
	messageID int,
	date time.Time,
	timeslots []entity.TimeOption,
) tgbotapi.EditMessageTextConfig {
	text := fmt.Sprintf("🆕 Вы выбрали дату: %s.\n\nТеперь выберите время:", tools.FormatRuDate(date))
	edit := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		text,
		buildSelectTimeMenuKeyboard(timeslots),
	)
	return edit
}

func buildSelectTimeMenuKeyboard(timeslots []entity.TimeOption) tgbotapi.InlineKeyboardMarkup {
	var selected []entity.TimeOption
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for _, t := range timeslots {
		text := t.ID
		if t.Selected {
			text = "✅ " + text
			selected = append(selected, t)
		}
		button := tgbotapi.NewInlineKeyboardButtonData(
			text,
			consts.SelectTime+":"+t.ID,
		)
		row = append(row, button)
		if len(row) == 4 { // 4 кнопки в ряд
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}
	if len(selected) > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Опубликовать", consts.Publish),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", consts.MainMenu),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
