package core

import (
	"database/sql"
	"github.com/youssef28m/LockIn/internal/blocker"
	"github.com/youssef28m/LockIn/internal/storage"
	"log"
	"time"
)

func InitializeScheduler(db *sql.DB, trigger chan struct{}) {

	// On startup, check for any active sessions and block sites
	checkUnblockedSites(db)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Periodic cleanup (unblocking expired sessions)
			checkSessionExpiration(db)
		case <-trigger:
			// Instant reaction
			checkUnblockedSites(db)
			checkSessionExpiration(db)
		}
	}
}

func checkSessionExpiration(db *sql.DB) {
	sessions, err := storage.GetAllSessions(db)
	if err != nil {
		log.Println("Error fetching sessions:", err)
		return
	}

	for _, session := range sessions {
		if session.Active && session.Expired() {
			session.Stop()

			// unblock websites/apps
			err := blocker.UnblockWebsites(db)
			if err != nil {
				log.Println("Error unblocking websites:", err)
			}

			err = storage.UpdateSession(db, session)
			if err != nil {
				log.Println("Error updating session:", err)
				continue
			}
		}
	}
}

func checkUnblockedSites(db *sql.DB) {
	sessions, err := storage.GetAllSessions(db)
	if err != nil {
		log.Println("Error fetching sessions:", err)
		return
	}
	for _, session := range sessions {
		if session.Active && !session.Expired() {
			// block websites/apps
			err := blocker.BlockWebsites(db)
			if err != nil {
				log.Println("Error blocking websites:", err)
			}
		}
	}
}
