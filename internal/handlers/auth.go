package handlers

import (
	"fmt"
	"net/http"

	jwt "github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/pkg"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SignupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

func (h *Handler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *Handler) SignupHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	confirm_password := c.PostForm("confirm_password")

	if username == "" || password == "" || confirm_password == "" || password != confirm_password {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"Error": "Invalid username, password or confirm password",
		})
		return
	}

	userId, err := h.Service.RegisterUser(c.Request.Context(), username, password)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			// TODO: maybe remove original error message after development
			"Error": "Could not register user, " + err.Error(),
		})
		return
	}

	token, err := jwt.GenerateToken(fmt.Sprint(userId), h.Service.Configs.JWT.SecretKey)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			"Error": "Could not generate token, " + err.Error(),
		})
		return
	}
	// cookie yum yum
	c.SetCookie("session_token", token, int(jwt.SessionMaxAge.Seconds()), "/", "", true, true)

	c.Redirect(http.StatusCreated, "/")
}

func (h *Handler) LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	token, err := h.Service.AuthenticateUser(c.Request.Context(), username, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error": "Invalid username or password",
		})
		return
	}

	// cookie yum yum
	// TODO: what should domain be?
	c.SetCookie("session_token", token, int(jwt.SessionMaxAge.Seconds()), "/", "", true, true)

	c.Redirect(http.StatusFound, "/")
}
