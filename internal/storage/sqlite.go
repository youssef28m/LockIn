package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
	"github.com/youssef28m/LockIn/internal/models"
)

func Connect() *sql.DB {

	home, _ := os.UserHomeDir()
	dbDir := filepath.Join(home, ".lockin")
	dbPath := filepath.Join(dbDir, "LockIn.db")

	// Create directory if it doesn't exist
	err := os.MkdirAll(dbDir, 0755)
	if err != nil {
		log.Fatal("Error creating database directory:", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateDB() {
	db := Connect()
	defer db.Close()

	// Create sessions table
	sessionsSQL := `CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		start_time INTEGER NOT NULL,
		duration_seconds INTEGER NOT NULL,
		active INTEGER NOT NULL
	);`

	_, err := db.Exec(sessionsSQL)
	if err != nil {
		log.Fatal("Error creating sessions table:", err)
	}
	fmt.Println("Sessions table created successfully")

	// Create blocked_sites table
	blockedSitesSQL := `CREATE TABLE IF NOT EXISTS blocked_sites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id INTEGER NOT NULL,
		domain TEXT NOT NULL,
		FOREIGN KEY(session_id) REFERENCES sessions(id) ON DELETE CASCADE
	);`

	_, err = db.Exec(blockedSitesSQL)
	if err != nil {
		log.Fatal("Error creating blocked_sites table:", err)
	}
	fmt.Println("Blocked sites table created successfully")

	// Create blocked_apps table
	blockedAppsSQL := `CREATE TABLE IF NOT EXISTS blocked_apps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id INTEGER NOT NULL,
		process_name TEXT NOT NULL,
		FOREIGN KEY(session_id) REFERENCES sessions(id) ON DELETE CASCADE
	);`

	_, err = db.Exec(blockedAppsSQL)
	if err != nil {
		log.Fatal("Error creating blocked_apps table:", err)
	}
	fmt.Println("Blocked apps table created successfully")
}

//************************************************************//
// Session CRUD Operations
//************************************************************//

func CreateSession(db *sql.DB, startTime int64, durationSeconds int, active bool) (int64, error) {
	activeInt := 0
	if active {
		activeInt = 1
	}

	// Execute the insert
	result, err := db.Exec(
		`INSERT INTO sessions (start_time, duration_seconds, active)
		 VALUES (?, ?, ?)`,
		startTime,
		durationSeconds,
		activeInt,
	)
	if err != nil {
		return 0, err
	}

	// Get the ID of the inserted row
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetAllSessions(db *sql.DB) ([]models.Session, error) {
	rows, err := db.Query("SELECT id, start_time, duration_seconds, active FROM sessions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		var activeInt int
		err := rows.Scan(&session.ID, &session.StartTime, &session.DurationSeconds, &activeInt)
		session.Active = activeInt != 0

		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

func GetSessionByID(db *sql.DB, id int64) (*models.Session, error) {
	
	row := db.QueryRow("SELECT id, start_time, duration_seconds, active FROM sessions WHERE id = ?", id)
	var session models.Session
	var activeInt int
	err := row.Scan(&session.ID, &session.StartTime, &session.DurationSeconds, &activeInt)
	session.Active = activeInt != 0
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func UpdateSession(db *sql.DB, session models.Session) error {
	query := `
	UPDATE sessions
	SET start_time = ?, duration_seconds = ?, active = ?
	WHERE id = ?
	`

	result ,err := db.Exec(query, session.StartTime, session.DurationSeconds, session.Active, session.ID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", session.ID)
	}
	
	return nil

}
