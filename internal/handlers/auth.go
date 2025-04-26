package handlers

import (
	"fmt"
	"net/http"

	jwt "github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/pkg"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SignupGet(c *gin.Context) {
	_, exists := c.Get("User")

	if !exists {
		c.HTML(http.StatusOK, "signup.html", nil)
	} else {
		// TODO: code must be 302 for get, and 303 for post
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
}

func (h *Handler) LoginGet(c *gin.Context) {
	_, exists := c.Get("User")

	if !exists {
		c.HTML(http.StatusOK, "login.html", nil)
	} else {
		// TODO: code must be 302 for get, and 303 for post
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
}

func (h *Handler) SignupPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	confirm_password := c.PostForm("confirm_password")

	if username == "" || password == "" || confirm_password == "" || password != confirm_password {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"Error": "Invalid username, password or confirm password",
		})
		c.Abort()
		return
	}

	userId, err := h.Service.RegisterUser(c.Request.Context(), username, password)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			// TODO: maybe remove original error message after development
			"Error": "Could not register user, " + err.Error(),
		})
		c.Abort()
		return
	}

	token, err := jwt.GenerateToken(fmt.Sprint(userId), h.Service.Configs.JWT.SecretKey)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			"Error": "Could not generate token, " + err.Error(),
		})
		c.Abort()
		return
	}
	// cookie yum yum
	c.SetCookie("session_token", token, int(jwt.SessionMaxAge.Seconds()), "/", "", true, true)

	c.Redirect(http.StatusCreated, "/")
}

func (h *Handler) LoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	token, err := h.Service.AuthenticateUser(c.Request.Context(), username, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error": "Invalid username or password",
		})
		c.Abort()
		return
	}

	// cookie yum yum
	// TODO: what should domain be?
	c.SetCookie("session_token", token, int(jwt.SessionMaxAge.Seconds()), "/", "", true, true)

	c.Redirect(http.StatusFound, "/")
}

func (h *Handler) Logout(c *gin.Context) {
	_, exists := c.Get("User")

	if !exists {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error": "You are not logged in",
		})
		c.Abort()
		return
	}

	ClearSessionCookie(c)

	c.Redirect(http.StatusFound, "/")
}

func ClearSessionCookie(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "", true, true)
}
