package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valuechaintool/valuechaintool/models"
	"github.com/valuechaintool/valuechaintool/utils"
)

func Register(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, "invalid JSON provided"))
		return
	}
	if err := models.NewUser(&u); err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, err.Error()))
		return
	}
	c.JSON(http.StatusOK, u)
}
