package services

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/models"
)

func GetUserProfile(userID int) (*models.User, error) {
	user, err := models.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func GetAllUsers() []models.User {
	users := models.AllUsers()

	// for i := range users {
	// 	users[i].Password = ""
	// }

	return users
}
