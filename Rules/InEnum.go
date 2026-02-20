package Rules

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// InEnum returns a ValidationRule that checks if a field value matches allowed enum values.
func InEnum(allowed []string) ValidationRule {
	return func(c *gin.Context, field_name string) (bool, string, error) {
		body, _ := c.Get("json_body")
		jsonBody := body.(map[string]interface{})

		value := jsonBody[field_name]
		// Convert to string, trim whitespace
		str := strings.TrimSpace(fmt.Sprintf("%v", value))
		for _, item := range allowed {
			if str == item {
				return true, "", nil
			}
		}
		return false, fmt.Sprintf("The %s field must be one of: %s", field_name, strings.Join(allowed, ", ")), nil
	}
}
