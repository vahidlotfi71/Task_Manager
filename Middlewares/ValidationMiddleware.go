package Middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Rules"
)

// ValidationMiddleware creates a Gin middleware for request validation using field-based rules.
func ValidationMiddleware(schema []Rules.FieldRules) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Parse JSON body once and cache in context for reuse by validation rules
		var jsonBody map[string]interface{}
		if err := c.ShouldBindJSON(&jsonBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid JSON body",
			})
			return
		}
		// Store parsed body in context for access by individual validation rules
		c.Set("json_body", jsonBody)

		// Execute all validation rules for each field
		for _, field_rules := range schema {
			for _, rule := range field_rules.Rules {
				passed, message, err := rule(c, field_rules.FieldName)

				// Handle system/internal errors
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
					return
				}
				// Handle validation failures
				if !passed {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"message": message,
					})
					return
				}
			}
		}

		c.Next()
	}
}
