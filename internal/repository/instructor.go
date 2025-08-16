package repository

import (
	"database/sql"

	"github.com/massivemadness/schedule-service/internal/database"
	"github.com/massivemadness/schedule-service/internal/entity"
)

type InstructorRepository interface {
	CreateIfNotExists(instructorID int64, instructorName string) error
	GetByUserID(instructorID int64) (*entity.Instructor, error)
}

type instructorRepositoryImpl struct {
	db *database.Database
}

func NewInstructorRepository(db *database.Database) InstructorRepository {
	return &instructorRepositoryImpl{db: db}
}

func (r *instructorRepositoryImpl) CreateIfNotExists(instructorID int64, instructorName string) error {
	_, err := r.db.Exec(`INSERT OR IGNORE INTO tbl_instructors (id, name) VALUES (?, ?)`, instructorID, instructorName)
	return err
}

func (r *instructorRepositoryImpl) GetByUserID(instructorID int64) (*entity.Instructor, error) {
	row := r.db.QueryRow(`SELECT id, name FROM tbl_instructors WHERE id = ?`, instructorID)

	var instr entity.Instructor
	err := row.Scan(&instr.ID, &instr.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &instr, nil
}
