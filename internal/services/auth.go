package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"fmt"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/models"
	jwt "github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/pkg"
)

func RegisterUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return models.CreateUser(user)
}

func AuthenticateUser(username, password string) (string, error) {
	user, err := models.FindUserByUsername(username)
	if err == nil {
		passwordMismatch := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if passwordMismatch == nil {
			return jwt.GenerateToken(fmt.Sprint(user.ID))
		}
	}
	return "", errors.New("invalid credentials")
}
