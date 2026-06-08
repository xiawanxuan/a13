package repository

import (
	"task-scheduler/internal/model"
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *model.Task) error
	GetByID(id uint64) (*model.Task, error)
	Update(task *model.Task) error
	Delete(id uint64) error
	List(page, pageSize int, status *model.TaskStatus) ([]model.Task, int64, error)
	GetPendingTasks(limit int) ([]model.Task, error)
	GetDueTasks(now time.Time, limit int) ([]model.Task, error)
	UpdateStatus(id uint64, status model.TaskStatus, version int) (bool, error)
	ClaimTask(id uint64, workerID string, version int) (bool, error)
	UpdateWithRetry(task *model.Task, version int) (bool, error)
	ClaimDueTasks(workerID string, limit int) ([]model.Task, error)
	RecoverTimeoutTasks(timeoutSeconds int) (int64, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) GetByID(id uint64) (*model.Task, error) {
	var task model.Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) Update(task *model.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Task{}, id).Error
}

func (r *taskRepository) List(page, pageSize int, status *model.TaskStatus) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := r.db.Model(&model.Task{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) GetPendingTasks(limit int) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where("status = ?", model.TaskStatusPending).
		Order("id ASC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetDueTasks(now time.Time, limit int) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where("status IN ? AND next_execute_time <= ?",
		[]model.TaskStatus{model.TaskStatusPending, model.TaskStatusFailed},
		now).
		Order("next_execute_time ASC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateStatus(id uint64, status model.TaskStatus, version int) (bool, error) {
	result := r.db.Model(&model.Task{}).
		Where("id = ? AND version = ?", id, version).
		Updates(map[string]interface{}{
			"status":  status,
			"version": version + 1,
		})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *taskRepository) ClaimTask(id uint64, workerID string, version int) (bool, error) {
	result := r.db.Model(&model.Task{}).
		Where("id = ? AND version = ?", id, version).
		Updates(map[string]interface{}{
			"status":            model.TaskStatusRunning,
			"worker_id":         workerID,
			"last_execute_time": time.Now(),
			"version":           version + 1,
		})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *taskRepository) UpdateWithRetry(task *model.Task, version int) (bool, error) {
	result := r.db.Model(&model.Task{}).
		Where("id = ? AND version = ?", task.ID, version).
		Updates(map[string]interface{}{
			"status":            task.Status,
			"retry_times":       task.RetryTimes,
			"next_execute_time": task.NextExecuteTime,
			"worker_id":         task.WorkerID,
			"version":           version + 1,
		})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *taskRepository) ClaimDueTasks(workerID string, limit int) ([]model.Task, error) {
	var claimedTasks []model.Task

	err := r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		startTime := now.Add(-1 * time.Second)

		result := tx.Model(&model.Task{}).
			Where("status IN ? AND next_execute_time <= ?",
				[]model.TaskStatus{model.TaskStatusPending, model.TaskStatusFailed},
				now).
			Order("next_execute_time ASC").
			Limit(limit).
			Updates(map[string]interface{}{
				"status":            model.TaskStatusRunning,
				"worker_id":         workerID,
				"last_execute_time": now,
				"version":           gorm.Expr("version + 1"),
			})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return nil
		}

		err := tx.Where("worker_id = ? AND status = ? AND last_execute_time >= ?",
			workerID, model.TaskStatusRunning, startTime).
			Order("next_execute_time ASC").
			Find(&claimedTasks).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return claimedTasks, nil
}

func (r *taskRepository) RecoverTimeoutTasks(timeoutSeconds int) (int64, error) {
	now := time.Now()
	timeoutThreshold := now.Add(-time.Duration(timeoutSeconds) * time.Second)

	result := r.db.Model(&model.Task{}).
		Where("status = ? AND last_execute_time <= ?",
			model.TaskStatusRunning,
			timeoutThreshold).
		Updates(map[string]interface{}{
			"status":  model.TaskStatusFailed,
			"version": gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
