package core

import (
	"database/sql"
	"log"
	"time"
	"github.com/youssef28m/LockIn/internal/blocker"
	"github.com/youssef28m/LockIn/internal/storage"
)



func InitializeScheduler(db *sql.DB) {
	
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
}

