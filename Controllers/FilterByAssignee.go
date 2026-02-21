package Controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Config"
	"github.com/vahidlotfi71/Task_Manager/Models/Repository"
	"github.com/vahidlotfi71/Task_Manager/Resources"
)

func FilterByAssignee(c *gin.Context) {
	assignee := strings.TrimSpace(c.Param("assignee"))

	if assignee == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assignee name is required"})
		return
	}

	tasks, meta, err := Repository.Paginate(
		Config.DB.Where("deleted_at IS NULL AND assignee = ?", assignee).Order("id"),
		c,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     Resources.Collection(tasks),
		"metadata": meta,
	})
}
