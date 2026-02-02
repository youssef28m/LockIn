package service

import (
	"database/sql"
	"fmt"
	"github.com/youssef28m/LockIn/internal/blocker"
	"github.com/youssef28m/LockIn/internal/storage"
	"github.com/youssef28m/LockIn/internal/validator"
	"time"
)

type AppService struct {
	db *sql.DB
	// This channel is the "trigger" we inject
	notifier chan struct{}
}

func NewAppService(db *sql.DB, n chan struct{}) *AppService {
	return &AppService{db: db, notifier: n}
}

func (s *AppService) CreateAndStartSession(duration int) error {

	validDuration := validator.IsValidDuration(duration)
	if !validDuration {
		return fmt.Errorf("invalid duration")
	}
	
	startTime := time.Now().Unix()
	_, err := storage.CreateSession(s.db, startTime, duration, true)
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

func HaveActiveSession(db *sql.DB) (bool, error) {
	sessions, err := storage.GetAllSessions(db)
	if err != nil {
		return false, err
	}
	for _, session := range sessions {
		if session.Active && !session.Expired() {
			return true, nil
		}
	}
	return false, nil
}
