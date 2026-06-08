package repository

import (
	"task-scheduler/internal/model"

	"gorm.io/gorm"
)

type TaskDependencyRepository interface {
	Create(dep *model.TaskDependency) error
	Delete(id uint64) error
	GetByTaskID(taskID uint64) ([]model.TaskDependency, error)
	GetByUpstreamID(upstreamID uint64) ([]model.TaskDependency, error)
	GetDownstreamTasks(upstreamID uint64) ([]model.Task, error)
	CheckAllUpstreamCompleted(taskID uint64) (bool, error)
	UpdateDependencyStatus(upstreamID uint64, status int) (int64, error)
	DeleteByTaskID(taskID uint64) error
	DeleteByUpstreamID(upstreamID uint64) error
}

type taskDependencyRepository struct {
	db *gorm.DB
}

func NewTaskDependencyRepository(db *gorm.DB) TaskDependencyRepository {
	return &taskDependencyRepository{db: db}
}

func (r *taskDependencyRepository) Create(dep *model.TaskDependency) error {
	return r.db.Create(dep).Error
}

func (r *taskDependencyRepository) Delete(id uint64) error {
	return r.db.Delete(&model.TaskDependency{}, id).Error
}

func (r *taskDependencyRepository) GetByTaskID(taskID uint64) ([]model.TaskDependency, error) {
	var deps []model.TaskDependency
	err := r.db.Where("task_id = ?", taskID).Find(&deps).Error
	return deps, err
}

func (r *taskDependencyRepository) GetByUpstreamID(upstreamID uint64) ([]model.TaskDependency, error) {
	var deps []model.TaskDependency
	err := r.db.Where("upstream_id = ?", upstreamID).Find(&deps).Error
	return deps, err
}

func (r *taskDependencyRepository) GetDownstreamTasks(upstreamID uint64) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Table("tasks").
		Select("tasks.*").
		Joins("JOIN task_dependencies td ON tasks.id = td.task_id").
		Where("td.upstream_id = ?", upstreamID).
		Scan(&tasks).Error
	return tasks, err
}

func (r *taskDependencyRepository) CheckAllUpstreamCompleted(taskID uint64) (bool, error) {
	var total int64
	err := r.db.Model(&model.TaskDependency{}).
		Where("task_id = ?", taskID).
		Count(&total).Error
	if err != nil {
		return false, err
	}

	if total == 0 {
		return true, nil
	}

	var completed int64
	err = r.db.Model(&model.TaskDependency{}).
		Where("task_id = ? AND status = ?", taskID, 1).
		Count(&completed).Error
	if err != nil {
		return false, err
	}

	return total == completed, nil
}

func (r *taskDependencyRepository) UpdateDependencyStatus(upstreamID uint64, status int) (int64, error) {
	result := r.db.Model(&model.TaskDependency{}).
		Where("upstream_id = ?", upstreamID).
		Update("status", status)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *taskDependencyRepository) DeleteByTaskID(taskID uint64) error {
	return r.db.Where("task_id = ?", taskID).Delete(&model.TaskDependency{}).Error
}

func (r *taskDependencyRepository) DeleteByUpstreamID(upstreamID uint64) error {
	return r.db.Where("upstream_id = ?", upstreamID).Delete(&model.TaskDependency{}).Error
}
