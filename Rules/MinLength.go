package Rules

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// the minimum required length for the field value
func MinLength(min uint) ValidationRule {
	return func(c *gin.Context, field_name string) (bool, string, error) {
		// Read cached JSON body from context
		body, _ := c.Get("json_body")
		jsonBody := body.(map[string]interface{})

		// Get field value from parsed body
		value := jsonBody[field_name]

		// Convert value to string and check minimum length
		str := fmt.Sprintf("%v", value)
		if len(str) < int(min) {
			return false, fmt.Sprintf("The %s field must be at least %d characters long", field_name, min), nil
		}
		return true, "", nil
	}
}
