package response

import (
	"go-api/core"
	"go-api/helper"
	"net/http"
)

// This send the data to the user, getting that with the core path
// this is like the controller in MVC

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	res, err := core.GetAllUsers()

	if err != nil {
		panic(err)
	}

	/* users, exist := res[1]
	if !exist {
		panic("This user not exists")
	} */

	// setting the header to the response
	/* w.Header().Set("Content-Type", "application/json")

	// setting the response
	json.NewEncoder(w).Encode(res) */

	helper.Response(res, w)

}
