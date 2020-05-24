package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valuechaintool/valuechaintool/models"
	"github.com/valuechaintool/valuechaintool/utils"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	type login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var l login
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, "invalid JSON provided"))
		return
	}
	user, err := models.GetUserByName(l.Username)
	if err != nil {
		c.JSON(utils.Error(http.StatusUnauthorized, "invalid login credentials"))
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(l.Password)) != nil {
		c.JSON(utils.Error(http.StatusUnauthorized, "invalid login credentials"))
		return
	}
	session, err := models.NewSession(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, session.ID)
}
