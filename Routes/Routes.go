// File: Routes/Routes.go
package Routes

import (
	"github.com/gin-gonic/gin"
	Controller "github.com/vahidlotfi71/Task_Manager/Controllers"
	Validation "github.com/vahidlotfi71/Task_Manager/Validations"
)

func SetupRoutes(r *gin.Engine) {

	// ---------- TASK GROUP ----------
	task := r.Group("/task")

	task.GET("/", Controller.Index)                                  // task.index
	task.GET("/show/:id", Controller.Show)                           // task.show
	task.POST("/store", Validation.Store(), Controller.Store)        // task.store
	task.POST("/update/:id", Validation.Update(), Controller.Update) // task.update
	task.POST("/delete/:id", Controller.Delete)                      // task.delete
	task.GET("/restore/:id", Controller.Restore)                     // task.restore
	task.GET("/trash", Controller.Trash)                             // task.trash
	task.POST("/clear-trash", Controller.ClearTrash)                 // task.clear-trash

	// فیلتر بر اساس Status
	task.GET("/filter/status/:status", Controller.FilterByStatus)

	// فیلتر بر اساس Assignee
	task.GET("/filter/assignee/:assignee", Controller.FilterByAssignee)
	// ---------- 404 HANDLER ----------
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Route Not Found"})
	})
}
