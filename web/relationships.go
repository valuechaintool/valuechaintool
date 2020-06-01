package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valuechaintool/valuechaintool/models"
)

// RelationshipsCreatePost parses a form to create a new relationship
func RelationshipsCreatePost(c *gin.Context) {
	leftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	rightID, err := uuid.Parse(c.PostForm("right"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if !isAllowedMany(c, []uuid.UUID{leftID, rightID}, "createRelationship") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
	tier, err := strconv.Atoi(c.PostForm("tier"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	relationship := models.Relationship{
		LeftID:   leftID,
		RightID:  rightID,
		LeftTier: tier,
	}
	if err := models.NewRelationship(&relationship); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/companies/%s?edit_mode=true", leftID.String()))
}

// RelationshipsUpdate allows to change relationships
func RelationshipsUpdate(c *gin.Context) {
	companyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	relationshipID, err := uuid.Parse(c.Param("rid"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	relationship, err := models.GetRelationship(relationshipID)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if !isAllowedMany(c, []uuid.UUID{relationship.LeftID, relationship.RightID}, "updateRelationship") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}

	if relationship.LeftID != companyID {
		*relationship = relationship.Reverse()
	}

	tier, err := strconv.Atoi(c.PostForm("tier"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	if err := relationship.Update(tier, c.PostForm("notes")); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/companies/%s?edit_mode=true", companyID.String()))
}

// RelationshipsDelete responds to /relationships/[ID]/delete url
func RelationshipsDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("rid"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	relationship, err := models.GetRelationship(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !isAllowedMany(c, []uuid.UUID{relationship.LeftID, relationship.RightID}, "deleteRelationship") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
	if err := relationship.Delete(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/companies/%s?edit_mode=true", c.Param("id")))
}
