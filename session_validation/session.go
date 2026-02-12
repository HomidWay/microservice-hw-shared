package sessionvalidation

import (
	"sync"
	"time"
)

type SessionState int

const (
	SessionStateUndefined SessionState = iota
	SessionStateActive
	SessionStateInactive
	SessionStateSuspended
)

type Session struct {
	userID    string
	sessionID string

	validUntil   time.Time
	sessionState SessionState

	mu sync.RWMutex
}

func NewSession(valid time.Time, userID, token string, state SessionState) *Session {
	return &Session{
		validUntil:   valid,
		userID:       userID,
		sessionID:    token,
		sessionState: state,
	}
}

func (s *Session) ValidUntil() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.validUntil
}

func (s *Session) IsValid() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.validUntil.After(time.Now()) && s.sessionState == SessionStateActive
}

func (s *Session) UserID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.userID
}

func (s *Session) SessionID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessionID
}

func (s *Session) SessionState() SessionState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessionState
}

func (s *Session) Renew(duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.validUntil = time.Now().Add(duration)
}
