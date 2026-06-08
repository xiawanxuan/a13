package dto

type CreateTaskRequest struct {
	Name            string   `json:"name" binding:"required"`
	Type            int      `json:"type"`
	CronExpr        string   `json:"cron_expr"`
	Payload         string   `json:"payload"`
	MaxRetryTimes   int      `json:"max_retry_times"`
	RetryInterval   int      `json:"retry_interval"`
	Timeout         int      `json:"timeout"`
	CallbackURL     string   `json:"callback_url"`
	CallbackTimeout int      `json:"callback_timeout"`
	UpstreamTaskIDs []uint64 `json:"upstream_task_ids"`
}

type UpdateTaskRequest struct {
	Name            string  `json:"name"`
	CronExpr        string  `json:"cron_expr"`
	Payload         string  `json:"payload"`
	MaxRetryTimes   *int    `json:"max_retry_times"`
	RetryInterval   *int    `json:"retry_interval"`
	Timeout         *int    `json:"timeout"`
	CallbackURL     string  `json:"callback_url"`
	CallbackTimeout *int    `json:"callback_timeout"`
}

type TaskListRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
	Status   *int `form:"status"`
}

type TaskLogListRequest struct {
	TaskID   uint64 `form:"task_id" binding:"required"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type AddDependencyRequest struct {
	TaskID     uint64 `json:"task_id" binding:"required"`
	UpstreamID uint64 `json:"upstream_id" binding:"required"`
}

type TaskMetrics struct {
	TotalTasks     int64            `json:"total_tasks"`
	StatusBreakdown map[string]int64 `json:"status_breakdown"`
	TotalExecutions int64           `json:"total_executions"`
	SuccessCount    int64           `json:"success_count"`
	FailedCount     int64           `json:"failed_count"`
	RunningCount    int64           `json:"running_count"`
	TodayExecutions int64           `json:"today_executions"`
	SuccessRate     float64         `json:"success_rate"`
	AvgDuration     float64         `json:"avg_duration_seconds"`
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

