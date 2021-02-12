package model

import (
	"neosmemo/backend/helper"
	"neosmemo/backend/util"
	"time"
)

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

// CreateUser CreateUser
func CreateUser(username string, password string, email string) (User, error) {
	user := User{
		ID:        util.GenUUID(),
		Username:  username,
		Password:  password,
		Email:     email,
		CreatedAt: util.GetNowTime(),
		UpdatedAt: util.GetNowTime(),
	}

	sqlStatement := `
		INSERT INTO users (id, username, password, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := helper.DBService.Exec(sqlStatement, util.IterStructFieldValue(&user)...)

	return user, err
}

// GetAllUsers GetAllUsers
func GetAllUsers() ([]User, error) {
	users := []User{}
	user := User{}

	rows, err := helper.DBService.Query("SELECT * FROM users")

	for rows.Next() {
		rows.Scan(util.IterStructFieldAddr(&user)...)
		users = append(users, user)
	}

	return users, err
}

// GetUserInfoByID GetUserInfoByID
func GetUserInfoByID(userID string) (User, error) {
	user := User{}

	row := helper.DBService.QueryRow("SELECT * FROM users WHERE id = $1", userID)
	err := row.Scan(util.IterStructFieldAddr(&user)...)

	return user, err
}

// TODO: UpdateUserInfo
// func UpdateUserInfo

// ValidUserSignin ValidSignin
func ValidUserSignin(username string, password string) (User, error) {
	user := User{}
	row := helper.DBService.QueryRow("SELECT * FROM users WHERE username = $1 AND password = $2", username, password)
	err := row.Scan(util.IterStructFieldAddr(&user)...)

	return user, err
}

// CheckUsernameUsable CheckUsernameUsable
func CheckUsernameUsable(username string) bool {
	count := 1
	err := helper.DBService.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)

	if err != nil || count != 0 {
		return false
	}

	return true
}
