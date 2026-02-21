package Controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Config"
	"github.com/vahidlotfi71/Task_Manager/Models"
	"github.com/vahidlotfi71/Task_Manager/Resources"
	"gorm.io/gorm"
)

func Restore(c *gin.Context) {
	// ۱) استخراج و اعتبارسنجی شناسه
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id param is required"})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id must be a positive integer"})
		return
	}

	// ۲) جستجوی task با Unscoped (شامل رکوردهای حذف‌شده)
	var task Models.Task
	if err := Config.DB.Unscoped().
		Where("id = ?", id).
		First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// ۳) بررسی اینکه task واقعاً حذف شده باشد
	if task.DeletedAt.Time.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Task is not deleted"})
		return
	}

	// ۴) بازگردانی با Unscoped
	result := Config.DB.Unscoped().
		Model(&Models.Task{}).
		Where("id = ?", id).
		Update("deleted_at", nil)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	// ۵) خواندن دوباره مدل برای پاسخ
	if err := Config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// ۶) پاسخ موفق
	c.JSON(http.StatusOK, gin.H{
		"message": "Task restored successfully",
		"data":    Resources.Single(task),
	})
}
