package menu

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api/consts"
	"github.com/massivemadness/schedule-service/internal/entity"
	"github.com/massivemadness/schedule-service/internal/tools"
)

func EditDeleteMenuMessage(
	chatID int64,
	messageID int,
	schedules []entity.Schedule,
) tgbotapi.EditMessageTextConfig {
	text := "🗑 Удаление расписания\n\n(нажмите для удаления)"
	edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, text, buildDeleteMenuKeyboard(schedules))
	return edit
}

func buildDeleteMenuKeyboard(schedules []entity.Schedule) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, s := range schedules {
		text := tools.FormatRuDate(s.Date)
		data := consts.DeleteId + fmt.Sprintf(":%d", s.ID)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, data),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", consts.MainMenu),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
