package TaskValidation

import (
	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Middlewares"
	"github.com/vahidlotfi71/Task_Manager/Rules"
)

func Store() gin.HandlerFunc {
	return Middlewares.ValidationMiddleware([]Rules.FieldRules{
		{
			FieldName: "title",
			Rules:     []Rules.ValidationRule{Rules.Required, Rules.MinLength(3), Rules.MaxLength(255)},
		},
		{
			FieldName: "description",
			Rules:     []Rules.ValidationRule{Rules.MinLength(3), Rules.MaxLength(255)},
		},
		{
			FieldName: "status",
			Rules:     []Rules.ValidationRule{Rules.InEnum([]string{"done", "in_pending", "pending"})},
		},
		{
			FieldName: "assignee",
			Rules:     []Rules.ValidationRule{Rules.MinLength(3), Rules.MaxLength(255)},
		},
	})
}

func Update() gin.HandlerFunc {
	return Middlewares.ValidationMiddleware([]Rules.FieldRules{
		{
			FieldName: "title",
			Rules:     []Rules.ValidationRule{Rules.Optional(), Rules.MinLength(3), Rules.MaxLength(255)},
		},
		{
			FieldName: "description",
			Rules:     []Rules.ValidationRule{Rules.Optional(), Rules.MinLength(3), Rules.MaxLength(255)},
		},
		{
			FieldName: "status",
			Rules:     []Rules.ValidationRule{Rules.Optional(), Rules.InEnum([]string{"done", "in_pending", "pending"})},
		},
		{
			FieldName: "assignee",
			Rules:     []Rules.ValidationRule{Rules.Optional(), Rules.MinLength(3), Rules.MaxLength(255)},
		},
	})
}
