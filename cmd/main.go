package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-scheduler/internal/config"
	"task-scheduler/internal/executor"
	"task-scheduler/internal/handler"
	"task-scheduler/internal/pkg/database"
	"task-scheduler/internal/repository"
	"task-scheduler/internal/scheduler"
	"task-scheduler/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)

	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	fmt.Println("Database connected successfully")

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database migration completed")

	taskRepo := repository.NewTaskRepository(db)
	logRepo := repository.NewTaskLogRepository(db)

	taskService := service.NewTaskService(taskRepo, logRepo)
	taskLogService := service.NewTaskLogService(logRepo)

	taskHandler := handler.NewTaskHandler(taskService)
	taskLogHandler := handler.NewTaskLogHandler(taskLogService)

	taskExecutor := executor.NewDefaultExecutor()

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		taskHandler.RegisterRoutes(api)
		taskLogHandler.RegisterRoutes(api)
	}

	sched := scheduler.NewScheduler(taskService, taskExecutor, cfg.Scheduler.WorkerCount)

	if cfg.Scheduler.Enable {
		sched.Start()
		fmt.Println("Scheduler started")
	}

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		fmt.Printf("Server starting on %s\n", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down...")

	if cfg.Scheduler.Enable {
		sched.Stop()
	}

	fmt.Println("Server exited")
}
