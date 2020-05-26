package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home renders the / page
func Home(c *gin.Context) {
	c.Redirect(http.StatusFound, "/companies")
}
