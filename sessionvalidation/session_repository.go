package sessionvalidation

type SessionRepository interface {
	ValidateSession(sessionID string) (*Session, error)
}
