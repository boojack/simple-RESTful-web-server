package model

import "time"

// User user
type User struct {
	ID        string
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}

// NewUser new user
func NewUser(id, username, password, email string) *User {

	return &User{
		// todo
	}
}
