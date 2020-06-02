package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valuechaintool/valuechaintool/models"
)

// PermissionsCreatePost parses a form to assign new roles to a user
func PermissionsCreatePost(c *gin.Context) {
	if !isAllowed(c, models.WildCardResource, "managePermission") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	resourceID, err := uuid.Parse(c.PostForm("resource"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	roleID, err := uuid.Parse(c.PostForm("role"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	permission := models.Permission{
		UserID:     userID,
		RoleID:     roleID,
		ResourceID: resourceID,
	}
	if err := models.NewPermission(&permission); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/users/%s?edit_mode=true", userID.String()))
}

// PermissionsDelete responds to /user/[ID]/permissions/[ID]/delete url
func PermissionsDelete(c *gin.Context) {
	if !isAllowed(c, models.WildCardResource, "managePermission") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
	permissionID, err := uuid.Parse(c.Param("permission_id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	permission, err := models.GetPermission(permissionID)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err := permission.Delete(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/users/%s?edit_mode=true", c.Param("user_id")))
}
