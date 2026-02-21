package Controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Config"
	"github.com/vahidlotfi71/Task_Manager/Models"
	"gorm.io/gorm"
)

func ClearTrash(c *gin.Context) {
	// ۱) خواندن تعداد درخواستی (با سقف امن)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	// ۲) شروع تراکنش
	tx := Config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "DB connection error"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ۳) خواندن رکوردهای حذف‌شده
	var tasks []Models.Task
	if err := tx.Unscoped().
		Where("deleted_at IS NOT NULL").
		Order("deleted_at ASC").
		Limit(limit).
		Find(&tasks).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// ۴) استخراج IDهای taskها برای حذف
	var taskIDs []uint
	for _, task := range tasks {
		taskIDs = append(taskIDs, uint(task.ID))
	}

	// ۵) حذف فیزیکی دسته‌ای
	var result *gorm.DB
	if len(taskIDs) > 0 {
		result = tx.Unscoped().Delete(&Models.Task{}, "id IN ?", taskIDs)
	} else {
		result = &gorm.DB{RowsAffected: 0}
	}

	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	// ۶) کامیت موفق
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Commit failed"})
		return
	}

	// ۷) پاسخ موفق
	c.JSON(http.StatusOK, gin.H{
		"message":       "Trash cleared successfully",
		"cleared_count": result.RowsAffected,
	})
}
