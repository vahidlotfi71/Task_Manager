// Package Rules provides custom validation utilities for Gin framework handlers.
// It enables field-based validation with custom rules and standardized error messages.
package Rules

import "github.com/gin-gonic/gin"

// ValidationRule defines a function type for validating a specific field in a Gin context.
type ValidationRule func(c *gin.Context, field_name string) (passed bool, message string, err error)

// FieldRules groups validation rules for a specific field.
type FieldRules struct {
	FieldName string
	Rules     []ValidationRule
}
