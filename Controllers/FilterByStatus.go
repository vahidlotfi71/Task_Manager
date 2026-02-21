package Controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Config"
	"github.com/vahidlotfi71/Task_Manager/Models"
	"github.com/vahidlotfi71/Task_Manager/Models/Repository"
	"github.com/vahidlotfi71/Task_Manager/Resources"
)

func FilterByStatus(c *gin.Context) {
	status := Models.TaskStatus(c.Param("status"))

	validStatuses := map[Models.TaskStatus]bool{
		Models.StatusPending:    true,
		Models.StatusInProgress: true,
		Models.StatusDone:       true,
	}

	if !validStatuses[status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid status. Valid values: pending, in_progress, done",
		})
		return
	}

	tasks, meta, err := Repository.Paginate(
		Config.DB.Where("deleted_at IS NULL AND status = ?", status).Order("id"),
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
