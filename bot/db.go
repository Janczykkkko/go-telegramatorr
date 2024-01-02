package bot

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SessionsByUser []struct {
	UserName string
	Sessions []ActiveSession
}

func DbExistsAndWorks(dblocation string) bool {
	if _, err := os.Stat(dblocation); os.IsNotExist(err) {
		return false
	}

	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		log.Println("Error opening database:", err)
		return false
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Println("Error creating test table:", err)
		return false
	}

	defer func() {
		_, err := db.Exec("DROP TABLE IF EXISTS test")
		if err != nil {
			log.Println("Error dropping test table:", err)
		}
	}()

	return true
}

func CreateDb(dblocation string) error {
	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS streams (
			user_name TEXT,
			item_name INTEGER,
			playback_method TEXT,
			service_name TEXT,
			device_name TEXT,
			substream TEXT,
			bitrate TEXT,
			started_at TEXT,
			ended_at TEXT,
			stream_id TEXT PRIMARY KEY,
			duration_minutes TEXT
		);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Database created successfully")
	return nil
}

func CleanDB(dblocation string) error {
	log.Println("Performing scheduled db clean after a week...")
	err := os.Remove(dblocation)
	if err != nil {
		log.Printf("Failed to remove db: %s", err)
		return err
	}
	log.Printf("Database file %s deleted successfully", dblocation)
	err = CreateDb(dblocation)
	if err != nil {
		log.Printf("Failed to create db: %s", err)
		return err
	}
	return nil
}

func InsertDataToDb(session ActiveSession, endTime time.Time, dblocation string) error {
	if !CheckAndUpdateEntry(session, endTime, dblocation) {
		err := InsertDataToDb(session, endTime, dblocation)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertNewDataToDb(session ActiveSession, endTime time.Time, dblocation string) error {
	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	insertSQL := `
        INSERT INTO streams (
            user_name, item_name, playback_method, service_name,
            device_name, substream, bitrate, started_at, ended_at, stream_id, duration_minutes
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	startedAtStr := session.StartTime.Format(time.RFC3339)
	endedAtStr := endTime.Format(time.RFC3339)

	duration := endTime.Sub(session.StartTime)
	durationMinutes := int(duration.Minutes())

	_, err = db.Exec(insertSQL, session.UserName, session.Name, session.PlayMethod,
		session.Service, session.DeviceName, session.SubStream, session.Bitrate,
		startedAtStr, endedAtStr, session.ID, durationMinutes)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}

	log.Printf("Stream %s data persisted successfully in db", session.ID)
	return nil
}

func CheckAndUpdateEntry(session ActiveSession, endTime time.Time, dblocation string) bool {
	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		log.Printf("error opening database: %v", err)
		return false
	}
	defer db.Close()

	checkAndUpdateSQL := `
        SELECT ended_at, duration_minutes 
        FROM streams 
        WHERE user_name = ? AND item_name = ? AND ended_at >= ? 
        ORDER BY ended_at DESC 
        LIMIT 1
    `

	endedAtThreshold := endTime.Add(-30 * time.Minute)
	rows, err := db.Query(checkAndUpdateSQL, session.UserName, session.Name, endedAtThreshold)
	if err != nil {
		log.Printf("error querying database: %v", err)
		return false
	}
	defer rows.Close()

	var endedAtStr string
	var durationMinutes int

	if rows.Next() {
		if err := rows.Scan(&endedAtStr, &durationMinutes); err != nil {
			log.Printf("error scanning rows: %v", err)
			return false
		}

		prevEndedAt, err := time.Parse(time.RFC3339, endedAtStr)
		if err != nil {
			log.Printf("error parsing time: %v", err)
			return false
		}

		// Update the found entry
		newDuration := endTime.Sub(prevEndedAt)
		newDurationMinutes := int(newDuration.Minutes())

		updateSQL := `
            UPDATE streams 
            SET ended_at = ?, duration_minutes = ? 
            WHERE user_name = ? AND item_name = ? AND ended_at = ?
        `
		_, err = db.Exec(updateSQL, endTime.Format(time.RFC3339), durationMinutes+newDurationMinutes,
			session.UserName, session.Name, endedAtStr)
		if err != nil {
			log.Printf("error updating database: %v", err)
			return false
		}

		log.Printf("Entry updated successfully in db: %s", session.Name)
		return true
	}

	return false, nil
}

func GetSessionsByUserFromDB(dblocation string, timeframe int) (SessionsByUser, error) {
	db, err := sql.Open("sqlite3", dblocation)
	if err != nil {
		log.Println("Error accessing db", err)
		return nil, err
	}
	defer db.Close()

	currentTime := time.Now()
	timeFrame := currentTime.Add(-time.Duration(timeframe) * time.Hour)

	query := `
        SELECT user_name, item_name, playback_method, service_name,
            device_name, substream, bitrate, started_at, ended_at, duration_minutes, stream_id
        FROM streams
        WHERE ended_at >= ?
        ORDER BY user_name
    `

	rows, err := db.Query(query, timeFrame.Format("2006-01-02T15:04:05-07:00"))
	if err != nil {
		log.Println("Error querying db", err)
		return nil, err
	}
	defer rows.Close()

	sessionsByUser := make(SessionsByUser, 0)

	for rows.Next() {
		var session ActiveSession
		err := rows.Scan(
			&session.UserName,
			&session.Name,
			&session.PlayMethod,
			&session.Service,
			&session.DeviceName,
			&session.SubStream,
			&session.Bitrate,
			&session.StartTimeStr,
			&session.EndTimeStr,
			&session.Duration,
			&session.ID,
		)
		if err != nil {
			log.Println("Error processing data from db", err)
			continue
		}

		startTime, err := time.Parse("2006-01-02T15:04:05-07:00", session.StartTimeStr)
		if err != nil {
			log.Println("Error retireiving (parsing) time from db", err)
		}
		session.StartTime = startTime
		endTime, err := time.Parse("2006-01-02T15:04:05-07:00", session.EndTimeStr)
		if err != nil {
			log.Println("Error retireiving (parsing) time from db", err)
		}
		session.EndTime = endTime

		userIndex := -1
		for i := range sessionsByUser {
			if sessionsByUser[i].UserName == session.UserName {
				userIndex = i
				break
			}
		}

		if userIndex == -1 {
			sessionsByUser = append(sessionsByUser, struct {
				UserName string
				Sessions []ActiveSession
			}{
				UserName: session.UserName,
				Sessions: []ActiveSession{session},
			})
		} else {
			sessionsByUser[userIndex].Sessions = append(sessionsByUser[userIndex].Sessions, session)
		}
	}
	//organise chronologically per user
	for i := range sessionsByUser {
		sort.SliceStable(sessionsByUser[i].Sessions, func(j, k int) bool {
			return sessionsByUser[i].Sessions[j].EndTime.Before(sessionsByUser[i].Sessions[k].EndTime)
		})
	}

	if err := rows.Err(); err != nil {
		log.Println("Error preparing data (scanning) from db", err)
		return nil, err
	}

	return sessionsByUser, nil
}
