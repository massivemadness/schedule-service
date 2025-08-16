package handler

import (
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api/consts"
	"github.com/massivemadness/schedule-service/internal/api/menu"
	"github.com/massivemadness/schedule-service/internal/service"
	"go.uber.org/zap"
)

type FormHandler struct {
	bot         *tgbotapi.BotAPI
	logger      *zap.Logger
	formService *service.FormService
}

func NewFormHandler(
	bot *tgbotapi.BotAPI,
	logger *zap.Logger,
	formService *service.FormService,
) *FormHandler {
	return &FormHandler{
		bot:         bot,
		logger:      logger,
		formService: formService,
	}
}

func (h *FormHandler) HandleCreateSchedule(cb *tgbotapi.CallbackQuery) {
	err := h.formService.CreateForm(cb.From.ID)
	if err != nil {
		h.logger.Error("Ошибка при создании формы", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(cb.From.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Получаем даты без расписаний
	availableDates, err := h.formService.GetAvailableDates(cb.From.ID)
	if err != nil {
		h.logger.Error("Ошибка при получении свободных дней", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(cb.From.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Отправляем список доступных дат
	msg := menu.EditSelectDateMenuMessage(
		cb.Message.Chat.ID,
		cb.Message.MessageID,
		availableDates,
	)
	h.bot.Send(msg)
}

func (h *FormHandler) HandleSelectDate(cb *tgbotapi.CallbackQuery) {
	selectedDate := strings.TrimPrefix(cb.Data, consts.SelectDate+":")

	// Парсим выбранную дату
	formattedDate, err := time.Parse(time.DateOnly, selectedDate)
	if err != nil {
		h.logger.Error("Не удалось распарсить дату", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(cb.From.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Обновляем дату в БД
	err = h.formService.SelectDate(cb.From.ID, selectedDate)
	if err != nil {
		h.logger.Error("Ошибка при обновлении даты", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(cb.From.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Получаем таймслоты для создания расписания
	timeslots, err := h.formService.GetAvailableTimeslots(cb.From.ID)
	if err != nil {
		h.logger.Error("Ошибка при получении доступных слотов", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(cb.From.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Отправляем список таймслотов
	msg := menu.EditSelectTimeMenuMessage(
		cb.Message.Chat.ID,
		cb.Message.MessageID,
		formattedDate,
		timeslots,
	)
	h.bot.Send(msg)
}

func (h *FormHandler) HandleSelectTime(cb *tgbotapi.CallbackQuery) {
	selectedTime := strings.TrimPrefix(cb.Data, consts.SelectTime+":")

	// Обновляем время в БД
	selectedDate, err := h.formService.SelectTime(cb.From.ID, selectedTime)
	if err != nil {
		h.logger.Error("Ошибка при обновлении времени", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(cb.From.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Парсим выбранную дату
	formattedDate, err := time.Parse(time.DateOnly, selectedDate)
	if err != nil {
		h.logger.Error("Не удалось распарсить дату", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(cb.From.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Получаем обновленные таймслоты
	timeslots, err := h.formService.GetAvailableTimeslots(cb.From.ID)
	if err != nil {
		h.logger.Error("Ошибка при получении доступных слотов", zap.Error(err))
		h.bot.Send(tgbotapi.NewMessage(cb.From.ID, "Произошла ошибка, попробуйте позже"))
		return
	}

	// Обновляем список таймслотов
	msg := menu.EditSelectTimeMenuMessage(
		cb.Message.Chat.ID,
		cb.Message.MessageID,
		formattedDate,
		timeslots,
	)
	h.bot.Send(msg)
}
