package api

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api/consts"
	"github.com/massivemadness/schedule-service/internal/api/handler"
)

type Router struct {
	bot             *tgbotapi.BotAPI
	startHandler    *handler.StartHandler
	linkerHandler   *handler.LinkerHandler
	menuHandler     *handler.MenuHandler
	formHandler     *handler.FormHandler
	scheduleHandler *handler.ScheduleHandler
	deleteHandler   *handler.DeleteHandler
}

func NewRouter(
	bot *tgbotapi.BotAPI,
	startHandler *handler.StartHandler,
	linkerHandler *handler.LinkerHandler,
	menuHandler *handler.MenuHandler,
	formHandler *handler.FormHandler,
	scheduleHandler *handler.ScheduleHandler,
	deleteHandler *handler.DeleteHandler,
) *Router {
	return &Router{
		bot:             bot,
		startHandler:    startHandler,
		linkerHandler:   linkerHandler,
		menuHandler:     menuHandler,
		formHandler:     formHandler,
		scheduleHandler: scheduleHandler,
		deleteHandler:   deleteHandler,
	}
}

func (r *Router) HandleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		r.handleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		r.handleCallback(update.CallbackQuery)
	}
}

func (r *Router) handleMessage(msg *tgbotapi.Message) {
	switch msg.Command() {
	case consts.StartCommand:
		r.startHandler.Handle(msg)
	case consts.LinkCommand:
		r.linkerHandler.Handle(msg)
	case consts.MenuCommand:
		r.menuHandler.Handle(msg)
	}
}

func (r *Router) handleCallback(cb *tgbotapi.CallbackQuery) {
	data := cb.Data

	switch {
	case data == consts.MainMenu:
		r.menuHandler.HandleCallback(cb)
	case data == consts.Create:
		r.formHandler.HandleCreateSchedule(cb)
	case strings.HasPrefix(data, consts.SelectDate+":"):
		r.formHandler.HandleSelectDate(cb)
	case strings.HasPrefix(data, consts.SelectTime+":"):
		r.formHandler.HandleSelectTime(cb)
	case data == consts.Publish:
		r.scheduleHandler.HandlePublish(cb)
	case strings.HasPrefix(data, consts.Book+":"):
		r.scheduleHandler.HandleBooking(cb)
	case data == consts.Delete:
		r.deleteHandler.HandleDeleteSchedule(cb)
	case strings.HasPrefix(data, consts.DeleteId+":"):
		r.deleteHandler.HandleDelete(cb)
	}

	r.bot.Request(tgbotapi.NewCallback(cb.ID, ""))
}
