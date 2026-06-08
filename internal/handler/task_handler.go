package handler

import (
	"net/http"
	"task-scheduler/internal/dto"
	"task-scheduler/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	task, err := h.taskService.CreateTask(&req)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	success(c, task)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := parseUint64Param(c, "id")
	if err != nil {
		fail(c, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := h.taskService.GetTask(id)
	if err != nil {
		fail(c, http.StatusNotFound, "task not found")
		return
	}

	success(c, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := parseUint64Param(c, "id")
	if err != nil {
		fail(c, http.StatusBadRequest, "invalid task id")
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	task, err := h.taskService.UpdateTask(id, &req)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	success(c, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := parseUint64Param(c, "id")
	if err != nil {
		fail(c, http.StatusBadRequest, "invalid task id")
		return
	}

	err = h.taskService.DeleteTask(id)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	success(c, nil)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	var req dto.TaskListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	result, err := h.taskService.ListTasks(&req)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	success(c, result)
}

func (h *TaskHandler) TriggerTask(c *gin.Context) {
	id, err := parseUint64Param(c, "id")
	if err != nil {
		fail(c, http.StatusBadRequest, "invalid task id")
		return
	}

	err = h.taskService.TriggerTask(id)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	success(c, nil)
}

func (h *TaskHandler) PauseTask(c *gin.Context) {
	id, err := parseUint64Param(c, "id")
	if err != nil {
		fail(c, http.StatusBadRequest, "invalid task id")
		return
	}

	err = h.taskService.PauseTask(id)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	success(c, nil)
}

func (h *TaskHandler) ResumeTask(c *gin.Context) {
	id, err := parseUint64Param(c, "id")
	if err != nil {
		fail(c, http.StatusBadRequest, "invalid task id")
		return
	}

	err = h.taskService.ResumeTask(id)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	success(c, nil)
}

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	tasks := router.Group("/tasks")
	{
		tasks.POST("", h.CreateTask)
		tasks.GET("", h.ListTasks)
		tasks.GET("/:id", h.GetTask)
		tasks.PUT("/:id", h.UpdateTask)
		tasks.DELETE("/:id", h.DeleteTask)
		tasks.POST("/:id/trigger", h.TriggerTask)
		tasks.POST("/:id/pause", h.PauseTask)
		tasks.POST("/:id/resume", h.ResumeTask)
	}
}
