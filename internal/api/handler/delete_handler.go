package handler

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api/consts"
	"github.com/massivemadness/schedule-service/internal/api/menu"
	"github.com/massivemadness/schedule-service/internal/service"
	"go.uber.org/zap"
)

type DeleteHandler struct {
	bot             *tgbotapi.BotAPI
	logger          *zap.Logger
	scheduleService *service.ScheduleService
}

func NewDeleteHandler(
	bot *tgbotapi.BotAPI,
	logger *zap.Logger,
	scheduleService *service.ScheduleService,
) *DeleteHandler {
	return &DeleteHandler{
		bot:             bot,
		logger:          logger,
		scheduleService: scheduleService,
	}
}

func (h *DeleteHandler) HandleDeleteSchedule(cb *tgbotapi.CallbackQuery) {
	// Загружаем ближайшие расписания
	schedules, err := h.scheduleService.LoadRecent(cb.From.ID)
	if err != nil {
		h.logger.Error("Ошибка при получении расписаний", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(cb.From.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Отправляем список расписаний
	msg := menu.EditDeleteMenuMessage(
		cb.Message.Chat.ID,
		cb.Message.MessageID,
		schedules,
	)
	h.bot.Send(msg)
}

func (h *DeleteHandler) HandleDelete(cb *tgbotapi.CallbackQuery) {
	// Парсим айди расписания
	scheduleID, err := strconv.ParseInt(strings.TrimPrefix(cb.Data, consts.DeleteId+":"), 10, 64)
	if err != nil {
		h.logger.Error("Не удалось распарсить ID", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(cb.From.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Удаляем расписание
	schedule, err := h.scheduleService.DeleteSchedule(scheduleID)
	if err != nil {
		h.logger.Error("Ошибка при удалении расписания", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(cb.From.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Удалить из группы
	del := tgbotapi.NewDeleteMessage(schedule.GroupID, int(schedule.MessageID.Int64))
	if _, err := h.bot.Send(del); err != nil {
		h.logger.Error("Ошибка при удалении сообщения:", zap.Error(err))
	}

	// Обновить меню
	h.bot.Send(menu.EditMainMenuMessage(cb.Message.Chat.ID, cb.Message.MessageID))

	// Уведомить об удалении
	h.bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "Расписание удалено"))
}
