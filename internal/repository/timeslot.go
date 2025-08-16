package repository

import (
	"github.com/massivemadness/schedule-service/internal/database"
	"github.com/massivemadness/schedule-service/internal/entity"
)

type TimeslotRepository interface {
	Save(slots []entity.TimeSlot) error
	GetByScheduleID(scheduleID int64) ([]entity.TimeSlot, error)
	BookSlot(timeslotID int64, userID int64, userName string) error
}

type timeslotRepositoryImpl struct {
	db *database.Database
}

func NewTimeslotRepository(db *database.Database) TimeslotRepository {
	return &timeslotRepositoryImpl{db: db}
}

func (r *timeslotRepositoryImpl) Save(slots []entity.TimeSlot) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, slot := range slots {
		r.db.Exec(`INSERT OR IGNORE INTO tbl_timeslots (schedule_id, time, user_id, user_name) VALUES (?, ?, ?, ?)`,
			slot.ScheduleID,
			slot.Time,
			slot.UserID,
			slot.UserName,
		)
	}

	return tx.Commit()
}

func (r *timeslotRepositoryImpl) GetByScheduleID(scheduleID int64) ([]entity.TimeSlot, error) {
	rows, err := r.db.Query(`
        SELECT id, schedule_id, time, user_id, user_name
        FROM tbl_timeslots
        WHERE schedule_id = ?
    `, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timeslots []entity.TimeSlot
	for rows.Next() {
		var slot entity.TimeSlot
		if err := rows.Scan(&slot.ID, &slot.ScheduleID, &slot.Time, &slot.UserID, &slot.UserName); err != nil {
			return nil, err
		}
		timeslots = append(timeslots, slot)
	}
	return timeslots, nil
}

func (r *timeslotRepositoryImpl) BookSlot(timeslotID int64, userID int64, userName string) error {
	_, err := r.db.Exec(`UPDATE tbl_timeslots SET user_id = ?, user_name = ? WHERE id = ?`,
		userID,
		userName,
		timeslotID,
	)
	return err
}
