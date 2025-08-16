package repository

import (
	"database/sql"

	"github.com/massivemadness/schedule-service/internal/database"
	"github.com/massivemadness/schedule-service/internal/entity"
)

type InstructorGroupRepository interface {
	LinkGroup(instructorID int64, groupID int64) error
	GetByGroupID(groupID int64) (*entity.InstructorGroup, error)
	GetByUserID(instructorID int64) (*entity.InstructorGroup, error)
}

type instructorGroupRepositoryImpl struct {
	db *database.Database
}

func NewInstructorGroupRepository(db *database.Database) InstructorGroupRepository {
	return &instructorGroupRepositoryImpl{db: db}
}

func (r *instructorGroupRepositoryImpl) LinkGroup(instructorID int64, groupID int64) error {
	_, err := r.db.Exec(`
        INSERT INTO tbl_instructor_groups (group_id, instructor_id)
        VALUES (?, ?)
    `, groupID, instructorID)
	return err
}

func (r *instructorGroupRepositoryImpl) GetByGroupID(groupID int64) (*entity.InstructorGroup, error) {
	row := r.db.QueryRow(`
        SELECT group_id, instructor_id
        FROM tbl_instructor_groups
        WHERE group_id = ?
    `, groupID)

	var group entity.InstructorGroup
	err := row.Scan(&group.GroupID, &group.InstructorID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *instructorGroupRepositoryImpl) GetByUserID(instructorID int64) (*entity.InstructorGroup, error) {
	row := r.db.QueryRow(`
        SELECT group_id, instructor_id
        FROM tbl_instructor_groups
        WHERE instructor_id = ?
    `, instructorID)

	var group entity.InstructorGroup
	err := row.Scan(&group.GroupID, &group.InstructorID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &group, nil
}
