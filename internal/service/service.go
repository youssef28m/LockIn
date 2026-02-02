package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/youssef28m/LockIn/internal/blocker"
	"github.com/youssef28m/LockIn/internal/models"
	"github.com/youssef28m/LockIn/internal/storage"
	"github.com/youssef28m/LockIn/internal/validator"
)

type AppService struct {
	db *sql.DB
	// This channel is the "trigger" we inject
	notifier chan struct{}
}

func NewAppService(db *sql.DB, n chan struct{}) *AppService {
	return &AppService{db: db, notifier: n}
}


func (s *AppService) CreateSession(duration int) error {

	validDuration := validator.IsValidDuration(duration)
	if !validDuration {
		return fmt.Errorf("invalid duration")
	}

	startTime := time.Now().Unix()
	_, err := storage.CreateSession(s.db, startTime, duration, true)
	if err != nil {
		return err
	}

	return nil
}

func (s *AppService) CreateAndStartSession(duration int) error {

	sessions, err := storage.GetAllSessions(s.db)
	if err != nil {
		return err
	}

	// check for any active session
	for _, session := range sessions {
		if session.Active && !session.Expired() {
			return fmt.Errorf("there is already an active session")
		}
	}

	err = s.CreateSession(duration)
	if err != nil {
		return err
	}

	// 3. Signal the Scheduler to wake up immediately
	// We use a non-blocking send so the UI doesn't hang if the scheduler is busy
	select {
	case s.notifier <- struct{}{}:
	default:
		// Scheduler is already busy/running, no need to queue another signal
	}
	return nil
}



func (s *AppService) AddBlockedSite(domain string) error {
	validDomain := validator.IsValidDomain(domain)
	if !validDomain {
		return fmt.Errorf("invalid domain format")
	}

	_, err := storage.CreateBlockedSite(s.db, domain)
	if err != nil {
		return err
	}

	// block the newly added site immediately if there is an active session
	haveActiveSession, err := HaveActiveSession(s.db)
	if err != nil {
		return err
	}
	if haveActiveSession {
		err := blocker.BlockSite(domain)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AppService) GetActiveSession() (models.Session, error) {
	sessions, err := storage.GetActiveSessions(s.db)
	if err != nil {
		return models.Session{}, err
	}

	if len(sessions) == 0 {
		return models.Session{}, fmt.Errorf("no active session found")
	}
	return sessions[0], nil

}

func HaveActiveSession(db *sql.DB) (bool, error) {
	sessions, err := storage.GetActiveSessions(db)
	if err != nil {
		return false, err
	}
	
	return len(sessions) > 0, nil

}

