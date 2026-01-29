package core

import (
	"database/sql"
	"log"
	"time"
	"github.com/youssef28m/LockIn/internal/storage"
)

// check for active sessions and manage them
// if session expired, unblock websites and apps

func InitializeScheduler(db *sql.DB) {
	// Process sessions

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		sessions, err := storage.GetAllSessions(db)
		if err != nil {
			log.Println("Error fetching sessions:", err)
			return
		}

		for _, session := range sessions {
			if session.Active && session.Expired() {
				session.Stop()
				
				// unblock websites/apps
				
				err := storage.UpdateSession(db, session)
				if err != nil {
					log.Println("Error updating session:", err)
                    continue
				}
			}
		}

	}
}
