package Controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Config"
	"github.com/vahidlotfi71/Task_Manager/Models"
	"gorm.io/gorm"
)

func Delete(c *gin.Context) {
	// ۱) استخراج و اعتبارسنجی شناسه
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id param is required"})
		return
	}
	num, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id must be a positive integer"})
		return
	}
	id := uint(num)

	// ۲) چک وجود رکورد (فقط حذف‌نشده‌ها)
	var task Models.Task
	if err := Config.DB.Where("deleted_at IS NULL").First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// ۳) Soft-Delete با ORM + RowsAffected
	result := Config.DB.Model(&Models.Task{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	// ۴) پاسخ موفق
	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}
