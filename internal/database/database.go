package database

import (
	"database/sql"
	"time"

	"github.com/HubertBel/lazyorg/internal/calendar"
	"github.com/HubertBel/lazyorg/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func (database *Database) InitDatabase(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	database.db = db

	return database.createTables()
}

func (database *Database) createTables() error {
	_, err := database.db.Exec(`
        CREATE TABLE IF NOT EXISTS events (
        id INTEGER NOT NULL PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT,
        location TEXT,
        time DATETIME NOT NULL,
        duration REAL NOT NULL,
        frequency INTEGER,
        occurence INTEGER
    )`)
	if err != nil {
		return err
	}

	_, err = database.db.Exec(`
        CREATE TABLE IF NOT EXISTS notes (
        id INTEGER NOT NULL PRIMARY KEY,
        content TEXT NOT NULL,
        time DATETIME NOT NULL
    )`)
	if err != nil {
		return err
	}

	return nil
}

func (database *Database) AddEvent(event calendar.Event) (int, error) {
	result, err := database.db.Exec(`
        INSERT INTO events (
            name, description, location, time, duration, frequency, occurence
        ) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		event.Name,
		event.Description,
		event.Location,
		event.Time,
		event.DurationHour,
		event.FrequencyDay,
		event.Occurence,
	)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), err
}

func (database *Database) GetEventById(id int) (*calendar.Event, error) {
    rows, err := database.db.Query(`
        SELECT * FROM events WHERE id = ?`,
        id,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    if rows.Next() {
        var event calendar.Event
        if err := rows.Scan(
            &event.Id,
            &event.Name,
            &event.Description,
            &event.Location,
            &event.Time,
            &event.DurationHour,
            &event.FrequencyDay,
            &event.Occurence,
        ); err != nil {
            return nil, err
        }
        return &event, nil
    }

    return nil, nil
}

func (database *Database) GetEventsByDate(date time.Time) ([]*calendar.Event, error) {
	formattedDate := utils.FormatDate(date)

	rows, err := database.db.Query(`
        SELECT * FROM events WHERE date(time) = ?`,
		formattedDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*calendar.Event
	for rows.Next() {
		var event calendar.Event

		if err := rows.Scan(
			&event.Id,
			&event.Name,
			&event.Description,
			&event.Location,
			&event.Time,
			&event.DurationHour,
			&event.FrequencyDay,
			&event.Occurence,
		); err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	return events, nil
}

func (database *Database) DeleteEventById(id int) error {
	_, err := database.db.Exec("DELETE FROM events WHERE id = ?", id)
	return err
}

func (database *Database) DeleteEventsByName(name string) error {
	_, err := database.db.Exec("DELETE FROM events WHERE name = ?", name)
	return err
}

func (database *Database) UpdateEventById(id int, event *calendar.Event) error {
	return nil
}

func (database *Database) UpdateEventByName(name string) error {
	return nil
}

func (database *Database) SaveNote(content string) error {
	_, err := database.db.Exec("DELETE FROM notes")
	if err != nil {
		return err
	}

	_, err = database.db.Exec(`INSERT INTO notes (
            content, time
        ) VALUES (?, datetime('now'))`, content)

	return err
}

func (database *Database) GetLatestNote() (string, error) {
	var content string
	err := database.db.QueryRow(
		"SELECT content FROM notes ORDER BY time DESC LIMIT 1",
	).Scan(&content)

	return content, err
}

func (database *Database) CloseDatabase() error {
	if database.db == nil {
		return nil
	}

	return database.db.Close()
}
