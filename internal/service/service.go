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

	// Block all sites immediately so the user gets feedback if it fails
	if err := blocker.BlockWebsites(s.db); err != nil {
		return fmt.Errorf("session created but failed to block sites: %w", err)
	}

	// Signal the Scheduler to wake up for ongoing management
	select {
	case s.notifier <- struct{}{}:
	default:
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

func (s *AppService) GetBlockedSites() ([]string, error) {
	blockedSites, err := storage.GetAllBlockedSites(s.db)
	if err != nil {
		return nil, err
	}

	var sites []string
	for _, site := range blockedSites {
		sites = append(sites, site.Domain)
	}

	return sites, nil
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

func (s *AppService) AddBlockedApp(processName string) error {
	if processName == "" {
		return fmt.Errorf("process name cannot be empty")
	}

	_, err := storage.CreateBlockedApp(s.db, processName)
	if err != nil {
		return err
	}

	haveActiveSession, err := HaveActiveSession(s.db)
	if err != nil {
		return err
	}
	if haveActiveSession {
		_ = blocker.BlockApps(s.db)
	}

	return nil
}

func (s *AppService) GetBlockedApps() ([]string, error) {
	blockedApps, err := storage.GetAllBlockedApps(s.db)
	if err != nil {
		return nil, err
	}

	var apps []string
	for _, app := range blockedApps {
		apps = append(apps, app.ProcessName)
	}
	return apps, nil
}

func (s *AppService) GetSessionHistory() ([]models.Session, error) {
	allSessions, err := storage.GetAllSessions(s.db)
	if err != nil {
		return nil, err
	}
	var completed []models.Session
	for _, session := range allSessions {
		if !session.Active {
			completed = append(completed, session)
		}
	}
	return completed, nil
}

func (s *AppService) RemoveBlockedApp(processName string) error {
	apps, err := storage.GetAllBlockedApps(s.db)
	if err != nil {
		return err
	}

	found := false
	var appId int64
	for _, app := range apps {
		if app.ProcessName == processName {
			appId = app.ID
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("app not found in blocked list")
	}

	return storage.DeleteBlockedApp(s.db, appId)
}

func (s *AppService) RemoveBlockedSite(domain string) error {

	var siteId int64

	sites, err := storage.GetAllBlockedSites(s.db)
	if err != nil {
		return err
	}

	// Check if the site is in the blocked list
	found := false
	for _, site := range sites {
		if site.Domain == domain {
			siteId = site.ID
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("site not found in blocked list")
	}


	err = storage.DeleteBlockedSite(s.db, siteId)
	if err != nil {
		return err
	}

	// unblock the site immediately if there is an active session
	haveActiveSession, err := HaveActiveSession(s.db)
	if err != nil {
		return err
	}
	if haveActiveSession {
		err := blocker.UnblockSite(domain)
		if err != nil {
			return err
		}
	}
	
	return nil
}