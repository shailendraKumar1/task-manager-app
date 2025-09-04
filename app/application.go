package app

import (
	"task-manager-app/config"
	"task-manager-app/controller"
	"task-manager-app/network/userManager"
	"task-manager-app/repo"
	"task-manager-app/services/taskManagerService"
	"task-manager-app/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

var (
	appName string
	version string
	router  = gin.Default()
)

func StartApplication() {

	// Load environment variables from .env file
	err := godotenv.Load("resources/.env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	utils.InitLogger()

	config.ApplicationConfig.SetAwsSecretValues()

	// Initialize database using GORM
	config.InitDB()

	// Initialize network clients
	userManager.InitNetworkClients()

	appName = config.ApplicationConfig.AppName
	version = config.ApplicationConfig.AppVersion
	utils.Sugar.Infow("Starting application: ", appName, version)

	// Initialize repositories, services, and controllers
	taskRepo := repo.NewTaskRepository(config.DB)
	taskService := taskManagerService.NewTaskService(taskRepo)
	taskController := controller.NewTaskController(taskService)
	healthController := controller.NewHealthController(config.DB)

	// Register routes
	RegisterTaskRoutes(router, taskController)
	RegisterHealthRoutes(router, healthController)

	runErr := router.Run(config.ApplicationConfig.AppHost + ":" + config.ApplicationConfig.AppPort)
	if runErr != nil {
		utils.Sugar.Fatal("Error starting application: ", runErr.Error())
	}
}
