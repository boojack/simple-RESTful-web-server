package memo

import (
	"encoding/json"
	"fmt"
	"neosmemo/backend/handler"
	"neosmemo/backend/helper"
	"neosmemo/backend/model"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetAllMemos get
func GetAllMemos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, ok := helper.GetUserIDFromSession(r)
	if !ok {
		panic("You have not sign in")
	}

	memos, err := model.GetMemosByUserID(userID)

	if err != nil {
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
	memo, err := model.GetMemoByID(id)

	if err != nil {
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

	memo, err := model.CreateMemo(userID, t.Content)

	if err != nil {
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

	if err := model.UpdateMemoByID(t.ID, t.Content, userID); err != nil {
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

	if err := model.DeleteMemoByID(t.ID, userID); err != nil {
		panic("delete memo error")
	}

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "memo delete succeed",
		Succeed:       true,
		Data:          &t,
	})
}
