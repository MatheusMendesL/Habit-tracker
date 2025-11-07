package core

import (
	"fmt"
	"go-api/db"
)

// example file
// On this path will have the structures
// Here is the rules about the handles, like verifications to send to the db and to de handler for the response
// the funcs will return a data what the handle will use to return the JSON to the user

func GetAllUsers() (map[int]db.User, error) {
	res, err := db.GetUserDetails()

	fmt.Println(res)

	if err != nil {
		fmt.Println(err)
	}

	return res, nil
}
