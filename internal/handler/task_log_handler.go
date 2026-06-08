package handler

import (
	"net/http"
	"task-scheduler/internal/dto"
	"task-scheduler/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskLogHandler struct {
	logService service.TaskLogService
}

func NewTaskLogHandler(logService service.TaskLogService) *TaskLogHandler {
	return &TaskLogHandler{logService: logService}
}

func (h *TaskLogHandler) GetTaskLog(c *gin.Context) {
	id, err := parseUint64Param(c, "id")
	if err != nil {
		fail(c, http.StatusBadRequest, "invalid log id")
		return
	}

	log, err := h.logService.GetTaskLog(id)
	if err != nil {
		fail(c, http.StatusNotFound, "log not found")
		return
	}

	success(c, log)
}

func (h *TaskLogHandler) ListTaskLogs(c *gin.Context) {
	var req dto.TaskLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	result, err := h.logService.ListTaskLogs(&req)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	success(c, result)
}

func (h *TaskLogHandler) RegisterRoutes(router *gin.RouterGroup) {
	logs := router.Group("/task-logs")
	{
		logs.GET("/:id", h.GetTaskLog)
		logs.GET("", h.ListTaskLogs)
	}
}
