package db

// example file
// this will have the querys to the db using this table

type User struct {
	Id    int
	Name  string
	Email string
}

func GetAllUsers() {

}

func GetUserDetails() (map[int]User, error) {
	return map[int]User{
		1: {
			Id:    1,
			Name:  "Matheus",
			Email: "matheuslopes131@gmail.com",
		},
	}, nil
}
