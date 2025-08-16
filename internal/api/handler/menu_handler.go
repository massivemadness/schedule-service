package handler

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api/menu"
	"github.com/massivemadness/schedule-service/internal/entity"
	"github.com/massivemadness/schedule-service/internal/service"
	"go.uber.org/zap"
)

type MenuHandler struct {
	bot               *tgbotapi.BotAPI
	logger            *zap.Logger
	instructorService *service.InstructorService
}

func NewMenuHandler(
	bot *tgbotapi.BotAPI,
	logger *zap.Logger,
	instructorService *service.InstructorService,
) *MenuHandler {
	return &MenuHandler{
		bot:               bot,
		logger:            logger,
		instructorService: instructorService,
	}
}

func (h *MenuHandler) Handle(msg *tgbotapi.Message) {
	chat := msg.Chat

	// Проверяем где была введена команда
	if chat.Type != "private" {
		h.bot.Send(tgbotapi.NewMessage(chat.ID, "Команда /menu доступна только в личных сообщениях с ботом"))
		return
	}

	// Проверяем условия: инструктор зарегистрировался и связал группу
	err := h.instructorService.CheckIsRegistered(chat.ID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			h.logger.Info("Инструктор не зарегистрирован", zap.Int64("chat_id", chat.ID))
			h.bot.Send(tgbotapi.NewMessage(chat.ID, "Вы не зарегистрированы как инструктор. Сначала запустите бота /start"))
		} else if errors.Is(err, entity.ErrNotLinked) {
			h.logger.Info("Группа не привязана", zap.Int64("chat_id", chat.ID))
			h.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Сначала привяжите группу командой /link_group"))
		} else {
			h.logger.Error("Ошибка при проверке регистрации", zap.Error(err))
			h.bot.Send(tgbotapi.NewMessage(chat.ID, "Произошла ошибка, попробуйте позже"))
		}
		return
	}

	// Отправляем меню управления
	h.bot.Send(menu.NewMainMenuMessage(chat.ID))
}

func (h *MenuHandler) HandleCallback(cb *tgbotapi.CallbackQuery) {
	// Удаляем временные записи
	// err := h.scheduleService.DeleteForm(cb.Message.Chat.ID)
	// if err != nil {
	// 	h.logger.Error("Ошибка при удалении временных записей", zap.Error(err))
	// }

	// Заменяем сообщение на меню управления
	h.bot.Send(menu.EditMainMenuMessage(cb.Message.Chat.ID, cb.Message.MessageID))
}
