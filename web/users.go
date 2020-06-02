package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valuechaintool/valuechaintool/models"
)

// UsersList renders the /users page
func UsersList(c *gin.Context) {
	if !isAllowed(c, models.WildCardResource, "readUser") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
	users, err := models.ListUsers(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	d := struct {
		PageTitle string
		Users     []models.User
	}{
		PageTitle: "User list - ValueChain",
		Users:     users,
	}
	c.HTML(http.StatusOK, "users-list.html", d)
}

// UsersRead renders the /users/[ID] page
func UsersRead(c *gin.Context) {
	if !isAllowed(c, models.WildCardResource, "readUser") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
	id, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	user, err := models.GetUser(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// Roles
	roles, err := models.ListRoles(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rls := make(map[string]models.Role)
	for _, role := range roles {
		rls[role.ID.String()] = role
	}

	// Permissions
	permissions, err := models.ListPermissionsByUser(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	permissionsGlobal := []models.Permission{}
	permissionsLocal := []models.Permission{}
	for _, permission := range permissions {
		if permission.ResourceID == models.WildCardResource {
			permissionsGlobal = append(permissionsGlobal, permission)
		} else {
			permissionsLocal = append(permissionsLocal, permission)
		}
	}

	// Companies
	companies, err := models.ListCompanies(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	cps := make(map[string]models.Company)
	for _, company := range companies {
		cps[company.ID.String()] = company
	}

	d := struct {
		PageTitle            string
		User                 models.User
		PermissionsGlobal    []models.Permission
		PermissionsLocal     []models.Permission
		Roles                map[string]models.Role
		Companies            map[string]models.Company
		CanUpdateUser        bool
		CanDeleteUser        bool
		CanManagePermissions bool
	}{
		PageTitle:            fmt.Sprintf("User %s information", (*user).Username),
		User:                 *user,
		PermissionsGlobal:    permissionsGlobal,
		PermissionsLocal:     permissionsLocal,
		Roles:                rls,
		Companies:            cps,
		CanUpdateUser:        isAllowed(c, models.WildCardResource, "updateUser"),
		CanDeleteUser:        isAllowed(c, models.WildCardResource, "deleteUser"),
		CanManagePermissions: isAllowed(c, models.WildCardResource, "managePermission"),
	}
	c.HTML(http.StatusOK, "users-single.html", d)
}

// UsersDelete responds to /users/[ID]/delete url
func UsersDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if !isAllowed(c, models.WildCardResource, "deleteUser") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
	user, err := models.GetUser(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err := user.Delete(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, "/users")
}
