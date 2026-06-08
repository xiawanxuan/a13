package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"task-scheduler/internal/model"
	"time"
)

type CallbackData struct {
	TaskID     uint64 `json:"task_id"`
	TaskName   string `json:"task_name"`
	Status     int    `json:"status"`
	RetryTimes int    `json:"retry_times"`
	Result     string `json:"result,omitempty"`
	ErrorMsg   string `json:"error_msg,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
	Duration   int64  `json:"duration,omitempty"`
}

type CallbackService struct {
	client *http.Client
}

func NewCallbackService() *CallbackService {
	return &CallbackService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *CallbackService) ExecuteCallback(callbackURL string, taskLog *model.TaskLog, timeout int) error {
	if callbackURL == "" {
		return nil
	}

	if timeout <= 0 {
		timeout = 10
	}

	data := CallbackData{
		TaskID:     taskLog.TaskID,
		TaskName:   taskLog.TaskName,
		Status:     int(taskLog.Status),
		RetryTimes: taskLog.RetryTimes,
		Result:     taskLog.Result,
		ErrorMsg:   taskLog.ErrorMsg,
		Duration:   taskLog.Duration,
	}

	if taskLog.StartTime != nil {
		data.StartTime = taskLog.StartTime.Format(time.RFC3339)
	}
	if taskLog.EndTime != nil {
		data.EndTime = taskLog.EndTime.Format(time.RFC3339)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal callback data failed: %w", err)
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	req, err := http.NewRequest(http.MethodPost, callbackURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create callback request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("callback request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("callback returned non-2xx status: %d", resp.StatusCode)
	}

	return nil
}
