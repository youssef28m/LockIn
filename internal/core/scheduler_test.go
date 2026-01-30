package core

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/youssef28m/LockIn/internal/models"
	"github.com/youssef28m/LockIn/internal/storage"
)

// setupTestDB creates a test database with all required tables
func setupSchedulerTestDB(t *testing.T) *sql.DB {
	home, _ := os.UserHomeDir()
	testDbPath := filepath.Join(home, ".lockin", "test_scheduler_LockIn.db")

	// Remove test db if it exists
	os.Remove(testDbPath)

	// Create new test database
	db, err := sql.Open("sqlite3", testDbPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create tables
	sessionsSQL := `CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		start_time INTEGER NOT NULL,
		duration_seconds INTEGER NOT NULL,
		active INTEGER NOT NULL
	);`

	blockedSitesSQL := `CREATE TABLE IF NOT EXISTS blocked_sites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain TEXT NOT NULL
	);`

	blockedAppsSQL := `CREATE TABLE IF NOT EXISTS blocked_apps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		process_name TEXT NOT NULL
	);`

	db.Exec(sessionsSQL)
	db.Exec(blockedSitesSQL)
	db.Exec(blockedAppsSQL)

	return db
}

// cleanupSchedulerTestDB removes the test database
func cleanupSchedulerTestDB(t *testing.T, db *sql.DB) {
	db.Close()
	home, _ := os.UserHomeDir()
	testDbPath := filepath.Join(home, ".lockin", "test_scheduler_LockIn.db")
	os.Remove(testDbPath)
}

// TestSessionExpiration tests if a session correctly identifies when it has expired
func TestSessionExpiration(t *testing.T) {
	db := setupSchedulerTestDB(t)
	defer cleanupSchedulerTestDB(t, db)

	// Create a session that expires immediately
	startTime := time.Now().Unix() - 10 // Started 10 seconds ago
	durationSeconds := 5                // Duration is 5 seconds
	active := true

	sessionID, err := storage.CreateSession(db, startTime, durationSeconds, active)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	session, err := storage.GetSessionByID(db, sessionID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	// Verify session has expired
	if !session.Expired() {
		t.Error("Session should have expired")
	}
	t.Logf("✓ Session correctly identified as expired (remaining: %d seconds)", session.Remaining())
}

// TestSessionNotExpired tests if an active, non-expired session is correctly identified
func TestSessionNotExpired(t *testing.T) {
	db := setupSchedulerTestDB(t)
	defer cleanupSchedulerTestDB(t, db)

	// Create a session that won't expire soon
	startTime := time.Now().Unix()
	durationSeconds := 3600 // 1 hour
	active := true

	sessionID, err := storage.CreateSession(db, startTime, durationSeconds, active)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	session, err := storage.GetSessionByID(db, sessionID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	// Verify session has NOT expired
	if session.Expired() {
		t.Error("Session should not have expired")
	}

	if !session.Active {
		t.Error("Session should be active")
	}

	remaining := session.Remaining()
	t.Logf("✓ Session is active with %d seconds remaining", remaining)

	if remaining <= 0 || remaining > int64(durationSeconds) {
		t.Errorf("Remaining time should be between 0 and %d, got %d", durationSeconds, remaining)
	}
}

// TestSessionStop tests stopping a session
func TestSessionStop(t *testing.T) {
	db := setupSchedulerTestDB(t)
	defer cleanupSchedulerTestDB(t, db)

	startTime := time.Now().Unix()
	sessionID, _ := storage.CreateSession(db, startTime, 3600, true)

	session, _ := storage.GetSessionByID(db, sessionID)

	// Verify session is initially active
	if !session.Active {
		t.Error("Session should be active initially")
	}

	// Stop the session
	session.Stop()
	if session.Active {
		t.Error("Session should be inactive after stopping")
	}

	// Update in database
	err := storage.UpdateSession(db, *session)
	if err != nil {
		t.Fatalf("Failed to update session: %v", err)
	}

	// Verify in database
	stoppedSession, _ := storage.GetSessionByID(db, sessionID)
	if stoppedSession.Active {
		t.Error("Stopped session should be inactive in database")
	}

	t.Logf("✓ Session successfully stopped")
}

// TestMultipleActiveSessionsWithExpired tests scheduler logic with multiple sessions
func TestMultipleActiveSessionsWithExpired(t *testing.T) {
	db := setupSchedulerTestDB(t)
	defer cleanupSchedulerTestDB(t, db)

	now := time.Now().Unix()

	// Create sessions: 1 expired, 2 active
	expiredSessionID, _ := storage.CreateSession(db, now-100, 50, true) // Expired
	_, _ = storage.CreateSession(db, now, 3600, true)                   // Active, 1 hour
	_, _ = storage.CreateSession(db, now, 1800, true)                   // Active, 30 minutes

	t.Logf("✓ Created 3 sessions: 1 expired, 2 active")

	// Retrieve all sessions
	sessions, err := storage.GetAllSessions(db)
	if err != nil {
		t.Fatalf("Failed to get sessions: %v", err)
	}

	if len(sessions) != 3 {
		t.Errorf("Expected 3 sessions, got %d", len(sessions))
	}

	// Count expired vs active
	expiredCount := 0
	activeCount := 0

	for _, session := range sessions {
		if session.Active && session.Expired() {
			expiredCount++
			t.Logf("  - Session %d: EXPIRED", session.ID)
		} else if session.Active && !session.Expired() {
			activeCount++
			t.Logf("  - Session %d: ACTIVE (%d seconds remaining)", session.ID, session.Remaining())
		}
	}

	if expiredCount != 1 {
		t.Errorf("Expected 1 expired session, found %d", expiredCount)
	}

	if activeCount != 2 {
		t.Errorf("Expected 2 active sessions, found %d", activeCount)
	}

	// Simulate scheduler: stop expired session
	expiredSession, _ := storage.GetSessionByID(db, expiredSessionID)
	if !expiredSession.Expired() {
		t.Fatal("Test session should be expired")
	}

	expiredSession.Stop()
	storage.UpdateSession(db, *expiredSession)
	t.Logf("✓ Stopped expired session")

	// Verify
	stoppedSession, _ := storage.GetSessionByID(db, expiredSessionID)
	if stoppedSession.Active {
		t.Error("Expired session should now be inactive")
	}

	t.Logf("✓ Session lifecycle correctly managed")
}

// TestSessionWithBlockedSites tests sessions with associated blocked sites
func TestSessionWithBlockedSites(t *testing.T) {
	db := setupSchedulerTestDB(t)
	defer cleanupSchedulerTestDB(t, db)

	// Create an active session
	startTime := time.Now().Unix()
	sessionID, _ := storage.CreateSession(db, startTime, 3600, true)
	t.Logf("✓ Created active session: %d", sessionID)

	// Add blocked sites
	sites := []string{"facebook.com", "twitter.com", "youtube.com"}
	for _, domain := range sites {
		storage.CreateBlockedSite(db, domain)
	}
	t.Logf("✓ Added %d blocked sites", len(sites))

	// Retrieve all blocked sites
	blockedSites, _ := storage.GetAllBlockedSites(db)
	if len(blockedSites) != 3 {
		t.Errorf("Expected 3 blocked sites, got %d", len(blockedSites))
	}

	// Get session and verify it's active with blocked sites
	session, _ := storage.GetSessionByID(db, sessionID)
	if !session.Active {
		t.Error("Session should be active")
	}

	if session.Expired() {
		t.Error("Session should not be expired")
	}

	t.Logf("✓ Session %d is active with %d blocked sites", sessionID, len(blockedSites))
}

// TestSessionRemainingTime tests remaining time calculations
func TestSessionRemainingTime(t *testing.T) {
	db := setupSchedulerTestDB(t)
	defer cleanupSchedulerTestDB(t, db)

	// Create a session that started 10 minutes ago and lasts 60 minutes
	startTime := time.Now().Unix() - 600 // 10 minutes ago
	durationSeconds := 3600              // 60 minutes
	active := true

	sessionID, _ := storage.CreateSession(db, startTime, durationSeconds, active)
	session, _ := storage.GetSessionByID(db, sessionID)

	remaining := session.Remaining()
	remainingMinutes := session.RemainingMinutes()
	remainingHours := session.RemainingHours()

	t.Logf("✓ Session remaining: %d seconds, %d minutes, %d hours", remaining, remainingMinutes, remainingHours)

	// Verify calculations are reasonable
	// Should have around 50 minutes remaining (3600 - 600 = 3000 seconds)
	expectedMin := int64(2900) // Allow 100 second margin
	expectedMax := int64(3100)

	if remaining < expectedMin || remaining > expectedMax {
		t.Errorf("Remaining time should be around 3000 seconds, got %d", remaining)
	}

	if remainingMinutes < 48 || remainingMinutes > 52 {
		t.Errorf("Remaining minutes should be around 50, got %d", remainingMinutes)
	}
}

// TestSessionStart tests the Start method
func TestSessionStart(t *testing.T) {
	db := setupSchedulerTestDB(t)
	defer cleanupSchedulerTestDB(t, db)

	// Create an inactive session
	session := &models.Session{
		StartTime:       0,
		DurationSeconds: 3600,
		Active:          false,
	}

	// Start the session
	beforeStart := time.Now().Unix()
	session.Start()
	afterStart := time.Now().Unix()

	if !session.Active {
		t.Error("Session should be active after Start()")
	}

	if session.StartTime < beforeStart || session.StartTime > afterStart {
		t.Errorf("StartTime should be between %d and %d, got %d", beforeStart, afterStart, session.StartTime)
	}

	t.Logf("✓ Session started at: %d (current time: %d)", session.StartTime, time.Now().Unix())
}

// TestSchedulerSessionFiltering tests filtering sessions for scheduler operations
func TestSchedulerSessionFiltering(t *testing.T) {
	db := setupSchedulerTestDB(t)
	defer cleanupSchedulerTestDB(t, db)

	now := time.Now().Unix()

	// Create various session states
	storage.CreateSession(db, now-100, 50, true)    // Expired, active
	storage.CreateSession(db, now, 3600, true)      // Active, not expired
	storage.CreateSession(db, now, 3600, false)     // Inactive (stopped)
	storage.CreateSession(db, now-1000, 500, false) // Expired, inactive

	sessions, _ := storage.GetAllSessions(db)

	// Filter for active + expired (what scheduler would unblock)
	var expiredActiveSessions []models.Session
	for _, session := range sessions {
		if session.Active && session.Expired() {
			expiredActiveSessions = append(expiredActiveSessions, session)
		}
	}

	if len(expiredActiveSessions) != 1 {
		t.Errorf("Expected 1 active+expired session, got %d", len(expiredActiveSessions))
	}

	// Filter for active + not expired (what scheduler would block)
	var activeNotExpiredSessions []models.Session
	for _, session := range sessions {
		if session.Active && !session.Expired() {
			activeNotExpiredSessions = append(activeNotExpiredSessions, session)
		}
	}

	if len(activeNotExpiredSessions) != 1 {
		t.Errorf("Expected 1 active+not-expired session, got %d", len(activeNotExpiredSessions))
	}

	t.Logf("✓ Correctly filtered %d expired+active and %d active+not-expired sessions",
		len(expiredActiveSessions), len(activeNotExpiredSessions))
}

// TestSchedulerInitialization tests the initialization logic
func TestSchedulerInitialization(t *testing.T) {
	db := setupSchedulerTestDB(t)
	defer cleanupSchedulerTestDB(t, db)

	// Create some active sessions
	now := time.Now().Unix()
	storage.CreateSession(db, now, 3600, true)
	storage.CreateSession(db, now, 1800, true)

	// Add blocked sites
	storage.CreateBlockedSite(db, "distraction.com")
	storage.CreateBlockedSite(db, "procrastination.com")

	sessions, _ := storage.GetAllSessions(db)
	blockedSites, _ := storage.GetAllBlockedSites(db)

	activeSessions := 0
	for _, session := range sessions {
		if session.Active && !session.Expired() {
			activeSessions++
		}
	}

	t.Logf("✓ Scheduler initialization state:")
	t.Logf("  - Total sessions: %d", len(sessions))
	t.Logf("  - Active sessions: %d", activeSessions)
	t.Logf("  - Blocked sites: %d", len(blockedSites))

	if activeSessions < 1 {
		t.Error("Should have at least 1 active session")
	}

	if len(blockedSites) < 1 {
		t.Error("Should have at least 1 blocked site")
	}
}
