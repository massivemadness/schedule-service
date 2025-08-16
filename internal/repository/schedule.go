package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/massivemadness/schedule-service/internal/database"
	"github.com/massivemadness/schedule-service/internal/entity"
)

type ScheduleRepository interface {
	GetRecent(instructorID int64) ([]entity.Schedule, error)
	Save(schedule *entity.Schedule) (int64, error)
	SetMessageId(scheduleID int64, messageId int64) error
	LoadById(scheduleID int64) (*entity.Schedule, error)
	DeleteById(scheduleID int64) error
}

type scheduleRepositoryImpl struct {
	db *database.Database
}

func NewScheduleRepository(db *database.Database) ScheduleRepository {
	return &scheduleRepositoryImpl{db: db}
}

func (r *scheduleRepositoryImpl) GetRecent(instructorID int64) ([]entity.Schedule, error) {
	today := time.Now().Format(time.DateOnly)
	rows, err := r.db.Query(`
        SELECT id, instructor_id, message_id, date
        FROM tbl_schedules
        WHERE date >= ? AND instructor_id = ?
        ORDER BY date ASC
    `, today, instructorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []entity.Schedule
	for rows.Next() {
		var dateStr string
		var s entity.Schedule
		if err := rows.Scan(&s.ID, &s.InstructorID, &s.MessageID, &dateStr); err != nil {
			return nil, err
		}
		s.Date, _ = time.Parse("2006-01-02", dateStr)
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *scheduleRepositoryImpl) Save(schedule *entity.Schedule) (int64, error) {
	res, err := r.db.Exec(`
	INSERT INTO tbl_schedules (instructor_id, date)
	VALUES (?, ?)`, schedule.InstructorID, schedule.Date.Format("2006-01-02"))
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func (r *scheduleRepositoryImpl) SetMessageId(scheduleID int64, messageId int64) error {
	_, err := r.db.Exec(`
	UPDATE tbl_schedules SET message_id = ? WHERE id = ?`, messageId, scheduleID)
	return err
}

func (r *scheduleRepositoryImpl) LoadById(scheduleID int64) (*entity.Schedule, error) {
	row := r.db.QueryRow(`
        SELECT id, instructor_id, message_id, date
        FROM tbl_schedules
        WHERE id = ? LIMIT 1
    `, scheduleID)

	var dateStr string
	var s entity.Schedule
	err := row.Scan(&s.ID, &s.InstructorID, &s.MessageID, &dateStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	s.Date, _ = time.Parse("2006-01-02", dateStr)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("schedule with id %d not found", scheduleID)
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *scheduleRepositoryImpl) DeleteById(scheduleID int64) error {
	_, err := r.db.Exec(`
	DELETE FROM tbl_schedules WHERE id = ?`, scheduleID)
	return err
}
