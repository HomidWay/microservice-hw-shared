package userserviceclient

type UserRole int

const (
	UserGroupUndefined UserRole = iota
	UserRoleFree
	UserRoleBasic
	UserRolePremium
	UserRoleTechSupport
)

type User struct {
	id       string
	username string
	fullName string
	userRole UserRole
}

func NewUser(id, name string, fullName string, userRole UserRole) *User {
	return &User{id: id, username: name, fullName: fullName, userRole: userRole}
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) FullName() string {
	return u.fullName
}

func (u *User) Role() UserRole {
	return u.userRole
}
