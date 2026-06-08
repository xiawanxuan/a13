package model

import (
	"time"
)

type TaskStatus int

const (
	TaskStatusPending   TaskStatus = iota
	TaskStatusRunning
	TaskStatusSuccess
	TaskStatusFailed
	TaskStatusPaused
)

type TaskType int

const (
	TaskTypeCron    TaskType = iota
	TaskTypeOneTime
	TaskTypeDependent
)

type Task struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"size:128;not null;index" json:"name"`
	Type            TaskType   `gorm:"not null;default:0" json:"type"`
	CronExpr        string     `gorm:"size:64" json:"cron_expr"`
	Payload         string     `gorm:"type:text" json:"payload"`
	MaxRetryTimes   int        `gorm:"not null;default:3" json:"max_retry_times"`
	RetryInterval   int        `gorm:"not null;default:60" json:"retry_interval"`
	Timeout         int        `gorm:"not null;default:300" json:"timeout"`
	CallbackURL     string     `gorm:"size:512" json:"callback_url"`
	CallbackTimeout int        `gorm:"not null;default:10" json:"callback_timeout"`
	Status          TaskStatus `gorm:"not null;default:0;index" json:"status"`
	RetryTimes      int        `gorm:"not null;default:0" json:"retry_times"`
	NextExecuteTime *time.Time `gorm:"index" json:"next_execute_time"`
	LastExecuteTime *time.Time `json:"last_execute_time"`
	WorkerID        string     `gorm:"size:64;index" json:"worker_id"`
	Version         int        `gorm:"not null;default:0" json:"-"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (Task) TableName() string {
	return "tasks"
}
