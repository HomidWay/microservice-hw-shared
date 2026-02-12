package usermanagement

type UserRepository interface {
	RegisterNewUser(userName, fullName, password, passwordConfirm string, userRole UserRole) (*User, error)
	GetUserData(userID string) (*User, error)
	ChangeUserRole(userID, roleName string) error
}
