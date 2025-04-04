package models

import "errors"

type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

// mock list of users
var userLastID = 0
var users = []User{}

// mock inserting into database
func CreateUser(user *User) error {
	userLastID++
	user.ID = userLastID
	users = append(users, *user)
	return nil
}

// mock selecting from database
func FindUserByUsername(username string) (*User, error) {
	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func FindUserByID(userID int) (*User, error) {
	for _, user := range users {
		if user.ID == userID {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func AllUsers() []User {
	if users == nil {
		return []User{}
	}
	return users
}
