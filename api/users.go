package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valuechaintool/valuechaintool/models"
	"github.com/valuechaintool/valuechaintool/utils"
)

func UsersList(c *gin.Context) {
	if allowed := utils.IsAllowed(c, models.WildCardResource, "readUser"); !allowed {
		c.JSON(utils.Error(http.StatusUnauthorized, "unauthorized"))
		return
	}

	users, err := models.ListUsers(nil)
	if err != nil {
		c.JSON(utils.Error(http.StatusInternalServerError, fmt.Sprintf("error while fetching the items: %s", err)))
		return
	}
	c.JSON(http.StatusOK, users)
}

func UsersRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("user"))
	if err != nil {
		c.JSON(utils.Error(http.StatusUnauthorized, "unauthorized"))
		return
	}
	if allowed := utils.IsAllowed(c, id, "readUser"); !allowed {
		c.JSON(utils.Error(http.StatusUnauthorized, "unauthorized"))
		return
	}

	user, err := models.GetUser(id)
	if err != nil {
		c.JSON(utils.Error(http.StatusInternalServerError, fmt.Sprintf("error while fetching the item: %s", err)))
		return
	}
	c.JSON(http.StatusOK, user)
}

func UsersUpdate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("user"))
	if err != nil {
		c.JSON(utils.Error(http.StatusUnauthorized, "unauthorized"))
		return
	}
	if allowed := utils.IsAllowed(c, id, "updateUser"); !allowed {
		c.JSON(utils.Error(http.StatusUnauthorized, "unauthorized"))
		return
	}

	user, err := models.GetUser(id)
	if err != nil {
		c.JSON(utils.Error(http.StatusNotFound, fmt.Sprintf("problem while retrive the users: %s", err)))
		return
	}

	var fields map[string]interface{}
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, fmt.Sprintf("invalid request: %s", err)))
		return
	}
	if err := utils.MapOnlyContains(fields, []string{"username", "real_name", "email", "password"}); err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, fmt.Sprintf("validation failed: %s", err)))
		return
	}
	if err := user.Update(fields); err != nil {
		c.JSON(utils.Error(http.StatusInternalServerError, fmt.Sprintf("error while saving the changes: %s", err)))
		return
	}
	c.JSON(http.StatusOK, user)
}

func UsersDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("user"))
	if err != nil {
		c.JSON(utils.Error(http.StatusUnauthorized, "unauthorized"))
		return
	}
	if allowed := utils.IsAllowed(c, id, "deleteUser"); !allowed {
		c.JSON(utils.Error(http.StatusUnauthorized, "unauthorized"))
		return
	}

	user, err := models.GetUser(id)
	if err != nil {
		c.JSON(utils.Error(http.StatusNotFound, fmt.Sprintf("problem while retrive the users: %s", err)))
		return
	}
	if err := user.Delete(); err != nil {
		c.JSON(utils.Error(http.StatusInternalServerError, fmt.Sprintf("error while deleting the item: %s", err)))
		return
	}
	c.JSON(utils.Error(http.StatusGone, "item successfully deleted"))
}
