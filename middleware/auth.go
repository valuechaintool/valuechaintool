package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valuechaintool/valuechaintool/models"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if len(authorization) < 8 {
			c.JSON(http.StatusUnauthorized, "unrecognized authorization method")
			c.Abort()
			return
		}
		if authorization[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, "unrecognized authorization method")
			c.Abort()
			return
		}
		session, err := models.GetSession(authorization[7:])
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("userID", session.UserID)
		capabilities, err := models.ListCapabilitiesByUser(session.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("capabilities", capabilities)
		c.Next()
	}
}
