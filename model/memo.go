package model

import (
	"neosmemo/backend/helper"
	"neosmemo/backend/util"
	"time"
)

// Memo memo
type Memo struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateMemo Create Memo
func CreateMemo(userID string, content string) (Memo, error) {
	// TODO: need to think about the other way to generate id
	memo := Memo{
		ID:        util.GenUUID(),
		UserID:    userID,
		Content:   content,
		CreatedAt: util.GetNowTime(),
		UpdatedAt: util.GetNowTime(),
	}

	sqlStatement := `
		INSERT INTO memos (id, user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	if _, err := helper.DBService.Exec(sqlStatement, util.IterStructFieldValue(&memo)...); err != nil {
		return memo, err
	}

	return memo, nil
}

// UpdateMemoByID Update memo
func UpdateMemoByID(id string, userID string, content string) error {
	sqlStatement := `
		UPDATE memos
		SET content = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4
	`
	_, err := helper.DBService.Exec(sqlStatement, content, util.GetNowTime(), id, userID)

	if err != nil {
		return err
	}
	return nil
}

// DeleteMemoByID Delete memo
func DeleteMemoByID(id string, userID string) error {
	sqlStatement := `
		DELETE FROM memos
		WHERE id = $1 AND user_id = $2
	`

	_, err := helper.DBService.Exec(sqlStatement, id, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetMemoByID Get Memos By UserID
func GetMemoByID(id string) (Memo, error) {
	memo := Memo{}

	row := helper.DBService.QueryRow("SELECT * FROM memos WHERE id = $1", id)
	err := row.Scan(util.IterStructFieldAddr(&memo)...)
	if err != nil {
		return memo, err
	}

	return memo, nil
}

// GetMemosByUserID Get Memos By UserID
func GetMemosByUserID(userID string) ([]Memo, error) {
	memos := []Memo{}
	memo := Memo{}

	if rows, err := helper.DBService.Query("SELECT * FROM memos WHERE user_id = $1", userID); err == nil {
		for rows.Next() {
			rows.Scan(util.IterStructFieldAddr(&memo)...)
			memos = append(memos, memo)
		}
	} else {
		return memos, err
	}

	return memos, nil
}
