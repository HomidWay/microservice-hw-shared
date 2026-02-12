package usermanagement

import (
	"sync/atomic"
	"time"
)

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

	lastAccessed *atomic.Int64
}

func NewUser(id, name string, fullName string, userRole UserRole) *User {
	lastAccessed := &atomic.Int64{}
	lastAccessed.Store(time.Now().Unix())

	return &User{id: id, username: name, fullName: fullName, userRole: userRole, lastAccessed: lastAccessed}
}

func (u *User) ID() string {
	u.lastAccessed.Store(time.Now().Unix())
	return u.id
}

func (u *User) Username() string {
	u.lastAccessed.Store(time.Now().Unix())
	return u.username
}

func (u *User) FullName() string {
	u.lastAccessed.Store(time.Now().Unix())
	return u.fullName
}

func (u *User) Role() UserRole {
	u.lastAccessed.Store(time.Now().Unix())
	return u.userRole
}
