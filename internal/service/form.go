package service

import (
	"database/sql"
	"slices"
	"strings"
	"time"

	"github.com/massivemadness/schedule-service/internal/entity"
	"github.com/massivemadness/schedule-service/internal/repository"
	"github.com/massivemadness/schedule-service/internal/tools"
)

type FormService struct {
	formRepo     repository.FormRepository
	scheduleRepo repository.ScheduleRepository
}

func NewFormService(
	formRepo repository.FormRepository,
	scheduleRepo repository.ScheduleRepository,
) *FormService {
	return &FormService{
		formRepo:     formRepo,
		scheduleRepo: scheduleRepo,
	}
}

func (s *FormService) CreateForm(instructorID int64) error {
	return s.formRepo.CreateForm(instructorID)
}

func (s *FormService) DeleteForm(instructorID int64) error {
	return s.formRepo.DeleteForm(instructorID)
}

func (s *FormService) SelectDate(instructorID int64, date string) error {
	form, err := s.formRepo.LoadForm(instructorID)
	if err != nil {
		return err
	}
	form.Date = sql.NullString{String: date, Valid: true}
	err = s.formRepo.UpdateForm(form)
	if err != nil {
		return err
	}
	return nil
}

func (s *FormService) SelectTime(instructorID int64, timeslot string) (string, error) {
	form, err := s.formRepo.LoadForm(instructorID)
	if err != nil {
		return "", err
	}

	var timeslots []string
	if form.Timeslots.Valid {
		timeslots = strings.Split(form.Timeslots.String, ",")
	}

	if slices.Contains(timeslots, timeslot) {
		timeslots = tools.RemoveValue(timeslots, timeslot)
		form.Timeslots = sql.NullString{
			String: strings.Join(timeslots, ","),
			Valid:  true,
		}
	} else {
		timeslots = append(timeslots, timeslot)
		form.Timeslots = sql.NullString{
			String: strings.Join(timeslots, ","),
			Valid:  true,
		}
	}

	err = s.formRepo.UpdateForm(form)
	if err != nil {
		return "", err
	}
	return form.Date.String, nil
}

func (s *FormService) GetAvailableDates(instructorID int64) ([]entity.DateOption, error) {
	recent, err := s.scheduleRepo.GetRecent(instructorID)
	if err != nil {
		return nil, err
	}

	existing := make(map[string]bool)
	for _, rt := range recent {
		dateID := rt.Date.Format(time.DateOnly)
		existing[dateID] = true
	}

	var availableDates []entity.DateOption
	for i := range 7 {
		date := time.Now().AddDate(0, 0, i)
		dateID := date.Format(time.DateOnly)

		// Пропускаем дату, если на этот день уже есть расписание
		if existing[dateID] {
			continue
		}

		availableDates = append(availableDates, entity.DateOption{
			ID:   dateID,
			Date: date,
		})
	}
	return availableDates, nil
}

func (s *FormService) GetAvailableTimeslots(instructorID int64) ([]entity.TimeOption, error) {
	form, err := s.formRepo.LoadForm(instructorID)
	if err != nil {
		return nil, err
	}

	selected := strings.Split(form.Timeslots.String, ",")
	selectedMap := make(map[string]bool)
	for _, time := range selected {
		selectedMap[time] = true
	}

	start := time.Date(0, 1, 1, 7, 0, 0, 0, time.UTC)
	end := time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)

	var options []entity.TimeOption
	for t := start; t.Before(end) || t.Equal(end); t = t.Add(30 * time.Minute) {
		timeID := t.Format("15:04")
		option := entity.TimeOption{
			ID:       timeID,
			Time:     t,
			Selected: selectedMap[timeID],
		}
		options = append(options, option)
	}

	return options, nil
}
