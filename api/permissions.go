package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valuechaintool/valuechaintool/models"
	"github.com/valuechaintool/valuechaintool/utils"
)

func PermissionsCreate(c *gin.Context) {
	if allowed := utils.IsAllowed(c, models.WildCardResource, "managePermission"); !allowed {
		c.JSON(utils.Error(http.StatusUnauthorized, "unauthorized"))
		return
	}

	id, err := uuid.Parse(c.Param("user"))
	if err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, fmt.Sprintf("invalid user value: %s", err)))
		return
	}
	var fields map[string]interface{}
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, fmt.Sprintf("invalid request: %s", err)))
		return
	}

	roleID, err := uuid.Parse(fields["role_id"].(string))
	if err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, fmt.Sprintf("invalid role value: %s", err)))
		return
	}
	resourceID, err := uuid.Parse(fields["resource_id"].(string))
	if err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, fmt.Sprintf("invalid resource value: %s", err)))
		return
	}

	p := models.Permission{
		UserID:     id,
		ResourceID: resourceID,
		RoleID:     roleID,
	}

	if err := models.NewPermission(&p); err != nil {
		c.JSON(utils.Error(http.StatusInternalServerError, fmt.Sprintf("error while creating the item: %s", err)))
		return
	}
	c.JSON(http.StatusOK, p)
}

func PermissionsList(c *gin.Context) {
	id, err := uuid.Parse(c.Param("user"))
	if err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, fmt.Sprintf("invalid user value: %s", err)))
		return
	}
	if allowed := utils.IsAllowed(c, id, "readUser"); !allowed {
		c.JSON(utils.Error(http.StatusUnauthorized, "unauthorized"))
		return
	}

	items, err := models.ListPermissionsByUser(id)
	if err != nil {
		c.JSON(utils.Error(http.StatusInternalServerError, fmt.Sprintf("error while fetching the items: %s", err)))
		return
	}
	c.JSON(http.StatusOK, items)
}

func PermissionsDelete(c *gin.Context) {
	if allowed := utils.IsAllowed(c, models.WildCardResource, "managePermission"); !allowed {
		c.JSON(utils.Error(http.StatusUnauthorized, "unauthorized"))
		return
	}

	id, err := uuid.Parse(c.Param("user"))
	if err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, fmt.Sprintf("invalid user value: %s", err)))
		return
	}
	permissionID, err := uuid.Parse(c.Param("permission"))
	if err != nil {
		c.JSON(utils.Error(http.StatusUnprocessableEntity, fmt.Sprintf("invalid permission value: %s", err)))
		return
	}

	items, err := models.ListPermissionsByUser(id)
	if err != nil {
		c.JSON(utils.Error(http.StatusInternalServerError, fmt.Sprintf("error while fetching the item: %s", err)))
		return
	}
	var item models.Permission
	for _, i := range items {
		if i.ID == permissionID {
			item = i
		}
	}
	if item.ID == uuid.Nil {
		c.JSON(utils.Error(http.StatusUnauthorized, "unauthorized"))
		return
	}

	if err := item.Delete(); err != nil {
		c.JSON(utils.Error(http.StatusInternalServerError, fmt.Sprintf("error while deleting the item: %s", err)))
		return
	}
	c.JSON(utils.Error(http.StatusGone, "item successfully deleted"))
}
