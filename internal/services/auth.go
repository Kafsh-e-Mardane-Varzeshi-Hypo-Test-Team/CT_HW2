package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"

	"fmt"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	jwt "github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/pkg"
)

func (s *Service) RegisterUser(c context.Context, username, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user, err := s.Queries.CreateUser(c, generated.CreateUserParams{
		Username:          username,
		EncryptedPassword: string(hashedPassword),
		Role:              generated.UserRoleNormal,
	})
	if err != nil {
		return 0, err
	}

	return int(user.ID), nil
}

func (s *Service) AuthenticateUser(c context.Context, username, password string) (string, error) {
	user, err := s.Queries.GetUserByUsername(c, username)
	if err == nil {
		passwordMismatch := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password))
		if passwordMismatch == nil {
			return jwt.GenerateToken(fmt.Sprint(user.ID), s.Configs.JWT.SecretKey)
		}
	}
	return "", errors.New("invalid credentials")
}
