package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/massivemadness/schedule-service/internal/api"
	"github.com/massivemadness/schedule-service/internal/api/handler"
	"github.com/massivemadness/schedule-service/internal/config"
	"github.com/massivemadness/schedule-service/internal/database"
	"github.com/massivemadness/schedule-service/internal/logger"
	"github.com/massivemadness/schedule-service/internal/repository"
	"github.com/massivemadness/schedule-service/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()
	zapLogger := logger.NewLogger(cfg.Env)
	db, err := database.New(cfg)
	if err != nil {
		zapLogger.Fatal("Не удалось подключиться к БД", zap.Error(err))
	}
	defer db.Close()

	instructorRepo := repository.NewInstructorRepository(db)
	groupRepo := repository.NewInstructorGroupRepository(db)
	formRepo := repository.NewFormRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	timeslotRepo := repository.NewTimeslotRepository(db)

	instructorService := service.NewInstructorService(instructorRepo, groupRepo)
	formService := service.NewFormService(formRepo, scheduleRepo)
	scheduleService := service.NewScheduleService(instructorRepo, groupRepo, formRepo, scheduleRepo, timeslotRepo)

	zapLogger.Info("Подключение к API телеграма...")

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		zapLogger.Fatal("Не удалось подключиться к Telegram", zap.Error(err))
	}

	startHandler := handler.NewStartHandler(bot, zapLogger, instructorService)
	linkerHandler := handler.NewLinkerHandler(bot, zapLogger, instructorService)
	menuHandler := handler.NewMenuHandler(bot, zapLogger, instructorService)
	formHandler := handler.NewFormHandler(bot, zapLogger, formService)
	scheduleHandler := handler.NewScheduleHandler(bot, zapLogger, scheduleService)
	deleteHandler := handler.NewDeleteHandler(bot, zapLogger, scheduleService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = int(cfg.Telegram.Timeout / time.Second)
	updates := bot.GetUpdatesChan(u)
	router := api.NewRouter(
		bot,
		startHandler,
		linkerHandler,
		menuHandler,
		formHandler,
		scheduleHandler,
		deleteHandler,
	)

	zapLogger.Info("Бот запущен", zap.String("username", bot.Self.UserName))

	for update := range updates {
		router.HandleUpdate(update)
	}

	zapLogger.Info("Бот остановлен")
}
