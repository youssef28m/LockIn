package core

import (
	"database/sql"
	"github.com/youssef28m/LockIn/internal/blocker"
	"github.com/youssef28m/LockIn/internal/storage"
	"time"
)

func InitializeScheduler(db *sql.DB, trigger <-chan struct{}, stop <-chan struct{}) {

	// On startup, check for any active sessions and block sites
	checkUnblockedSites(db)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkSessionExpiration(db)
		case <-trigger:
			checkUnblockedSites(db)
			checkSessionExpiration(db)
		case <-stop:
			return
		}
	}
}

func checkSessionExpiration(db *sql.DB) {
	sessions, err := storage.GetAllSessions(db)
	if err != nil {
		return
	}

	for _, session := range sessions {
		if session.Active && session.Expired() {
			session.Stop()

			blocker.UnblockWebsites(db)
			storage.UpdateSession(db, session)
		}
	}
}

func checkUnblockedSites(db *sql.DB) {
	sessions, err := storage.GetAllSessions(db)
	if err != nil {
		return
	}
	for _, session := range sessions {
		if session.Active && !session.Expired() {
			blocker.BlockWebsites(db)
		}
	}
}
