package handler

import (
	"errors"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api/menu"
	"github.com/massivemadness/schedule-service/internal/entity"
	"github.com/massivemadness/schedule-service/internal/service"
	"go.uber.org/zap"
)

type ScheduleHandler struct {
	bot             *tgbotapi.BotAPI
	logger          *zap.Logger
	scheduleService *service.ScheduleService
}

func NewScheduleHandler(
	bot *tgbotapi.BotAPI,
	logger *zap.Logger,
	scheduleService *service.ScheduleService,
) *ScheduleHandler {
	return &ScheduleHandler{
		bot:             bot,
		logger:          logger,
		scheduleService: scheduleService,
	}
}

func (h *ScheduleHandler) HandlePublish(cb *tgbotapi.CallbackQuery) {
	user := cb.From

	// Сохраняем расписание
	schedule, err := h.scheduleService.CreateSchedule(user.ID)
	if err != nil {
		h.logger.Error("Ошибка при сохранении расписания", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(user.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Публикуем сообщение в группу
	msg := menu.NewScheduleMenuMessage(schedule)
	res, err := h.bot.Send(msg)
	if err != nil {
		h.logger.Error("Ошибка при отправке сообщения", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(user.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Сохраняем ID сообщения в расписание, удаляем форму
	err = h.scheduleService.ConfirmPublished(schedule.InstructorID, schedule.ID, int64(res.MessageID))
	if err != nil {
		h.logger.Error("Ошибка при обновлении данных о сообщении", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(user.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Заменяем форму на главное меню
	h.bot.Send(menu.EditMainMenuMessage(cb.Message.Chat.ID, cb.Message.MessageID))

	// Уведомляем об отправке
	h.bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "Расписание отправлено в канал!"))
}

func (h *ScheduleHandler) HandleBooking(cb *tgbotapi.CallbackQuery) {
	parts := strings.Split(cb.Data, ":")
	if len(parts) != 3 {
		h.logger.Error("Неверный формат вызова", zap.String("data", cb.Data))
		return
	}
	scheduleID, _ := strconv.ParseInt(parts[1], 10, 64)
	timeslotID, _ := strconv.ParseInt(parts[2], 10, 64)

	userID := cb.From.ID
	userName := cb.From.FirstName
	if cb.From.LastName != "" {
		userName += " " + cb.From.LastName
	}

	// Бронируем время
	schedule, err := h.scheduleService.BookTime(scheduleID, timeslotID, userID, userName)
	if err != nil {
		if errors.Is(err, entity.ErrNotAllowed) {
			h.logger.Info("Пользователь уже записан на этот день", zap.Int64("user_id", userID), zap.Int64("schedule_id", scheduleID))
		} else if errors.Is(err, entity.ErrAlreadyBooked) {
			h.logger.Info("Время уже занято", zap.Int64("user_id", userID), zap.Int64("schedule_id", scheduleID), zap.Int64("timeslot_id", timeslotID))
		} else {
			h.logger.Error("Ошибка при бронировании времени", zap.Error(err))
		}
		return
	}

	// Обновляем расписание
	msg := menu.EditScheduleMenuMessage(schedule)
	_, err = h.bot.Send(msg)
	if err != nil {
		h.logger.Error("Ошибка при обновлении сообщения", zap.Error(err))
	}
}
