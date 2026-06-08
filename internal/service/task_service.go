package service

import (
	"errors"
	"task-scheduler/internal/dto"
	"task-scheduler/internal/model"
	"task-scheduler/internal/repository"
	"time"

	"github.com/robfig/cron/v3"
)

type TaskService interface {
	CreateTask(req *dto.CreateTaskRequest) (*model.Task, error)
	GetTask(id uint64) (*model.Task, error)
	UpdateTask(id uint64, req *dto.UpdateTaskRequest) (*model.Task, error)
	DeleteTask(id uint64) error
	ListTasks(req *dto.TaskListRequest) (*dto.PageResult, error)
	TriggerTask(id uint64) error
	PauseTask(id uint64) error
	ResumeTask(id uint64) error
	GetDueTasks(limit int) ([]model.Task, error)
	ClaimTask(taskID uint64, workerID string) (bool, error)
	CompleteTask(taskID uint64, success bool, result, errorMsg string) error
	CalculateNextTime(cronExpr string) (*time.Time, error)
}

type taskService struct {
	taskRepo repository.TaskRepository
	logRepo  repository.TaskLogRepository
}

func NewTaskService(taskRepo repository.TaskRepository, logRepo repository.TaskLogRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
		logRepo:  logRepo,
	}
}

func (s *taskService) CreateTask(req *dto.CreateTaskRequest) (*model.Task, error) {
	task := &model.Task{
		Name:          req.Name,
		Type:          model.TaskType(req.Type),
		CronExpr:      req.CronExpr,
		Payload:       req.Payload,
		MaxRetryTimes: req.MaxRetryTimes,
		RetryInterval: req.RetryInterval,
		Timeout:       req.Timeout,
		Status:        model.TaskStatusPending,
		RetryTimes:    0,
	}

	if task.MaxRetryTimes == 0 {
		task.MaxRetryTimes = 3
	}
	if task.RetryInterval == 0 {
		task.RetryInterval = 60
	}
	if task.Timeout == 0 {
		task.Timeout = 300
	}

	if task.Type == model.TaskTypeCron && task.CronExpr != "" {
		nextTime, err := s.CalculateNextTime(task.CronExpr)
		if err != nil {
			return nil, err
		}
		task.NextExecuteTime = nextTime
	} else if task.Type == model.TaskTypeOneTime {
		now := time.Now()
		task.NextExecuteTime = &now
	}

	err := s.taskRepo.Create(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) GetTask(id uint64) (*model.Task, error) {
	return s.taskRepo.GetByID(id)
}

func (s *taskService) UpdateTask(id uint64, req *dto.UpdateTaskRequest) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		task.Name = req.Name
	}
	if req.CronExpr != "" {
		task.CronExpr = req.CronExpr
		if task.Type == model.TaskTypeCron {
			nextTime, err := s.CalculateNextTime(task.CronExpr)
			if err != nil {
				return nil, err
			}
			task.NextExecuteTime = nextTime
		}
	}
	if req.Payload != "" {
		task.Payload = req.Payload
	}
	if req.MaxRetryTimes != nil {
		task.MaxRetryTimes = *req.MaxRetryTimes
	}
	if req.RetryInterval != nil {
		task.RetryInterval = *req.RetryInterval
	}
	if req.Timeout != nil {
		task.Timeout = *req.Timeout
	}

	err = s.taskRepo.Update(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) DeleteTask(id uint64) error {
	return s.taskRepo.Delete(id)
}

func (s *taskService) ListTasks(req *dto.TaskListRequest) (*dto.PageResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	var status *model.TaskStatus
	if req.Status != nil {
		s := model.TaskStatus(*req.Status)
		status = &s
	}

	tasks, total, err := s.taskRepo.List(req.Page, req.PageSize, status)
	if err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:     tasks,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *taskService) TriggerTask(id uint64) error {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	task.NextExecuteTime = &now
	task.Status = model.TaskStatusPending
	task.RetryTimes = 0

	return s.taskRepo.Update(task)
}

func (s *taskService) PauseTask(id uint64) error {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return err
	}

	if task.Status == model.TaskStatusRunning {
		return errors.New("task is running, cannot pause")
	}

	task.Status = model.TaskStatusPaused
	return s.taskRepo.Update(task)
}

func (s *taskService) ResumeTask(id uint64) error {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return err
	}

	if task.Status != model.TaskStatusPaused {
		return errors.New("task is not paused")
	}

	if task.Type == model.TaskTypeCron && task.CronExpr != "" {
		nextTime, err := s.CalculateNextTime(task.CronExpr)
		if err != nil {
			return err
		}
		task.NextExecuteTime = nextTime
	}

	task.Status = model.TaskStatusPending
	return s.taskRepo.Update(task)
}

func (s *taskService) GetDueTasks(limit int) ([]model.Task, error) {
	return s.taskRepo.GetDueTasks(time.Now(), limit)
}

func (s *taskService) ClaimTask(taskID uint64, workerID string) (bool, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return false, err
	}

	if task.Status == model.TaskStatusRunning {
		return false, nil
	}

	success, err := s.taskRepo.ClaimTask(taskID, workerID, task.Version)
	if err != nil {
		return false, err
	}

	if success {
		taskLog := &model.TaskLog{
			TaskID:     task.ID,
			TaskName:   task.Name,
			Status:     model.TaskLogStatusRunning,
			RetryTimes: task.RetryTimes,
			WorkerID:   workerID,
			StartTime:  &time.Time{},
		}
		now := time.Now()
		taskLog.StartTime = &now
		err = s.logRepo.Create(taskLog)
		if err != nil {
			return false, err
		}
	}

	return success, nil
}

func (s *taskService) CompleteTask(taskID uint64, success bool, result, errorMsg string) error {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return err
	}

	var logStatus model.TaskLogStatus
	var taskStatus model.TaskStatus

	if success {
		logStatus = model.TaskLogStatusSuccess
		taskStatus = model.TaskStatusSuccess
	} else {
		logStatus = model.TaskLogStatusFailed
		taskStatus = model.TaskStatusFailed
	}

	now := time.Now()
	taskLog := &model.TaskLog{
		TaskID:     task.ID,
		TaskName:   task.Name,
		Status:     logStatus,
		RetryTimes: task.RetryTimes,
		WorkerID:   task.WorkerID,
		EndTime:    &now,
		Result:     result,
		ErrorMsg:   errorMsg,
	}
	taskLogs, _, err := s.logRepo.ListByTaskID(taskID, 1, 1)
	if err == nil && len(taskLogs) > 0 {
		latestLog := taskLogs[0]
		if latestLog.Status == model.TaskLogStatusRunning {
			latestLog.Status = logStatus
			latestLog.EndTime = &now
			latestLog.Result = result
			latestLog.ErrorMsg = errorMsg
			if latestLog.StartTime != nil {
				latestLog.Duration = int64(now.Sub(*latestLog.StartTime).Seconds())
			}
			err = s.logRepo.Update(&latestLog)
			if err != nil {
				return err
			}
		} else {
			err = s.logRepo.Create(taskLog)
			if err != nil {
				return err
			}
		}
	} else {
		err = s.logRepo.Create(taskLog)
		if err != nil {
			return err
		}
	}

	if success {
		if task.Type == model.TaskTypeCron && task.CronExpr != "" {
			nextTime, err := s.CalculateNextTime(task.CronExpr)
			if err != nil {
				return err
			}
			task.NextExecuteTime = nextTime
			task.Status = model.TaskStatusPending
			task.RetryTimes = 0
		} else {
			task.Status = taskStatus
		}
	} else {
		if task.RetryTimes < task.MaxRetryTimes {
			task.RetryTimes++
			nextTime := now.Add(time.Duration(task.RetryInterval) * time.Second)
			task.NextExecuteTime = &nextTime
			task.Status = model.TaskStatusPending
		} else {
			task.Status = taskStatus
		}
	}

	return s.taskRepo.Update(task)
}

func (s *taskService) CalculateNextTime(cronExpr string) (*time.Time, error) {
	schedule, err := cron.ParseStandard(cronExpr)
	if err != nil {
		return nil, err
	}
	next := schedule.Next(time.Now())
	return &next, nil
}
