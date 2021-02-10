package memo

import (
	"encoding/json"
	"fmt"
	"neosmemo/backend/handler"
	"neosmemo/backend/helper"
	"neosmemo/backend/model"
	"neosmemo/backend/util"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetAllMemos get
func GetAllMemos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, ok := helper.GetUserIDFromSession(r)
	if !ok {
		panic("You have not sign in")
	}

	memos := []model.Memo{}
	memo := model.Memo{}

	if rows, err := helper.DBService.Query("SELECT * FROM memos WHERE user_id = $1", userID); err == nil {
		for rows.Next() {
			rows.Scan(util.IterStructFieldAddr(&memo)...)
			memos = append(memos, memo)
		}
	} else {
		panic("fetch failed, try later plz")
	}

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "memo get all succeed",
		Succeed:       true,
		Data:          &memos,
	})
}

// GetMemoByID do not need user session
func GetMemoByID(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	memo := model.Memo{}

	row := helper.DBService.QueryRow("SELECT * FROM memos WHERE id = $1", id)
	err := row.Scan(util.IterStructFieldAddr(&memo)...)
	if err != nil {
		// no rows in result set
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "memo get succeed",
		Succeed:       true,
		Data:          &memo,
	})
}

// CreateMemo post
func CreateMemo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, ok := helper.GetUserIDFromSession(r)
	if !ok {
		panic("You have not sign in")
	}

	t := struct {
		Content string
	}{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&t); err != nil {
		panic("request data type error")
	}

	// TODO: need to polish the way to generate id
	memo := model.Memo{
		ID:        util.GenUUID(),
		UserID:    userID,
		Content:   t.Content,
		CreatedAt: util.GetNowTime(),
		UpdatedAt: util.GetNowTime(),
	}

	sqlStatement := `
		INSERT INTO memos (id, user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	if _, err := helper.DBService.Exec(sqlStatement, util.IterStructFieldValue(&memo)...); err != nil {
		panic("create memo error")
	}

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "memo create succeed",
		Succeed:       true,
		Data:          &memo,
	})
}

// UpdateMemo post
func UpdateMemo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, ok := helper.GetUserIDFromSession(r)
	if !ok {
		panic("You have not sign in")
	}

	t := struct {
		ID      string
		Content string
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		panic("request data type error")
	}

	sqlStatement := `
		UPDATE memos
		SET content = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4
	`
	_, err := helper.DBService.Exec(sqlStatement, t.Content, util.GetNowTime(), t.ID, userID)
	if err != nil {
		panic("update memo error")
	}

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "memo update succeed",
		Succeed:       true,
		Data:          &t,
	})
}

// DeleteMemo post
func DeleteMemo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, ok := helper.GetUserIDFromSession(r)
	if !ok {
		panic("You have not sign in")
	}

	t := struct {
		ID string
	}{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&t); err != nil {
		panic("request data type error")
	}

	sqlStatement := `
		DELETE FROM memos
		WHERE id = $1 AND user_id = $2
	`

	_, err := helper.DBService.Exec(sqlStatement, t.ID, userID)
	if err != nil {
		panic("delete memo error")
	}

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "memo delete succeed",
		Succeed:       true,
		Data:          &t,
	})
}
