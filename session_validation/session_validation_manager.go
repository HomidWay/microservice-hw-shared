package sessionvalidation

type SessionRepository interface {
	CreateNewSession(session *Session) error
	GetSessionData(sessionID, userID string) (*Session, error)
}
