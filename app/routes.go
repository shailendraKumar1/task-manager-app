package app

import (
	"task-manager-app/controller"
	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(router *gin.Engine, taskController *controller.TaskController) {
	tasks := router.Group("/tasks")
	{
		tasks.POST("", taskController.CreateTask)
		tasks.GET("", taskController.ListTasks)
		tasks.GET("/:uuid", taskController.GetTask)
		tasks.PUT("/:uuid", taskController.UpdateTask)
		tasks.DELETE("/:uuid", taskController.DeleteTask)
	}
}

func RegisterHealthRoutes(router *gin.Engine, healthController *controller.HealthController) {
	// Health check endpoint
	router.GET("/health", healthController.HealthCheck)
}
