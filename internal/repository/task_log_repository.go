package repository

import (
	"task-scheduler/internal/model"

	"gorm.io/gorm"
)

type TaskLogRepository interface {
	Create(log *model.TaskLog) error
	GetByID(id uint64) (*model.TaskLog, error)
	ListByTaskID(taskID uint64, page, pageSize int) ([]model.TaskLog, int64, error)
	Update(log *model.TaskLog) error
	CountByStatus(status model.TaskLogStatus) (int64, error)
	CountAll() (int64, error)
	CountToday() (int64, error)
	AvgDuration() (float64, error)
}

type taskLogRepository struct {
	db *gorm.DB
}

func NewTaskLogRepository(db *gorm.DB) TaskLogRepository {
	return &taskLogRepository{db: db}
}

func (r *taskLogRepository) Create(log *model.TaskLog) error {
	return r.db.Create(log).Error
}

func (r *taskLogRepository) GetByID(id uint64) (*model.TaskLog, error) {
	var log model.TaskLog
	err := r.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *taskLogRepository) ListByTaskID(taskID uint64, page, pageSize int) ([]model.TaskLog, int64, error) {
	var logs []model.TaskLog
	var total int64

	query := r.db.Model(&model.TaskLog{}).Where("task_id = ?", taskID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *taskLogRepository) Update(log *model.TaskLog) error {
	return r.db.Save(log).Error
}
