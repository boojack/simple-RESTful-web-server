package user

import (
	"encoding/json"
	"fmt"
	"neosmemo/backend/dbservice"
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

	rows, err := dbservice.DB.Query("SELECT * FROM users")

	for rows.Next() {
		// TODO: 顺序必须和数据库字段定义的顺序一致。如何抽象？
		rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, user.UpdatedAt)
		users = append(users, user)
	}

	if err != nil {
		// no rows in result set
		fmt.Println(err.Error())
	}

	json.NewEncoder(w).Encode(&users)
}

// GetUserByID just for test
func GetUserByID(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	user := model.User{}

	err := dbservice.DB.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, user.UpdatedAt)

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

	err := dbservice.DB.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, user.UpdatedAt)

	if err != nil {
		// no rows in result set
		fmt.Println(err.Error())
	}

	fmt.Fprintf(w, "hello, %#v!\n", &user)
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
		// request data type error
		fmt.Println(err.Error())
	}

	usenameUsable := checkUsernameUsable(t.Username)

	if !usenameUsable {
		fmt.Println(t.Username, "不可用")
		// TODO: 错误处理
		return
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
		VALUES ($1, $2, $3, $4, $5)`
	_, err = dbservice.DB.Exec(sqlStatement, user.ID, user.Username, user.Password, user.Email, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		fmt.Println(err.Error())
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "userid", Value: user.ID, Expires: expiration}
	http.SetCookie(w, &cookie)
}

// DoSignIn post
func DoSignIn(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	t := struct {
		Username string
		Password string
	}{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)

	if err != nil {
		// request data type error
		fmt.Println(err.Error())
	}

	userID := ""
	err = dbservice.DB.QueryRow("SELECT id FROM users WHERE username = $1 AND password = $2", t.Username, t.Password).Scan(&userID)

	if err != nil || userID == "" {
		// fmt.Println(err.Error())
		fmt.Println("Sign in failed")
		return
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "userid", Value: userID, Expires: expiration}
	http.SetCookie(w, &cookie)

	// TODO:
	user := model.User{
		ID:       userID,
		Username: t.Username,
	}
	json.NewEncoder(w).Encode(&user)
}

func checkUsernameUsable(u string) bool {
	count := 1
	err := dbservice.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", u).Scan(&count)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if count != 0 {
		return false
	}

	return true
}
