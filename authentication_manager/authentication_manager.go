package authenticationmanager

type AuthenticationManager interface {
	AuthorizeUser(userName, password string) (string, error)
}
