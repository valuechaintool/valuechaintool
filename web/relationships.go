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
	}
	rightID, err := uuid.Parse(c.PostForm("right"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	tier, err := strconv.Atoi(c.PostForm("tier"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	relationship := models.Relationship{
		LeftID:  leftID,
		RightID: rightID,
		Tier:    tier,
	}
	if err := models.NewRelationship(&relationship); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/companies/%s", leftID.String()))
}

// RelationshipsDelete responds to /relationships/[ID]/delete url
func RelationshipsDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("rid"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	relationship, err := models.GetRelationship(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	if err := relationship.Delete(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/companies/%s", c.Param("id")))
}
