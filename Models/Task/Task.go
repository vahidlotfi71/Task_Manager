package Task

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vahidlotfi71/Task_Manager/Models"
	"github.com/vahidlotfi71/Task_Manager/Utils"
	"gorm.io/gorm"
)

type TaskCreateDTO struct {
	Title       string
	Description string
	Status      Models.TaskStatus
	Assignee    string
}

type TaskUpdateDTO struct {
	Title       string
	Description string
	Status      Models.TaskStatus
	Assignee    string
}

func Paginate(tx *gorm.DB, c *gin.Context) (tasks []Models.Task, meta Utils.PaginationMetadata, err error) {
	if tx.Statement.Table == "" {
		tx = tx.Model(&Models.Task{})
	}
	tx, meta = Utils.Paginate(tx, c)
	err = tx.Find(&tasks).Error
	return
}
func FindByID(tx *gorm.DB, id int64) (task Models.Task, err error) {
	err = tx.Where("deleted_at IS NULL").First(&task, id).Error
	return
}

func FindByStatus(tx *gorm.DB, status Models.TaskStatus, c *gin.Context) (tasks []Models.Task, meta Utils.PaginationMetadata, err error) {
	tx = tx.Model(&Models.Task{}).Where("status = ?", status)
	tx, meta = Utils.Paginate(tx, c)
	err = tx.Find(&tasks).Error
	return
}

func FindByAssignee(tx *gorm.DB, assignee string, c *gin.Context) (tasks []Models.Task, meta Utils.PaginationMetadata, err error) {
	tx = tx.Model(&Models.Task{}).Where("assignee = ?", assignee)
	tx, meta = Utils.Paginate(tx, c)
	err = tx.Find(&tasks).Error
	return
}

func Create(tx *gorm.DB, dto TaskCreateDTO) (task Models.Task, err error) {
	status := dto.Status
	if status == "" {
		status = Models.StatusPending
	}

	task = Models.Task{
		Title:       dto.Title,
		Description: dto.Description,
		Status:      status,
		Assignee:    dto.Assignee,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = tx.Create(&task).Error
	return
}

func Update(tx *gorm.DB, id int64, dto TaskUpdateDTO) error {
	updates := map[string]interface{}{
		"title":       dto.Title,
		"description": dto.Description,
		"assignee":    dto.Assignee,
		"updated_at":  time.Now(),
	}

	if dto.Status != "" {
		updates["status"] = dto.Status
	}

	result := tx.Model(&Models.Task{}).
		Where("id = ? AND deleted_at IS NULL", id).Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found or already deleted")
	}
	return nil
}

func SoftDelete(tx *gorm.DB, id int64) error {
	result := tx.Model(&Models.Task{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", time.Now())

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found or already deleted")
	}
	return nil
}
