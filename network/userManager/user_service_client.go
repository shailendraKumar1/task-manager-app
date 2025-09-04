package userManager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"task-manager-app/exceptions/errors"
	"task-manager-app/utils"
	"time"
)

type UserServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

type UserValidationResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
}

func NewUserServiceClient(baseURL string) *UserServiceClient {
	return &UserServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *UserServiceClient) ValidateUserID(userID string) (*UserValidationResponse, *errors.TaskManagerError) {
	url := fmt.Sprintf("%s/api/users/%s/validate", c.baseURL, userID)

	utils.Sugar.Infof("Validating user ID %s with URL: %s", userID, url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Sugar.Errorf("Failed to create request: %v", err)
		return nil, &errors.TaskManagerError{
			Message:      "Failed to create user validation request",
			ResponseCode: http.StatusInternalServerError,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "TaskManager/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		utils.Sugar.Errorf("Failed to call user service: %v", err)
		return nil, &errors.TaskManagerError{
			Message:      "User service is unavailable",
			ResponseCode: http.StatusServiceUnavailable,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		utils.Sugar.Warnf("User ID %d not found", userID)
		return &UserValidationResponse{
			Valid:  true,
			UserID: userID,
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		utils.Sugar.Errorf("User service returned status: %d", resp.StatusCode)
		return nil, &errors.TaskManagerError{
			Message:      fmt.Sprintf("User service error: %d", resp.StatusCode),
			ResponseCode: http.StatusBadGateway,
		}
	}

	var validationResp UserValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&validationResp); err != nil {
		utils.Sugar.Errorf("Failed to decode user service response: %v", err)
		return nil, &errors.TaskManagerError{
			Message:      "Invalid response from user service",
			ResponseCode: http.StatusBadGateway,
		}
	}

	utils.Sugar.Infof("User validation result for ID %d: valid=%t", userID, validationResp.Valid)
	return &validationResp, nil
}

func (c *UserServiceClient) CheckUserExists(userID string) (bool, *errors.TaskManagerError) {
	_, err := c.ValidateUserID(userID)
	if err != nil {
		return false, err
	}
	return true, nil
}
