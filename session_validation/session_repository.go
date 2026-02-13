package session_validation

type SessionRepository interface {
	ValidateSession(sessionID string) (*Session, error)
}
