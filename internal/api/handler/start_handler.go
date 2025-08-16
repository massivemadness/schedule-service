package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/service"
	"go.uber.org/zap"
)

type StartHandler struct {
	bot               *tgbotapi.BotAPI
	logger            *zap.Logger
	instructorService *service.InstructorService
}

func NewStartHandler(
	bot *tgbotapi.BotAPI,
	logger *zap.Logger,
	instructorService *service.InstructorService,
) *StartHandler {
	return &StartHandler{
		bot:               bot,
		logger:            logger,
		instructorService: instructorService,
	}
}

func (h *StartHandler) Handle(msg *tgbotapi.Message) {
	chat := msg.Chat

	// Проверяем где была введена команда
	if chat.Type != "private" {
		h.bot.Send(tgbotapi.NewMessage(chat.ID, "Команду /start нужно вводить в личных сообщениях боту"))
		return
	}

	instructorID := chat.ID
	instructorName := chat.FirstName
	if chat.LastName != "" {
		instructorName += " " + chat.LastName
	}

	// Регистрируем инструктора в БД
	err := h.instructorService.Register(instructorID, instructorName)
	if err != nil {
		h.logger.Error(
			"Ошибка при регистрации инструктора",
			zap.Int64("user_id", instructorID),
			zap.String("user_name", instructorName),
			zap.Error(err),
		)
		h.bot.Send(tgbotapi.NewMessage(chat.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	h.logger.Info(
		"Инструктор зарегистрирован",
		zap.Int64("user_id", instructorID),
		zap.String("user_name", instructorName),
	)

	h.bot.Send(tgbotapi.NewMessage(chat.ID, "Вы зарегистрированы как инструктор.\nДобавьте бота в группу и выполните команду /link_group"))
}
