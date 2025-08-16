package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api/consts"
)

func NewMainMenuMessage(chatID int64) tgbotapi.MessageConfig {
	text := "📅 Главное меню управления расписанием"
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = buildMainMenuKeyboard()
	return msg
}

func EditMainMenuMessage(chatID int64, messageID int) tgbotapi.EditMessageTextConfig {
	text := "📅 Главное меню управления расписанием"
	edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, text, buildMainMenuKeyboard())
	return edit
}

func buildMainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать расписание", consts.Create),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить расписание", consts.Edit),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить расписание", consts.Delete),
		),
	)
}
