package service

import (
	"task-scheduler/internal/dto"
	"task-scheduler/internal/model"
	"task-scheduler/internal/repository"
)

type TaskLogService interface {
	GetTaskLog(id uint64) (*model.TaskLog, error)
	ListTaskLogs(req *dto.TaskLogListRequest) (*dto.PageResult, error)
}

type taskLogService struct {
	logRepo repository.TaskLogRepository
}

func NewTaskLogService(logRepo repository.TaskLogRepository) TaskLogService {
	return &taskLogService{logRepo: logRepo}
}

func (s *taskLogService) GetTaskLog(id uint64) (*model.TaskLog, error) {
	return s.logRepo.GetByID(id)
}

func (s *taskLogService) ListTaskLogs(req *dto.TaskLogListRequest) (*dto.PageResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	logs, total, err := s.logRepo.ListByTaskID(req.TaskID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:     logs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
