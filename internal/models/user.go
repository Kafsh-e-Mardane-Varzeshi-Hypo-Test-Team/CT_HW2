package models

type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

// mock list of users
var users = []User{
	{ID: 1, Username: "user1", Password: "hashpassword1"},
	{ID: 2, Username: "user2", Password: "hashpassword2"},
}

// mock inserting into database
func CreateUser(user *User) error {
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
	return nil, nil
}
