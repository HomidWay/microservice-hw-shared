package session_validation

import (
	"errors"
	"time"
)

var (
	ErrSessionNotFound           = errors.New("session not found")
	ErrSessionBelongsToOtherUser = errors.New("session belongs to another user")
)

type SessionState int

const (
	SessionStateUndefined SessionState = iota
	SessionStateActive
	SessionStateSuspended
)

type Session struct {
	sessionID    string
	userID       string
	sessionState SessionState
	validUntil   time.Time
}

func NewSession(sessionID, userID string, sessionState SessionState, validUntil time.Time) *Session {
	return &Session{
		sessionID:    sessionID,
		userID:       userID,
		sessionState: sessionState,
		validUntil:   validUntil,
	}
}

func (s *Session) SessionID() string {
	return s.sessionID
}

func (s *Session) UserID() string {
	return s.userID
}

func (s *Session) ValidUntil() time.Time {
	return s.validUntil
}

func (s *Session) SessionState() SessionState {
	return s.sessionState
}

func (s *Session) IsValid() bool {
	return s.validUntil.After(time.Now()) && s.sessionState == SessionStateActive
}
