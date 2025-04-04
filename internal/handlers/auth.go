package handlers

import (
	"fmt"
	"net/http"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/models"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/services"
	jwt "github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/pkg"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := services.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register"})
		return
	}

	token, err := jwt.GenerateToken(fmt.Sprint(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	// cookie yum yum
	c.SetCookie("session_token", token, int(jwt.SessionMaxAge.Seconds()), "/", "", true, true)

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully and logged in"})
}

func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := services.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// cookie yum yum
	// TODO: what should domain be?
	c.SetCookie("session_token", token, int(jwt.SessionMaxAge.Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}
