package db

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"telegramatorr/bot"
	"time"
)

func checkDb(dblocation string) (exists bool) {
	if _, err := os.Stat(dblocation); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func createDb(dblocation string) error {
	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS streams (
			user_name TEXT PRIMARY KEY,
			item_name INTEGER,
			playback_method TEXT,
			service_name TEXT,
			device_name TEXT,
			substream TEXT,
			bitrate TEXT,
			started_at TEXT,
			ended_at TEXT,
			stream_id TEXT
		);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Database created successfully")
	return nil
}

func cleanDB(dblocation string) error {
	log.Println("Performing scheduled db clean after a week...")
	err := os.Remove(dblocation)
	if err != nil {
		return err
	}
	log.Printf("Database file %s deleted successfully", dblocation)
	return nil
}

func getDBCreationTime(dblocation string) (time.Time, error) {
	fileInfo, err := os.Stat(dblocation)
	if err != nil {
		return time.Time{}, err
	}
	creationTime := fileInfo.ModTime()
	return creationTime, nil
}

func insertDataToDb(session bot.ActiveSession, endTime time.Time, dblocation string) error {
	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		return err
	}
	defer db.Close()

	insertSQL := `
		INSERT INTO streams (
			user_name, item_name, playback_method, service_name,
			device_name, substream, bitrate, started_at, ended_at, duration_minutes, stream_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	startedAtStr := session.StartTime.Format(time.RFC3339)
	endedAtStr := endTime.Format(time.RFC3339)

	duration := endTime.Sub(session.StartTime)
	durationMinutes := int(duration.Minutes())

	_, err = db.Exec(insertSQL, session.UserName, session.Name, session.PlayMethod,
		session.Service, session.DeviceName, session.SubStream, session.Bitrate,
		startedAtStr, endedAtStr, durationMinutes, session.ID)
	if err != nil {
		return err
	}

	log.Printf("Stream %s data persisted successfully in db", session.ID)
	return nil
}

func getSessionDataFromDB(dblocation string) ([]bot.ActiveSession, error) {
	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT user_name, item_name, playback_method, service_name,
			device_name, substream, bitrate, started_at, duration_minutes, ended_at, stream_id
		FROM streams
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []bot.ActiveSession

	for rows.Next() {
		var session bot.ActiveSession
		var startedAtStr, endedAtStr string
		var durationMinutes int

		err := rows.Scan(
			&session.UserName, &session.Name, &session.PlayMethod,
			&session.Service, &session.DeviceName, &session.SubStream,
			&session.Bitrate, &startedAtStr, &durationMinutes, &endedAtStr, &session.ID,
		)
		if err != nil {
			return nil, err
		}

		session.StartTime, err = time.Parse(time.RFC3339, startedAtStr)
		if err != nil {
			return nil, err
		}

		session.Duration = strconv.Itoa(durationMinutes)

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}
