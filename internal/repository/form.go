package repository

import (
	"database/sql"

	"github.com/massivemadness/schedule-service/internal/database"
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
	_, err := r.db.Exec(`
		UPDATE tbl_schedule_form SET date = ?, timeslots = ? WHERE id = ?`,
		form.Date,
		form.Timeslots,
		form.ID,
	)
	return err
}

func (r *formRepositoryImpl) LoadForm(instructorID int64) (*entity.Form, error) {
	row := r.db.QueryRow(
		`SELECT id, instructor_id, date, timeslots FROM tbl_schedule_form WHERE instructor_id = ? LIMIT 1`,
		instructorID,
	)
	s := &entity.Form{}
	err := row.Scan(&s.ID, &s.InstructorID, &s.Date, &s.Timeslots)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func (r *formRepositoryImpl) DeleteForm(instructorID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM tbl_schedule_form WHERE instructor_id = ?`,
		instructorID,
	)
	return err
}
