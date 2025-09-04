package controller

import (
	"net/http"
	"strconv"
	"task-manager-app/constants"
	"task-manager-app/exceptions"
	"task-manager-app/request"
	"task-manager-app/services/taskManagerService"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	service taskManagerService.TaskService
}

func NewTaskController(service taskManagerService.TaskService) *TaskController {
	return &TaskController{service: service}
}

func (c *TaskController) CreateTask(ctx *gin.Context) {
	var req request.ReqCreateOrUpdateTasks
	if err := ctx.ShouldBindJSON(&req); err != nil {
		taskErr := exceptions.NewBadRequestException(constants.ErrInvalidRequestBody + ": " + err.Error())
		ctx.JSON(taskErr.ResponseCode, taskErr)
		return
	}

	resp, taskErr := c.service.CreateTask(&req)
	if taskErr != nil {
		ctx.JSON(taskErr.ResponseCode, taskErr)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

func (c *TaskController) GetTask(ctx *gin.Context) {
	uuid := ctx.Param(constants.URLParamUUID)
	resp, taskErr := c.service.GetTaskByUUID(uuid)
	if taskErr != nil {
		ctx.JSON(taskErr.ResponseCode, taskErr)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *TaskController) UpdateTask(ctx *gin.Context) {
	uuid := ctx.Param(constants.URLParamUUID)
	var req request.ReqCreateOrUpdateTasks
	if err := ctx.ShouldBindJSON(&req); err != nil {
		taskErr := exceptions.NewBadRequestException(constants.ErrInvalidRequestBody + ": " + err.Error())
		ctx.JSON(taskErr.ResponseCode, taskErr)
		return
	}

	// Update task with validation in service
	resp, taskErr := c.service.UpdateTask(uuid, &req)
	if taskErr != nil {
		ctx.JSON(taskErr.ResponseCode, taskErr)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *TaskController) DeleteTask(ctx *gin.Context) {
	uuid := ctx.Param(constants.URLParamUUID)
	if taskErr := c.service.DeleteTask(uuid); taskErr != nil {
		ctx.JSON(taskErr.ResponseCode, taskErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *TaskController) ListTasks(ctx *gin.Context) {
	status := ctx.Query(constants.QueryParamStatus)
	userID := ctx.Query(constants.QueryParamUserID)
	priority := ctx.Query(constants.QueryParamPriority)
	pageStr := ctx.DefaultQuery(constants.QueryParamPage, constants.DefaultPageStr)
	sizeStr := ctx.DefaultQuery(constants.QueryParamPageSize, constants.DefaultPageSizeStr)

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(sizeStr)

	resp, taskErr := c.service.ListTasks(status, userID, priority, page, pageSize)
	if taskErr != nil {
		ctx.JSON(taskErr.ResponseCode, taskErr)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
