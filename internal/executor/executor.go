package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"task-scheduler/internal/model"
	"time"
)

type TaskExecutor interface {
	Execute(ctx context.Context, task *model.Task) (string, error)
}

type DefaultExecutor struct{}

func NewDefaultExecutor() *DefaultExecutor {
	return &DefaultExecutor{}
}

func (e *DefaultExecutor) Execute(ctx context.Context, task *model.Task) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	result := fmt.Sprintf("Task %d executed successfully at %s", task.ID, time.Now().Format(time.RFC3339))

	if task.Payload != "" {
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(task.Payload), &payload); err == nil {
			if action, ok := payload["action"].(string); ok {
				result += fmt.Sprintf(", action: %s", action)
			}
		}
	}

	time.Sleep(100 * time.Millisecond)

	return result, nil
}
