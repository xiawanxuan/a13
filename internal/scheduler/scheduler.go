package scheduler

import (
	"context"
	"fmt"
	"math/rand"
	"task-scheduler/internal/executor"
	"task-scheduler/internal/model"
	"task-scheduler/internal/service"
	"time"
)

type Scheduler struct {
	taskService  service.TaskService
	executor     executor.TaskExecutor
	workerID     string
	workerCount  int
	taskQueue    chan *model.Task
	pollInterval time.Duration
	maxFetchSize int
	running      bool
	stopChan     chan struct{}
}

func NewScheduler(taskService service.TaskService, exec executor.TaskExecutor, workerCount int) *Scheduler {
	if workerCount <= 0 {
		workerCount = 5
	}

	workerID := generateWorkerID()

	return &Scheduler{
		taskService:  taskService,
		executor:     exec,
		workerID:     workerID,
		workerCount:  workerCount,
		taskQueue:    make(chan *model.Task, workerCount*2),
		pollInterval: 5 * time.Second,
		maxFetchSize: 10,
		stopChan:     make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	if s.running {
		return
	}
	s.running = true

	fmt.Printf("[Scheduler] Starting scheduler with worker_id=%s, worker_count=%d\n", s.workerID, s.workerCount)

	for i := 0; i < s.workerCount; i++ {
		go s.worker(i)
	}

	go s.pollTasks()

	fmt.Println("[Scheduler] Scheduler started successfully")
}

func (s *Scheduler) Stop() {
	if !s.running {
		return
	}
	s.running = false
	close(s.stopChan)
	fmt.Println("[Scheduler] Scheduler stopped")
}

func (s *Scheduler) pollTasks() {
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.fetchAndDispatchTasks()
		}
	}
}

func (s *Scheduler) fetchAndDispatchTasks() {
	tasks, err := s.taskService.GetDueTasks(s.maxFetchSize)
	if err != nil {
		fmt.Printf("[Scheduler] Failed to get due tasks: %v\n", err)
		return
	}

	for i := range tasks {
		task := &tasks[i]
		claimed, err := s.taskService.ClaimTask(task.ID, s.workerID)
		if err != nil {
			fmt.Printf("[Scheduler] Failed to claim task %d: %v\n", task.ID, err)
			continue
		}
		if claimed {
			fmt.Printf("[Scheduler] Task %d claimed by worker %s\n", task.ID, s.workerID)
			select {
			case s.taskQueue <- task:
			default:
				fmt.Printf("[Scheduler] Task queue is full, skipping task %d\n", task.ID)
			}
		}
	}
}

func (s *Scheduler) worker(id int) {
	fmt.Printf("[Worker-%d] Worker started\n", id)

	for task := range s.taskQueue {
		s.executeTask(id, task)
	}

	fmt.Printf("[Worker-%d] Worker stopped\n", id)
}

func (s *Scheduler) executeTask(workerID int, task *model.Task) {
	fmt.Printf("[Worker-%d] Executing task %d (%s)\n", workerID, task.ID, task.Name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(task.Timeout)*time.Second)
	defer cancel()

	result, err := s.executor.Execute(ctx, task)

	var success bool
	var errorMsg string

	if err != nil {
		success = false
		errorMsg = err.Error()
		fmt.Printf("[Worker-%d] Task %d failed: %v\n", workerID, task.ID, err)
	} else {
		success = true
		fmt.Printf("[Worker-%d] Task %d completed successfully\n", workerID, task.ID)
	}

	completeErr := s.taskService.CompleteTask(task.ID, success, result, errorMsg)
	if completeErr != nil {
		fmt.Printf("[Worker-%d] Failed to complete task %d: %v\n", workerID, task.ID, completeErr)
	}
}

func (s *Scheduler) GetWorkerID() string {
	return s.workerID
}

func generateWorkerID() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 8)
	for i := range b {
		b[i] = byte(rand.Intn(16))
	}
	return fmt.Sprintf("%x", b)
}
