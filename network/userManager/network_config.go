package userManager

import (
	"os"
	"task-manager-app/utils"
)

var (
	UserClient *UserServiceClient
)

func InitNetworkClients() {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:8081"
	}

	UserClient = NewUserServiceClient(userServiceURL)

	utils.Sugar.Infof("Network clients initialized successfully")
}
