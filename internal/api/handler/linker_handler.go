package handler

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api/menu"
	"github.com/massivemadness/schedule-service/internal/entity"
	"github.com/massivemadness/schedule-service/internal/service"
	"go.uber.org/zap"
)

type LinkerHandler struct {
	bot               *tgbotapi.BotAPI
	logger            *zap.Logger
	instructorService *service.InstructorService
}

func NewLinkerHandler(
	bot *tgbotapi.BotAPI,
	logger *zap.Logger,
	instructorService *service.InstructorService,
) *LinkerHandler {
	return &LinkerHandler{
		bot:               bot,
		logger:            logger,
		instructorService: instructorService,
	}
}

func (h *LinkerHandler) Handle(msg *tgbotapi.Message) {
	chat := msg.Chat
	user := msg.From

	// Проверяем где была введена команда
	if chat.Type != "supergroup" && chat.Type != "group" {
		h.bot.Send(tgbotapi.NewMessage(chat.ID, "Команда /link_group должна вызываться в группе"))
		return
	}

	// Проверяем, что пользователь админ или создатель группы
	chatMember, err := h.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chat.ID,
			UserID: user.ID,
		},
	})
	if err != nil {
		h.logger.Error("Ошибка при получении статуса пользователя в группе", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(chat.ID, "Не удалось проверить ваши права в группе"))
		return
	}
	if chatMember.Status != "administrator" && chatMember.Status != "creator" {
		h.bot.Send(tgbotapi.NewMessage(chat.ID, "Вы должны быть администратором группы, чтобы связать её с вашим аккаунтом"))
		return
	}

	// Связываем инструктора с группой
	err = h.instructorService.LinkGroup(user.ID, chat.ID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			h.logger.Info("Инструктор не зарегистрирован", zap.Int64("user_id", user.ID), zap.Int64("chat_id", chat.ID))
			h.bot.Send(tgbotapi.NewMessage(chat.ID, "Вы не зарегистрированы как инструктор. Сначала запустите бота /start"))
		} else if errors.Is(err, entity.ErrAlreadyLinked) {
			h.logger.Info("Группа уже привязана", zap.Int64("user_id", user.ID), zap.Int64("chat_id", chat.ID))
			h.bot.Send(tgbotapi.NewMessage(chat.ID, "Группа уже привязана, перейдите в меню бота для управления"))
		} else if errors.Is(err, entity.ErrOtherUserLinked) {
			h.logger.Info("Группа уже привязана", zap.Int64("user_id", user.ID), zap.Int64("chat_id", chat.ID))
			h.bot.Send(tgbotapi.NewMessage(chat.ID, "Эта группа уже связана с другим инструктором"))
		} else {
			h.logger.Error("Ошибка при привязке группы", zap.Error(err))
			h.bot.Send(tgbotapi.NewMessage(chat.ID, "Произошла ошибка, попробуйте позже"))
		}
		return
	}

	h.logger.Info("Группа успешно привязана", zap.Int64("user_id", user.ID), zap.Int64("chat_id", chat.ID))
	h.bot.Send(tgbotapi.NewMessage(chat.ID, "Группа привязана! Теперь вы можете публиковать сюда расписания"))

	// Отправляем меню управления
	h.bot.Send(menu.NewMainMenuMessage(user.ID))
}
