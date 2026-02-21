package Controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Config"
	"github.com/vahidlotfi71/Task_Manager/Models/Repository"
	"github.com/vahidlotfi71/Task_Manager/Resources"
)

func Index(c *gin.Context) {
	tasks, meta, err := Repository.Paginate(
		Config.DB.Where("deleted_at IS NULL").Order("id"),
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
