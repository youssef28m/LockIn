package core

import "time"

// TODO: Add session functionality

type Session struct {
    ID              int64
    StartTime       int64
    DurationSeconds int64
    Active          bool
}

func (s *Session) Remaining() int64 {
	now := time.Now().Unix()
	end := s.StartTime + int64(s.DurationSeconds)
	remaining := end - now

	if remaining < 0 {
		remaining = 0
	}
	return remaining
}

func (s *Session) Expired() bool {
	return s.Remaining() == 0
}


func (s *Session) RemainingMinutes() int64 {
	return s.Remaining() / 60
}

func (s *Session) RemainingHours() int64 {
	return s.Remaining() / 3600
}

func (s *Session) Start() {
	s.Active = true
	s.StartTime = time.Now().Unix()
}

func (s *Session) Stop() {
	s.Active = false
}