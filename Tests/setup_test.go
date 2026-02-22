// File PATH: D:\myProject\Task Manager\Tests\setup_test.go
package Tests

import (
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite" // درایور جدید که به CGO نیازی ندارد
	"github.com/vahidlotfi71/Task_Manager/Config"
	"github.com/vahidlotfi71/Task_Manager/Models"
	"github.com/vahidlotfi71/Task_Manager/Routes"
	"gorm.io/gorm"
)

// SetupTestRouter یک روتر تستی و دیتابیس ایزوله برای تست می‌سازد
func SetupTestRouter() *gin.Engine {
	// ۱) قرار دادن Gin در حالت تست
	gin.SetMode(gin.TestMode)

	// ۲) اتصال به دیتابیس SQLite در حافظه با استفاده از درایور glebarez
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	// ۳) جایگزینی دیتابیس اصلی با دیتابیس تستی
	Config.DB = db

	// ۴) AutoMigrate برای ساخت جداول در حافظه
	err = Config.DB.AutoMigrate(&Models.Task{})
	if err != nil {
		panic("Failed to migrate test database: " + err.Error())
	}

	// ۵) راه‌اندازی روترها بر اساس ساختار پروژه شما
	r := gin.Default()
	Routes.SetupRoutes(r)

	return r
}

// ClearDatabase جداول را بعد از هر تست پاک می‌کند
func ClearDatabase() {
	Config.DB.Exec("DELETE FROM tasks")
}
