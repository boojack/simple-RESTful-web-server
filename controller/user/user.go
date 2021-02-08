package user

import (
	"encoding/json"
	"fmt"
	"neosmemo/backend/dbservice"
	"neosmemo/backend/model"
	"neosmemo/backend/util"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetAllUser just for test
func GetAllUser(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	users := []model.User{}
	var user model.User

	rows, err := dbservice.DB.Query("SELECT * FROM users")
	for rows.Next() {
		rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
		fmt.Println(&user)
		users = append(users, user)
	}
	if err != nil {
		// no rows in result set
		func() {
			fmt.Println(err.Error())
		}()
	}

	fmt.Fprintf(w, "hello, %#+v!\n", &users)
}

// GetUserByID just for test
func GetUserByID(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id, _ := util.ParseInt(ps.ByName("id"))
	var user model.User

	err := dbservice.DB.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
	if err != nil {
		// no rows in result set
		func() {
			fmt.Println(err.Error())
		}()
	}

	fmt.Fprintf(w, "hello, %#+v!\n", &user)
}

// CheckUsernameUsed check Username was Used
func CheckUsernameUsed(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var t struct {
		Username string
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)

	if err != nil {
		// request data type error
		func() {
			fmt.Println(err.Error())
		}()
	}

	var count int
	err = dbservice.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", t.Username).Scan(&count)

	if count != 0 {
		fmt.Println(t.Username, "不可用")
	} else {
		fmt.Println(t.Username, "可用")
	}
}

// SignUp post
func SignUp(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var t struct {
		Username string
		Password string
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)

	if err != nil {
		// request data type error
		func() {
			fmt.Println(err.Error())
		}()
	}

	var count int
	err = dbservice.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", t.Username).Scan(&count)

	if count != 0 {
		fmt.Println(t.Username, "不可用")
		return
	}
	fmt.Println(t.Username, "可用")

	user := model.User{
		ID:        util.GenUUID(),
		Username:  t.Username,
		Password:  t.Password,
		Email:     "",
		CreatedAt: util.GetNowTime(),
	}

	sqlStatement := `
		INSERT INTO users (id, username, password, email, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = dbservice.DB.Exec(sqlStatement, user.ID, user.Username, user.Password, user.Email, user.CreatedAt)

	if err != nil {
		panic(err)
	}

	fmt.Println(user)
}
