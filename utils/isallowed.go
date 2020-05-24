package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valuechaintool/valuechaintool/models"
)

func IsAllowed(c *gin.Context, resourceID uuid.UUID, desiredCapability string) bool {
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
