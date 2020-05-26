package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valuechaintool/valuechaintool/models"
)

func Logout(c *gin.Context) {
	token, err := c.Cookie("session_id")
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, errors.New("invalid login credentials"))
		return
	}
	session, err := models.GetSession(token)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, errors.New("invalid login credentials"))
		return
	}
	if err := session.Delete(); err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, errors.New("invalid login credentials"))
		return
	}
	c.SetCookie("session_id", "", 0, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}
