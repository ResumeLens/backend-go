package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/resumelens/authservice/internal/db"
	"github.com/resumelens/authservice/internal/models"
)

func OrganizationAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jobID := c.Param("id")
		if jobID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
			c.Abort()
			return
		}

		if !isValidUUID(jobID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format"})
			c.Abort()
			return
		}

		userOrgID, exists := c.Get("organizationID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User organization not found"})
			c.Abort()
			return
		}

		if orgIDStr, ok := userOrgID.(string); !ok || !isValidUUID(orgIDStr) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user organization ID"})
			c.Abort()
			return
		}

		var job models.Job
		if err := db.DB.Where("id = ? AND organization_id = ? AND is_active = ?", jobID, userOrgID, true).First(&job).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Job not found or does not belong to your organization"})
			c.Abort()
			return
		}

		c.Set("job", job)
		c.Next()
	}
}

func isValidUUID(id string) bool {
	if len(id) != 36 {
		return false
	}
	parts := strings.Split(id, "-")
	return len(parts) == 5 && len(parts[0]) == 8 && len(parts[1]) == 4 && len(parts[2]) == 4 && len(parts[3]) == 4 && len(parts[4]) == 12
}
