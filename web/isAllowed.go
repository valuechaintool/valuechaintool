package web

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valuechaintool/valuechaintool/models"
)

func isAllowedMany(c *gin.Context, resourceIDs []uuid.UUID, desiredCapability string) bool {
	for _, resourceID := range resourceIDs {
		if isAllowed(c, resourceID, desiredCapability) {
			return true
		}
	}
	return false
}

func isAllowed(c *gin.Context, resourceID uuid.UUID, desiredCapability string) bool {
	capabilities, ok := c.Get("capabilities")
	if !ok {
		return false
	}
	if resourceID != models.WildCardResource {
		for _, c := range capabilities.(map[uuid.UUID][]string)[models.WildCardResource] {
			if c == desiredCapability {
				return true
			}
		}
	}
	for _, c := range capabilities.(map[uuid.UUID][]string)[resourceID] {
		if c == desiredCapability {
			return true
		}
	}
	return false
}
