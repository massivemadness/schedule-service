package repository

import (
	"database/sql"
	"strings"

	"github.com/massivemadness/schedule-service/internal/database"
	"github.com/massivemadness/schedule-service/internal/database/model"
	"github.com/massivemadness/schedule-service/internal/entity"
)

type FormRepository interface {
	CreateForm(instructorID int64) error
	UpdateForm(form *entity.Form) error
	LoadForm(instructorID int64) (*entity.Form, error)
	DeleteForm(instructorID int64) error
}

type formRepositoryImpl struct {
	db *database.Database
}

func NewFormRepository(db *database.Database) FormRepository {
	return &formRepositoryImpl{db: db}
}

func (r *formRepositoryImpl) CreateForm(instructorID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO tbl_schedule_form (instructor_id) VALUES (?)`,
		instructorID,
	)
	return err
}

func (r *formRepositoryImpl) UpdateForm(form *entity.Form) error {
	var selectedDate sql.NullString
	if form.Date == "" {
		selectedDate = sql.NullString{Valid: false}
	} else {
		selectedDate = sql.NullString{
			String: form.Date,
			Valid:  true,
		}
	}
	var selectedTime sql.NullString
	if len(form.Timeslots) == 0 {
		selectedTime = sql.NullString{Valid: false}
	} else {
		selectedTime = sql.NullString{
			String: strings.Join(form.Timeslots, ","),
			Valid:  true,
		}
	}

	data := model.Form{
		ID:           form.ID,
		InstructorID: form.InstructorID,
		Date:         selectedDate,
		Timeslots:    selectedTime,
	}
	_, err := r.db.Exec(`
		UPDATE tbl_schedule_form SET date = ?, timeslots = ? WHERE id = ?`,
		data.Date,
		data.Timeslots,
		data.ID,
	)
	return err
}

func (r *formRepositoryImpl) LoadForm(instructorID int64) (*entity.Form, error) {
	row := r.db.QueryRow(
		`SELECT id, instructor_id, date, timeslots FROM tbl_schedule_form WHERE instructor_id = ? LIMIT 1`,
		instructorID,
	)
	data := &model.Form{}
	err := row.Scan(&data.ID, &data.InstructorID, &data.Date, &data.Timeslots)
	if err != nil {
		return nil, err
	}

	var selectedDate string
	if data.Date.Valid {
		selectedDate = data.Date.String
	} else {
		selectedDate = ""
	}

	var selectedTime []string
	if data.Timeslots.Valid {
		selectedTime = strings.Split(data.Timeslots.String, ",")
	}

	domain := &entity.Form{
		ID:           data.ID,
		InstructorID: data.InstructorID,
		Date:         selectedDate,
		Timeslots:    selectedTime,
	}
	return domain, nil
}

func (r *formRepositoryImpl) DeleteForm(instructorID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM tbl_schedule_form WHERE instructor_id = ?`,
		instructorID,
	)
	return err
}
