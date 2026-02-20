package Rules

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// This function parses the entire JSON body. For multiple validations
func Required(c *gin.Context, field_name string) (passed bool, message string, err error) {
	var jsonBody map[string]interface{}
	// Parse JSON body into a map for dynamic field access
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		return false, "", fmt.Errorf("failed to parse JSON body: %w", err)
	}

	// Check field existence and non-empty value
	value, exists := jsonBody[field_name]
	if !exists || value == nil || value == "" {
		return false, fmt.Sprintf("The %s field is required", field_name), nil
	}

	return true, "", nil
}
