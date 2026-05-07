package core

import (
	"database/sql"
	"github.com/youssef28m/LockIn/internal/blocker"
	"github.com/youssef28m/LockIn/internal/storage"
	"time"
)

func InitializeScheduler(db *sql.DB, trigger <-chan struct{}, stop <-chan struct{}) {

	// On startup, check for any active sessions and block sites/apps
	checkUnblockedSites(db)
	blockAppsIfActive(db)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkSessionExpiration(db)
			blockAppsIfActive(db)
		case <-trigger:
			checkUnblockedSites(db)
			checkSessionExpiration(db)
			blockAppsIfActive(db)
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

	expiredFound := false
	for _, session := range sessions {
		if session.Active && session.Expired() {
			session.Stop()

			blocker.UnblockWebsites(db)
			storage.UpdateSession(db, session)
			expiredFound = true
		}
	}

	if expiredFound {
		storage.DeleteExpiredSessions(db)
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

func blockAppsIfActive(db *sql.DB) {
	sessions, err := storage.GetAllSessions(db)
	if err != nil {
		return
	}
	for _, session := range sessions {
		if session.Active && !session.Expired() {
			blocker.BlockApps(db)
			return
		}
	}
}
