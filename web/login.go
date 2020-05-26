package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valuechaintool/valuechaintool/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginPost(c *gin.Context) {
	user, err := models.GetUserByName(c.PostForm("username"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, errors.New("invalid login credentials"))
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.PostForm("password"))) != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, errors.New("invalid login credentials"))
		return
	}
	session, err := models.NewSession(user.ID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.SetCookie("session_id", session.ID, 60*60*24, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}
