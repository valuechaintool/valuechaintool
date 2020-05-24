package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valuechaintool/valuechaintool/models"
)

func User(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, "invalid userID key")
		return
	}
	user, err := models.GetUser(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("problem while retrive the user: %s", err))
		return
	}
	c.JSON(http.StatusOK, user)
}
