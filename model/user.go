package model

import "time"

// User user
// NOTE: 注意属性的先后顺序务必和数据库字段顺序保持一致
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewUser new user
func NewUser(id, username, password, email string) *User {

	return &User{
		// todo
	}
}
