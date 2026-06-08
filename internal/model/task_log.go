package model

import "time"

type TaskLogStatus int

const (
	TaskLogStatusRunning TaskLogStatus = iota
	TaskLogStatusSuccess
	TaskLogStatusFailed
)

type TaskLog struct {
	ID         uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID     uint64        `gorm:"not null;index:idx_task_id" json:"task_id"`
	TaskName   string        `gorm:"size:128;not null" json:"task_name"`
	Status     TaskLogStatus `gorm:"not null;index" json:"status"`
	RetryTimes int           `gorm:"not null;default:0" json:"retry_times"`
	WorkerID   string        `gorm:"size:64;index" json:"worker_id"`
	StartTime  *time.Time    `json:"start_time"`
	EndTime    *time.Time    `json:"end_time"`
	Duration   int64         `json:"duration"`
	Result     string        `gorm:"type:text" json:"result"`
	ErrorMsg   string        `gorm:"type:text" json:"error_msg"`
	CreatedAt  time.Time     `json:"created_at"`
}

func (TaskLog) TableName() string {
	return "task_logs"
}
