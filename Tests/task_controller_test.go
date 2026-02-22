package Tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vahidlotfi71/Task_Manager/Config"
	"github.com/vahidlotfi71/Task_Manager/Models"
	"github.com/vahidlotfi71/Task_Manager/Models/Repository"
	"gorm.io/gorm"
)

// TestStoreTask: Test successful creation of a task
func TestStoreTask(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	reqBody := map[string]interface{}{
		"title":       "Task Creation Test",
		"description": "This is a test task",
		"status":      "pending",
		"assignee":    "Ali",
	}
	jsonValue, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/task/store", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Task created successfully", response["message"])
}

// TestStoreTaskValidationError: Test validation error when creating a task
func TestStoreTaskValidationError(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	reqBody := map[string]interface{}{"description": "Only description"} // Title is not sent
	jsonValue, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/task/store", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestIndexTasks: Test getting the list of tasks with Pagination
func TestIndexTasks(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	Config.DB.Create(&Models.Task{Title: "Task 1"})
	Config.DB.Create(&Models.Task{Title: "Task 2"})

	req, _ := http.NewRequest("GET", "/task/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotNil(t, response["data"])
	assert.NotNil(t, response["metadata"])
	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))
}

// TestShowTask: Test showing a specific task
func TestShowTask(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	task := Models.Task{Title: "Display Task"}
	Config.DB.Create(&task)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/task/show/%d", task.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Display Task", data["title"])
}

// TestShowTaskNotFound: Test showing a task that does not exist
func TestShowTaskNotFound(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	req, _ := http.NewRequest("GET", "/task/show/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestFilterTasksByStatus: Test filtering by status
func TestFilterTasksByStatus(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	Config.DB.Create(&Models.Task{Title: "Done Task", Status: "done"})
	Config.DB.Create(&Models.Task{Title: "In Progress Task", Status: "in_progress"})

	req, _ := http.NewRequest("GET", "/task/filter/status/done", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Equal(t, 1, len(data))
	assert.Equal(t, "done", data[0].(map[string]interface{})["status"])
}

// TestFilterTasksByAssignee: Test filtering by assignee
func TestFilterTasksByAssignee(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	Config.DB.Create(&Models.Task{Title: "Ali's Task", Assignee: "Ali"})
	Config.DB.Create(&Models.Task{Title: "Reza's Task", Assignee: "Reza"})

	req, _ := http.NewRequest("GET", "/task/filter/assignee/Ali", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Equal(t, 1, len(data))
	assert.Equal(t, "Ali", data[0].(map[string]interface{})["assignee"])
}

// TestUpdateTask: Test successful update of a task
func TestUpdateTask(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	task := Models.Task{Title: "Old Title", Status: "pending"}
	Config.DB.Create(&task)

	reqBody := map[string]interface{}{"title": "New Title", "status": "done"}
	jsonValue, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/task/update/%d", task.ID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedTask Models.Task
	Config.DB.First(&updatedTask, task.ID)
	assert.Equal(t, "New Title", updatedTask.Title)
	assert.Equal(t, "done", string(updatedTask.Status))
}

// TestUpdateTaskNotFound: Test updating a task that does not exist
func TestUpdateTaskNotFound(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	reqBody := map[string]interface{}{"title": "New Title", "status": "pending"}
	jsonValue, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/task/update/999", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestUpdateTaskValidation: Test validation error when updating
func TestUpdateTaskValidation(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	task := Models.Task{Title: "Initial Task"}
	Config.DB.Create(&task)

	reqBody := map[string]interface{}{"status": "invalid_status"}
	jsonValue, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/task/update/%d", task.ID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestSoftDeleteTask: Test soft delete
func TestSoftDeleteTask(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	task := Models.Task{Title: "Task to delete"}
	Config.DB.Create(&task)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/task/delete/%d", task.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var deletedTask Models.Task
	err := Config.DB.First(&deletedTask, task.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	err = Config.DB.Unscoped().First(&deletedTask, task.ID).Error
	assert.NoError(t, err)
	assert.True(t, deletedTask.DeletedAt.Valid)
}

// TestRestoreTask: Test restoring a deleted task
func TestRestoreTask(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	task := Models.Task{Title: "Task to restore"}
	Config.DB.Create(&task)
	Config.DB.Delete(&task) // Soft delete

	req, _ := http.NewRequest("GET", fmt.Sprintf("/task/restore/%d", task.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var restoredTask Models.Task
	Config.DB.First(&restoredTask, task.ID)
	assert.False(t, restoredTask.DeletedAt.Valid)
}

// TestRestoreNotDeletedTask: Test restoring a task that is not deleted
func TestRestoreNotDeletedTask(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	task := Models.Task{Title: "Undeleted task"}
	Config.DB.Create(&task)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/task/restore/%d", task.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestTrash: Test viewing deleted tasks
func TestTrash(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	task1 := Models.Task{Title: "Deleted 1"}
	task2 := Models.Task{Title: "Not deleted"}
	Config.DB.Create(&task1)
	Config.DB.Create(&task2)
	Config.DB.Delete(&task1)

	req, _ := http.NewRequest("GET", "/task/trash", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.Equal(t, 1, len(data))
}

// TestClearTrash: Test completely clearing the trash
func TestClearTrash(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	task1 := Models.Task{Title: "For permanent deletion"}
	Config.DB.Create(&task1)
	Config.DB.Delete(&task1)

	reqBody := map[string]interface{}{"limit": "10"}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/task/clear-trash", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(1), response["cleared_count"])

	var count int64
	Config.DB.Unscoped().Model(&Models.Task{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

// ==========================================
// Repository Tests
// ==========================================

// TestRepositorySoftDelete: Test direct soft delete from repository
// TestRepositorySoftDelete: Test direct soft delete from repository
func TestRepositorySoftDelete(t *testing.T) {
	SetupTestRouter()
	defer ClearDatabase()

	// Create DTO instead of direct model
	dto := Repository.TaskCreateDTO{
		Title: "Repository Test",
	}

	// Get created task and error
	task, err := Repository.Create(Config.DB, dto)
	assert.NoError(t, err)

	err = Repository.SoftDelete(Config.DB, task.ID)
	assert.NoError(t, err)

	_, err = Repository.FindByID(Config.DB, task.ID)
	assert.Error(t, err) // Should not be found
}

// TestRepositoryRestore: Test direct restore from repository
// TestRepositoryRestore: Test direct restore from repository
func TestRepositoryRestore(t *testing.T) {
	SetupTestRouter()
	defer ClearDatabase()

	// Create DTO instead of direct model
	dto := Repository.TaskCreateDTO{
		Title: "Repository Test",
	}

	// Get created task and error
	task, err := Repository.Create(Config.DB, dto)
	assert.NoError(t, err)

	Repository.SoftDelete(Config.DB, task.ID)

	err = Repository.Restore(Config.DB, task.ID)
	assert.NoError(t, err)

	restoredTask, err := Repository.FindByID(Config.DB, task.ID)
	assert.NoError(t, err)
	assert.False(t, restoredTask.DeletedAt.Valid)
}

// TestRepositoryClearTrash: Test clearing trash from repository
// TestRepositoryClearTrash: Test clearing trash from repository
// TestRepositoryClearTrash: Test clearing trash from repository
// TestRepositoryClearTrash: Test clearing trash from repository
func TestRepositoryClearTrash(t *testing.T) {
	SetupTestRouter()
	defer ClearDatabase()

	dto := Repository.TaskCreateDTO{
		Title: "Repository Test",
	}

	task, err := Repository.Create(Config.DB, dto)
	assert.NoError(t, err)

	Repository.SoftDelete(Config.DB, task.ID)

	// Final fix: This function only takes the database connection (without limit parameter)
	err = Repository.ClearTrash(Config.DB)
	assert.NoError(t, err) // Expect no error to occur

	var totalCount int64
	Config.DB.Unscoped().Model(&Models.Task{}).Count(&totalCount)
	assert.Equal(t, int64(0), totalCount)
}

// ==========================================
// Additional tests to increase Coverage
// ==========================================

// TestSoftDeleteTaskNotFound: Attempt to delete a task that does not exist
func TestSoftDeleteTaskNotFound(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	// Send invalid ID for deletion
	req, _ := http.NewRequest("POST", "/task/delete/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Expect a 404 error to be returned
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestRestoreTaskNotFound: Attempt to restore a task that is not in the database
func TestRestoreTaskNotFound(t *testing.T) {
	r := SetupTestRouter()
	defer ClearDatabase()

	// Send invalid ID for restore
	req, _ := http.NewRequest("GET", "/task/restore/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRepositoryFindByIDNotFound(t *testing.T) {
	SetupTestRouter()
	defer ClearDatabase()
	_, err := Repository.FindByID(Config.DB, 999)

	assert.Error(t, err)
}
