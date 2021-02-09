package user

import (
	"encoding/json"
	"fmt"
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

	json.NewEncoder(w).Encode(&users)
}

// GetMyUserInfo check for user login status
func GetMyUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := util.GetKeyValueFromCookie("user_id", r)
	if err != nil {
		panic("You have not sign in")
	}

	user := model.User{}
	row := helper.DBService.QueryRow("SELECT * FROM users WHERE id = $1", userID)
	err = row.Scan(util.IterStructFieldAddr(&user)...)

	if err != nil {
		// no rows in result set
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(&user)
}

// GetUserInfo just for test
func GetUserInfo(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := ps.ByName("id")
	user := model.User{}

	row := helper.DBService.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(util.IterStructFieldAddr(&user)...)

	if err != nil {
		// no rows in result set
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(&user)
}

// CheckUsernameUsed check Username was Used
func CheckUsernameUsed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	json.NewEncoder(w).Encode(data)
}

// DoSignUp post
func DoSignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	usenameUsable := checkUsernameUsable(t.Username)

	if !usenameUsable {
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
		panic("Sign in failed, plz check your password")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    user.ID,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})
	json.NewEncoder(w).Encode(&user)
}

// DoSignIn post
func DoSignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    user.ID,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})
	json.NewEncoder(w).Encode(&user)
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
