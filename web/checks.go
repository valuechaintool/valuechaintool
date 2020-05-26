package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Healthz(c *gin.Context) {
	c.HTML(http.StatusOK, "ok", nil)
}

func Readiness(c *gin.Context) {
	c.HTML(http.StatusOK, "ok", nil)
}
