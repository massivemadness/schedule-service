package database

import (
	"database/sql"

	"github.com/massivemadness/schedule-service/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

func New(cfg *config.Config) (*Database, error) {
	db, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		return nil, err
	}

	// if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
	// 	log.Fatalf("Ошибка включения внешних ключей: %v", err)
	// }

	schema := `
	-- Инструкторы
	CREATE TABLE IF NOT EXISTS tbl_instructors (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);

	-- Привязанные группы
	CREATE TABLE IF NOT EXISTS tbl_instructor_groups (
		group_id INTEGER PRIMARY KEY,
		instructor_id INTEGER NOT NULL,
		FOREIGN KEY (instructor_id) REFERENCES tbl_instructors (id)
	);

	-- Форма создания расписания
	CREATE TABLE IF NOT EXISTS tbl_schedule_form (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		instructor_id INTEGER NOT NULL,
		date TEXT, -- формат YYYY-MM-DD
		timeslots TEXT -- слоты через запятую: "07:00,07:30"
	);

	-- Расписание
	CREATE TABLE IF NOT EXISTS tbl_schedules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		instructor_id INTEGER NOT NULL,
		message_id INTEGER,
		date TEXT NOT NULL, -- формат YYYY-MM-DD
		FOREIGN KEY (instructor_id) REFERENCES tbl_instructors (id)
	);

	-- Таймслоты
	CREATE TABLE IF NOT EXISTS tbl_timeslots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		schedule_id INTEGER NOT NULL,
		time TEXT NOT NULL, -- формат HH:MM
		user_id INTEGER,
		user_name TEXT,
		FOREIGN KEY (schedule_id) REFERENCES tbl_schedules (id) ON DELETE CASCADE
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}
