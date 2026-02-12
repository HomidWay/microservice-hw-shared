package authorizationserviceclient

import "time"

type SessionState int

const (
	SessionStateUndefined SessionState = iota
	SessionStateActive
	SessionStateSuspended
)

// MARK: Authorization Response
type AuthorizationResponse struct {
	sessionID string
	userID    string
}

func NewAuthorizationResponse(sessionID string, userID string) *AuthorizationResponse {
	return &AuthorizationResponse{
		sessionID: sessionID,
		userID:    userID,
	}
}

func (ar *AuthorizationResponse) SessionID() string {
	return ar.sessionID
}

func (ar *AuthorizationResponse) UserID() string {
	return ar.userID
}

// MARK: Validate Session Response

type ValidateSessionResponse struct {
	sessionID  string
	userID     string
	state      SessionState
	validUntil time.Time
}

func NewValidateSessionResponse(sessionID, userID string, state SessionState, validUntil time.Time) *ValidateSessionResponse {
	return &ValidateSessionResponse{sessionID: sessionID, userID: userID, state: state, validUntil: validUntil}
}
