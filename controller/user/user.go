package user

import (
	"fmt"
	"neosmemo/backend/dbservice"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Test just for test
func Test() {

}

// GetUserByID just for test
func GetUserByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Prepare statement for reading data
	stmtOut, err := dbservice.DB.Prepare("SELECT username FROM users WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	var username string
	// Query the user-id of 1
	err = stmtOut.QueryRow(ps.ByName("id")).Scan(&username)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "hello, %s!\n", username)
}
