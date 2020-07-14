package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/valuechaintool/valuechaintool/models"
)

func MiddlewareHSTS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%v", viper.GetInt("tls.hsts")))
		c.Next()
	}
}

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
