package Controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Config"
	"github.com/vahidlotfi71/Task_Manager/Models"
	Task "github.com/vahidlotfi71/Task_Manager/Models/Repository"
	Resource "github.com/vahidlotfi71/Task_Manager/Resources"
)

/* ---------- DTO درخواست ---------- */
type TaskUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Assignee    string `json:"assignee"`
}

func Update(c *gin.Context) {
	// ۱) پارس ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid task ID"})
		return
	}

	// ۲) خواندن بدنه درخواست
	body, _ := c.Get("json_body")
	jsonBody := body.(map[string]interface{})

	var req TaskUpdateRequest
	req.Title = strings.TrimSpace(getString(jsonBody, "title"))
	req.Description = strings.TrimSpace(getString(jsonBody, "description"))
	req.Status = strings.TrimSpace(getString(jsonBody, "status"))
	req.Assignee = strings.TrimSpace(getString(jsonBody, "assignee"))

	// ۳) شروع تراکنش
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

	// ۴) چک وجود task (فقط حذف‌نشده‌ها)
	existing, err := Task.FindByID(tx, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	// ۵) اگر فیلدی خالی بود، از مدل فعلی بخوان
	if req.Title == "" {
		req.Title = existing.Title
	}
	if req.Description == "" {
		req.Description = existing.Description
	}
	if req.Status == "" {
		req.Status = string(existing.Status)
	}
	if req.Assignee == "" {
		req.Assignee = existing.Assignee
	}

	// ۶) ساخت DTO
	dto := Task.TaskUpdateDTO{
		Title:       req.Title,
		Description: req.Description,
		Status:      Models.TaskStatus(req.Status),
		Assignee:    req.Assignee,
	}

	// ۷) به‌روزرسانی
	if err := Task.Update(tx, id, dto); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// ۸) کامیت تراکنش
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Commit failed"})
		return
	}

	// ۹) بارگذاری مجدد برای داشتن داده‌های تازه
	freshTask, err := Task.FindByID(Config.DB, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to reload task"})
		return
	}

	// ۱۰) پاسخ
	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"data":    Resource.Single(freshTask),
	})
}
