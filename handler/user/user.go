package user

import (
	"encoding/json"
	"fmt"
	"neosmemo/backend/handler"
	"neosmemo/backend/helper"
	"neosmemo/backend/model"
	"neosmemo/backend/util"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// GetAllUser just for test
func GetAllUser(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	users := []model.User{}
	user := model.User{}

	rows, err := helper.DBService.Query("SELECT * FROM users")

	for rows.Next() {
		rows.Scan(util.IterStructFieldAddr(&user)...)
		users = append(users, user)
	}

	if err != nil {
		// no rows in result set
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "get all users succeed",
		Succeed:       true,
		Data:          &users,
	})
}

// GetMyUserInfo check for user login status
func GetMyUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, ok := helper.GetUserIDFromSession(r)
	if !ok {
		panic("You have not sign in")
	}

	user := model.User{}
	row := helper.DBService.QueryRow("SELECT * FROM users WHERE id = $1", userID)
	err := row.Scan(util.IterStructFieldAddr(&user)...)

	if err != nil {
		panic("no result set")
	}

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "get userinfo succeed",
		Succeed:       true,
		Data:          &user,
	})
}

// UpdateInfo TODO: just for test
func UpdateInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, ok := helper.GetUserIDFromSession(r)
	if !ok {
		panic("You have not sign in")
	}

	user := model.User{}

	row := helper.DBService.QueryRow("SELECT * FROM users WHERE id = $1", userID)
	err := row.Scan(util.IterStructFieldAddr(&user)...)

	if err != nil {
		// no rows in result set
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(&user)
}

// CheckUsernameUsed check Username was Used
func CheckUsernameUsed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := struct {
		Username string
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)

	if err != nil {
		// request data type error
		fmt.Println(err.Error())
	}

	usenameUsable := checkUsernameUsable(t.Username)

	if !usenameUsable {
		fmt.Println(t.Username, "不可用")
	} else {
		fmt.Println(t.Username, "可用")
	}

	data := struct {
		usable bool
	}{
		usable: usenameUsable,
	}

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "succeed",
		Succeed:       true,
		Data:          &data,
	})
}

// DoSignUp post
func DoSignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := struct {
		Username string
		Password string
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)

	if err != nil {
		fmt.Println("request data type error", err.Error())
		panic("request data type error")
	}

	isUsable := checkUsernameUsable(t.Username)

	if !isUsable {
		panic(t.Username + " is unusable")
	}

	user := model.User{
		ID:        util.GenUUID(),
		Username:  t.Username,
		Password:  t.Password,
		Email:     "",
		CreatedAt: util.GetNowTime(),
		UpdatedAt: util.GetNowTime(),
	}

	sqlStatement := `
		INSERT INTO users (id, username, password, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = helper.DBService.Exec(sqlStatement, util.IterStructFieldValue(&user)...)

	if err != nil {
		fmt.Println(err.Error())
		panic("Sign up failed, redo later plz")
	}

	var sessionID string = util.GenUUID()
	helper.SessionManager[sessionID] = helper.Session{
		UserID:    user.ID,
		SessionID: sessionID,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(365 * 24 * time.Hour),
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "sign up succeed",
		Succeed:       true,
		Data:          &user,
	})
}

// DoSignIn post
func DoSignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := struct {
		Username string
		Password string
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)

	if err != nil {
		fmt.Println("request data type err", err.Error())
		panic("request data type error")
	}

	user := model.User{}
	row := helper.DBService.QueryRow("SELECT * FROM users WHERE username = $1 AND password = $2", t.Username, t.Password)
	err = row.Scan(util.IterStructFieldAddr(&user)...)

	if err != nil || user.ID == "" {
		fmt.Println("Sign in failed", err.Error())
		panic("Sign in failed, plz check your password")
	}

	var sessionID string = util.GenUUID()
	helper.SessionManager[sessionID] = helper.Session{
		UserID:    user.ID,
		SessionID: sessionID,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(365 * 24 * time.Hour),
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "sign in succeed",
		Succeed:       true,
		Data:          &user,
	})
}

// DoSignOut post
func DoSignOut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
	})

	json.NewEncoder(w).Encode(&handler.Response{
		StatusCode:    http.StatusOK,
		StatusMessage: "sign out succeed",
		Succeed:       true,
	})
}

func checkUsernameUsable(u string) bool {
	count := 1
	err := helper.DBService.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", u).Scan(&count)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if count != 0 {
		return false
	}

	return true
}
