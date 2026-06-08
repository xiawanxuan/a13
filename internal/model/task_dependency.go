package model

import "time"

type TaskDependency struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID     uint64    `gorm:"not null;index:idx_task_id" json:"task_id"`
	UpstreamID uint64    `gorm:"not null;index:idx_upstream_id" json:"upstream_id"`
	Status     int       `gorm:"not null;default:0" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (TaskDependency) TableName() string {
	return "task_dependencies"
}
