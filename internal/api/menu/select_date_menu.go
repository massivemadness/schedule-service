package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api/consts"
	"github.com/massivemadness/schedule-service/internal/entity"
	"github.com/massivemadness/schedule-service/internal/tools"
)

func EditSelectDateMenuMessage(
	chatID int64,
	messageID int,
	availableDates []entity.DateOption,
) tgbotapi.EditMessageTextConfig {
	text := "🆕 Создание расписания\n\nВыберите дату:"
	edit := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		text,
		buildSelectDateMenuKeyboard(availableDates),
	)
	return edit
}

func buildSelectDateMenuKeyboard(availableDates []entity.DateOption) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, d := range availableDates {
		button := tgbotapi.NewInlineKeyboardButtonData(
			tools.FormatRuDate(d.Date),
			consts.SelectDate+":"+d.ID,
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", consts.MainMenu),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
