package Rules

import "github.com/gin-gonic/gin"

// Optional returns a ValidationRule that always passes
func Optional() ValidationRule {
	return func(c *gin.Context, field_name string) (bool, string, error) {
		return true, "", nil
	}
}
