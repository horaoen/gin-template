package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/horaoen/go-backend-clean-architecture/domain"
)

type TaskController struct {
	TaskUsecase domain.TaskUsecase
}

// Create godoc
// @Summary Create Task
// @Description Create a new task
// @Tags Task
// @Security BearerAuth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param request formData domain.Task true "Task Details"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /task [post]
func (tc *TaskController) Create(c *gin.Context) {
	var task domain.Task

	err := c.ShouldBind(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	userID := c.GetString("x-user-id")

	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	task.UserID = uint(uid)

	err = tc.TaskUsecase.Create(c, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "Task created successfully",
	})
}

// Fetch godoc
// @Summary Get Tasks
// @Description Get all tasks for the user
// @Tags Task
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Task
// @Failure 500 {object} domain.ErrorResponse
// @Router /task [get]
func (tc *TaskController) Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")

	tasks, err := tc.TaskUsecase.FetchByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
