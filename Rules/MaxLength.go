package Rules

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// the maximum required length for the field value
func MaxLength(min uint) ValidationRule {
	return func(c *gin.Context, field_name string) (passed bool, message string, err error) {
		var jsonBody map[string]interface{}

		// Parse JSON body into a map for dynamic field access
		if err := c.ShouldBindJSON(&jsonBody); err != nil {
			return false, "", fmt.Errorf("failed to parse JSON body: %w", err)
		}

		value := jsonBody[field_name]
		// Convert value to string and check minimum length
		str := fmt.Sprintf("%v", value)
		if len(str) > int(min) {
			return false, fmt.Sprintf("The %s field must be at least %d characters long", field_name, min), nil
		}

		return true, "", nil
	}
}
