package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valuechaintool/valuechaintool/models"
)

func MiddlewareAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_id")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		session, err := models.GetSession(token)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		c.Set("userID", session.UserID)
		capabilities, err := models.ListCapabilitiesByUser(session.UserID)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		c.Set("capabilities", capabilities)
		c.Next()
	}
}
