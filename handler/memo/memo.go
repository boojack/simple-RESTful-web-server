package memo

import (
	"encoding/json"
	"fmt"
	"neosmemo/backend/dbservice"
	"neosmemo/backend/model"
	"neosmemo/backend/util"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetAllMemos get
func GetAllMemos(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userID, err := req.Cookie("user_id")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Fprint(w, userID.Value)
	}

	memos := []model.Memo{}
	memo := model.Memo{}

	rows, err := dbservice.DB.Query("SELECT * FROM memos WHERE user_id = $1", userID)

	if err != nil {
		// no rows in result set
		fmt.Println(err.Error())
	} else {
		for rows.Next() {
			rows.Scan(&memo.ID, &memo.UserID, &memo.Content, &memo.CreatedAt, memo.UpdatedAt)
			memos = append(memos, memo)
		}
	}

	json.NewEncoder(w).Encode(&memos)
}

// GetMemoByID just for test
func GetMemoByID(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	id := ps.ByName("id")
	memo := model.Memo{}

	err := dbservice.DB.QueryRow("SELECT * FROM memos WHERE id = $1", id).Scan(&memo.ID, &memo.UserID, &memo.Content, &memo.CreatedAt, memo.UpdatedAt)
	if err != nil {
		// no rows in result set
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(&memo)
}

// CreateMemo post
func CreateMemo(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userID, err := req.Cookie("user_id")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Fprint(w, userID.Value)
	}

	t := struct {
		Content string
	}{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&t)

	if err != nil {
		// request data type error
		fmt.Println(err.Error())
	}

	memo := model.Memo{
		ID:        util.GenUUID(),
		UserID:    userID.Value,
		Content:   t.Content,
		CreatedAt: util.GetNowTime(),
		UpdatedAt: util.GetNowTime(),
	}

	sqlStatement := `
		INSERT INTO memos (id, user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = dbservice.DB.Exec(sqlStatement, memo.ID, memo.UserID, memo.Content, memo.CreatedAt, memo.UpdatedAt)

	if err != nil {
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(&memo)
}

// UpdateMemo post
func UpdateMemo(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userID, err := req.Cookie("user_id")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Fprint(w, userID.Value)
	}

	t := struct {
		ID      string
		Content string
	}{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&t)

	if err != nil {
		// request data type error
		fmt.Println(err.Error())
	}

	sqlStatement := `
		UPDATE memos
		SET content = $1, updated_at = $2
		WHERE id = $3
	`
	_, err = dbservice.DB.Exec(sqlStatement, t.Content, util.GetNowTime(), t.ID)

	if err != nil {
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(&t)
}

// DeleteMemo post
func DeleteMemo(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userID, err := req.Cookie("user_id")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Fprint(w, userID.Value)
	}

	t := struct {
		ID string
	}{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&t)

	if err != nil {
		// request data type error
		fmt.Println(err.Error())
	}

	sqlStatement := `
		DELETE FROM memos
		WHERE id = $1 AND user_id = $2
	`
	_, err = dbservice.DB.Exec(sqlStatement, t.ID, userID)

	if err != nil {
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(&t)
}
