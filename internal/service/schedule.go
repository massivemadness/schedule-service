package service

import (
	"database/sql"
	"sort"
	"strings"
	"time"

	"github.com/massivemadness/schedule-service/internal/entity"
	"github.com/massivemadness/schedule-service/internal/repository"
)

type ScheduleService struct {
	instructorRepo repository.InstructorRepository
	groupRepo      repository.InstructorGroupRepository
	formRepo       repository.FormRepository
	scheduleRepo   repository.ScheduleRepository
	timeslotRepo   repository.TimeslotRepository
}

func NewScheduleService(
	instructorRepo repository.InstructorRepository,
	groupRepo repository.InstructorGroupRepository,
	formRepo repository.FormRepository,
	scheduleRepo repository.ScheduleRepository,
	timeslotRepo repository.TimeslotRepository,
) *ScheduleService {
	return &ScheduleService{
		instructorRepo: instructorRepo,
		groupRepo:      groupRepo,
		formRepo:       formRepo,
		scheduleRepo:   scheduleRepo,
		timeslotRepo:   timeslotRepo,
	}
}

func (s *ScheduleService) CreateSchedule(instructorID int64) (*entity.Schedule, error) {
	// Загружаем заполненную форму
	form, err := s.formRepo.LoadForm(instructorID)
	if err != nil {
		return nil, err
	}

	// Загружаем связанную группу
	group, err := s.groupRepo.GetByUserID(form.InstructorID)
	if err != nil {
		return nil, err
	}

	selectedDate, err := time.Parse(time.DateOnly, form.Date.String)
	if err != nil {
		return nil, err
	}

	// Сортируем слоты по времени
	rawStringSlots := strings.Split(form.Timeslots.String, ",")
	sort.Slice(rawStringSlots, func(i, j int) bool {
		ti, _ := time.Parse("15:04", rawStringSlots[i])
		tj, _ := time.Parse("15:04", rawStringSlots[j])
		return ti.Before(tj)
	})

	// Сохраняем расписание в БД
	schedule := entity.Schedule{
		InstructorID: form.InstructorID,
		GroupID:      group.GroupID,
		Date:         selectedDate,
	}
	scheduleId, err := s.scheduleRepo.Save(&schedule)
	if err != nil {
		return nil, err
	}
	schedule.ID = scheduleId

	// Переводим строку в модель слота
	var timeslots []entity.TimeSlot
	for _, timeStr := range sort.StringSlice(rawStringSlots) {
		timeslot := entity.TimeSlot{
			ScheduleID: scheduleId,
			Time:       timeStr,
			UserID:     sql.NullInt64{Valid: false},  // свободно
			UserName:   sql.NullString{Valid: false}, // свободно
		}
		timeslots = append(timeslots, timeslot)
	}

	// Сохраняем таймслоты в БД
	err = s.timeslotRepo.Save(timeslots)
	if err != nil {
		return nil, err
	}

	// Загружаем слоты, уже с заполненным ID
	timeslots, err = s.timeslotRepo.GetByScheduleID(scheduleId)
	if err != nil {
		return nil, err
	}
	schedule.Timeslots = timeslots

	return &schedule, nil
}

func (s *ScheduleService) ConfirmPublished(instructorID int64, scheduleID int64, messageID int64) error {
	err := s.scheduleRepo.SetMessageId(scheduleID, messageID)

	// Очищаем старую форму (допускаются ошибки)
	_ = s.formRepo.DeleteForm(instructorID)

	return err
}

func (s *ScheduleService) LoadRecent(instructorID int64) ([]entity.Schedule, error) {
	return s.scheduleRepo.GetRecent(instructorID)
}

func (s *ScheduleService) DeleteSchedule(scheduleID int64) (*entity.Schedule, error) {
	schedule, err := s.scheduleRepo.LoadById(scheduleID)
	if err != nil {
		return nil, err
	}

	group, err := s.groupRepo.GetByUserID(schedule.InstructorID)
	if err != nil {
		return nil, err
	}
	schedule.GroupID = group.GroupID

	err = s.scheduleRepo.DeleteById(schedule.ID)
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (s *ScheduleService) BookTime(scheduleID int64, timeslotID int64, userID int64, userName string) (*entity.Schedule, error) {
	// Проверяем, свободна ли запись на этот слот
	timeslots, err := s.timeslotRepo.GetByScheduleID(scheduleID)
	if err != nil {
		return nil, err
	}

	// Проверяем, записан ли пользователь на другой слот
	for _, timeslot := range timeslots {
		if timeslot.UserID.Valid && timeslot.UserID.Int64 == userID {
			return nil, entity.ErrNotAllowed
		}
	}

	// Проверяем занят ли уже этот слот
	var selected *entity.TimeSlot
	for _, timeslot := range timeslots {
		if timeslot.ID == timeslotID {
			if !timeslot.UserID.Valid {
				selected = &timeslot
				break
			}
		}
	}
	if selected == nil {
		return nil, entity.ErrAlreadyBooked
	}

	// Обновляем запись в БД
	err = s.timeslotRepo.BookSlot(selected.ID, userID, userName)
	if err != nil {
		return nil, err
	}

	// Формируем новое расписание
	schedule, err := s.scheduleRepo.LoadById(scheduleID)
	if err != nil {
		return nil, err
	}

	group, err := s.groupRepo.GetByUserID(schedule.InstructorID)
	if err != nil {
		return nil, err
	}
	schedule.GroupID = group.GroupID

	timeslots, err = s.timeslotRepo.GetByScheduleID(scheduleID)
	if err != nil {
		return nil, err
	}
	schedule.Timeslots = timeslots

	return schedule, nil
}
