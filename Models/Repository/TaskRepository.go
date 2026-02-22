package Repository

import (
	"errors"

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
	err = tx.First(&task, id).Error
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
		// CreatedAt and UpdatedAt are managed automatically by GORM
	}
	err = tx.Create(&task).Error
	return
}

func Update(tx *gorm.DB, id int64, dto TaskUpdateDTO) error {
	updates := map[string]interface{}{}

	if dto.Title != "" {
		updates["title"] = dto.Title
	}
	if dto.Description != "" {
		updates["description"] = dto.Description
	}
	if dto.Assignee != "" {
		updates["assignee"] = dto.Assignee
	}
	if dto.Status != "" {
		updates["status"] = dto.Status
	}

	if len(updates) == 0 {
		return errors.New("no fields provided to update")
	}

	result := tx.Model(&Models.Task{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found or already deleted")
	}
	return nil
}

func SoftDelete(tx *gorm.DB, id int64) error {
	// Use GORM's built-in Delete which respects gorm.DeletedAt (soft delete)
	result := tx.Where("id = ?", id).Delete(&Models.Task{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found or already deleted")
	}
	return nil
}



func Restore(tx *gorm.DB, id int64) error {
	result := tx.Unscoped().
		Model(&Models.Task{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found in trash")
	}
	return nil
}

func ClearTrash(tx *gorm.DB) error {
	return tx.Unscoped().
		Where("deleted_at IS NOT NULL").
		Delete(&Models.Task{}).Error
}
