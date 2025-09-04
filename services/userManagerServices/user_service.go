package userManagerServices

import (
	"fmt"
	"task-manager-app/network/userManager"
)

type UserService interface {
	ValidateUser(userID string) (bool, error)
}

type userService struct {
	userClient *userManager.UserServiceClient
}

func NewUserService() UserService {
	return &userService{
		userClient: userManager.UserClient,
	}
}

// ValidateUser validates if a user ID exists in the user service
func (s *userService) ValidateUser(userID string) (bool, error) {
	if s.userClient == nil {
		return false, fmt.Errorf("user service client not initialized")
	}

	exists, clientErr := s.userClient.CheckUserExists(userID)
	if clientErr != nil {
		return false, fmt.Errorf("failed to validate user ID: %s", clientErr.Message)
	}

	return exists, nil
}
