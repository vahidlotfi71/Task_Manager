package Rules

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// InEnum returns a ValidationRule that checks if a field value matches allowed enum values.
func InEnum(allowed []string) ValidationRule {
	return func(c *gin.Context, field_name string) (passed bool, message string, err error) {
		var jsonBody map[string]interface{}

		// Parse JSON body into a map for dynamic field access
		if err := c.ShouldBindJSON(&jsonBody); err != nil {
			return false, "", fmt.Errorf("failed to parse JSON body: %w", err)
		}

		value := jsonBody[field_name]
		// Convert to string, trim whitespace and check against allowed values
		str := strings.TrimSpace(fmt.Sprintf("%v", value))
		for _, item := range allowed {
			if str == item {
				return true, "", nil
			}
		}

		return false, fmt.Sprintf("The %s field must be one of: %s", field_name, strings.Join(allowed, ", ")), nil
	}
}
