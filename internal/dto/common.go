package dto

type CreateTaskRequest struct {
	Name          string `json:"name" binding:"required"`
	Type          int    `json:"type"`
	CronExpr      string `json:"cron_expr"`
	Payload       string `json:"payload"`
	MaxRetryTimes int    `json:"max_retry_times"`
	RetryInterval int    `json:"retry_interval"`
	Timeout       int    `json:"timeout"`
}

type UpdateTaskRequest struct {
	Name          string `json:"name"`
	CronExpr      string `json:"cron_expr"`
	Payload       string `json:"payload"`
	MaxRetryTimes *int   `json:"max_retry_times"`
	RetryInterval *int   `json:"retry_interval"`
	Timeout       *int   `json:"timeout"`
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
