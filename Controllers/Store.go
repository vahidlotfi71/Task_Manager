package Controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Config"
	"github.com/vahidlotfi71/Task_Manager/Models"
	"github.com/vahidlotfi71/Task_Manager/Models/Repository"
	"github.com/vahidlotfi71/Task_Manager/Resources"
)

/* ---------- DTO درخواست ---------- */
type TaskCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Assignee    string `json:"assignee"`
}

func Store(c *gin.Context) {
	// ۱) خواندن بدنه درخواست
	body, _ := c.Get("json_body")
	jsonBody := body.(map[string]interface{})

	req := TaskCreateRequest{
		Title:       strings.TrimSpace(getString(jsonBody, "title")),
		Description: strings.TrimSpace(getString(jsonBody, "description")),
		Status:      strings.TrimSpace(getString(jsonBody, "status")),
		Assignee:    strings.TrimSpace(getString(jsonBody, "assignee")),
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

	// ۳) ساخت DTO برای Repository
	dto := Repository.TaskCreateDTO{
		Title:       req.Title,
		Description: req.Description,
		Status:      Models.TaskStatus(req.Status),
		Assignee:    req.Assignee,
	}

	// ۴) درج در DB
	task, err := Repository.Create(tx, dto)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// ۵) کامیت موفق
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Commit failed"})
		return
	}

	// ۶) پاسخ استاندارد
	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"data":    Resources.Single(task),
	})
}
