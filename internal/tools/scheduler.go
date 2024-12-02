package tools

import (
	"fmt"
	"sync"
	"time"

	"gay/plugintools/internal/config"
	"gay/plugintools/internal/core"
)

// Task 表示一个日程任务
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueTime     time.Time `json:"due_time"`
	Status      string    `json:"status"` // pending, completed, cancelled
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Scheduler 日程管理工具
type Scheduler struct {
	tasks map[string]*Task
	mu    sync.RWMutex
}

// NewScheduler 创建新的日程管理工具实例
func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make(map[string]*Task),
	}
}

// GetInfo 实现Tool接口
func (s *Scheduler) GetInfo() core.ToolInfo {
	return core.ToolInfo{
		ID:          "scheduler",
		Name:        "Task Scheduler",
		Description: "Manage tasks and schedules",
		Version:     "1.0.0",
		Category:    "Productivity",
	}
}

// GetParams 实现Tool接口
func (s *Scheduler) GetParams() []core.ParamSpec {
	return []core.ParamSpec{
		{
			Name:        "operation",
			Type:        "string",
			Required:    true,
			Description: "Operation to perform (create, update, delete, list, get)",
		},
		{
			Name:        "task_id",
			Type:        "string",
			Required:    false,
			Description: "Task ID for update, delete, get operations",
		},
		{
			Name:        "title",
			Type:        "string",
			Required:    false,
			Description: "Task title for create/update operations",
		},
		{
			Name:        "description",
			Type:        "string",
			Required:    false,
			Description: "Task description for create/update operations",
		},
		{
			Name:        "due_time",
			Type:        "string",
			Required:    false,
			Description: "Task due time in RFC3339 format",
		},
		{
			Name:        "status",
			Type:        "string",
			Required:    false,
			Description: "Task status (pending, completed, cancelled)",
		},
	}
}

// Execute 实现Tool接口
func (s *Scheduler) Execute(params map[string]interface{}) (interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter is required")
	}

	switch operation {
	case "create":
		return s.createTask(params)
	case "update":
		return s.updateTask(params)
	case "delete":
		return s.deleteTask(params)
	case "list":
		return s.listTasks()
	case "get":
		return s.getTask(params)
	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}
}

// createTask 创建新任务
func (s *Scheduler) createTask(params map[string]interface{}) (*Task, error) {
	// 检查任务数量限制
	s.mu.RLock()
	cfg := config.Get()
	if len(s.tasks) >= cfg.Tools.Scheduler.MaxTasks {
		s.mu.RUnlock()
		return nil, fmt.Errorf("maximum number of tasks (%d) reached", cfg.Tools.Scheduler.MaxTasks)
	}
	s.mu.RUnlock()

	title, ok := params["title"].(string)
	if !ok || title == "" {
		return nil, fmt.Errorf("title is required for create operation")
	}

	task := &Task{
		ID:          fmt.Sprintf("task_%d", time.Now().UnixNano()),
		Title:       title,
		Description: params["description"].(string),
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if dueTime, ok := params["due_time"].(string); ok {
		t, err := time.Parse(time.RFC3339, dueTime)
		if err != nil {
			return nil, fmt.Errorf("invalid due_time format: %v", err)
		}
		task.DueTime = t
	}

	s.mu.Lock()
	s.tasks[task.ID] = task
	s.mu.Unlock()

	// 如果启用了通知，发送创建通知
	if cfg.Tools.Scheduler.EnableNotifications {
		go s.sendNotification("task_created", task)
	}

	return task, nil
}

// updateTask 更新任务
func (s *Scheduler) updateTask(params map[string]interface{}) (*Task, error) {
	taskID, ok := params["task_id"].(string)
	if !ok || taskID == "" {
		return nil, fmt.Errorf("task_id is required for update operation")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	if title, ok := params["title"].(string); ok && title != "" {
		task.Title = title
	}
	if desc, ok := params["description"].(string); ok {
		task.Description = desc
	}
	if status, ok := params["status"].(string); ok && status != "" {
		task.Status = status
	}
	if dueTime, ok := params["due_time"].(string); ok {
		t, err := time.Parse(time.RFC3339, dueTime)
		if err != nil {
			return nil, fmt.Errorf("invalid due_time format: %v", err)
		}
		task.DueTime = t
	}

	task.UpdatedAt = time.Now()

	// 如果启用了通知，发送更新通知
	if config.Get().Tools.Scheduler.EnableNotifications {
		go s.sendNotification("task_updated", task)
	}

	return task, nil
}

// deleteTask 删除任务
func (s *Scheduler) deleteTask(params map[string]interface{}) (interface{}, error) {
	taskID, ok := params["task_id"].(string)
	if !ok || taskID == "" {
		return nil, fmt.Errorf("task_id is required for delete operation")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	delete(s.tasks, taskID)

	// 如果启用了通知，发送删除通知
	if config.Get().Tools.Scheduler.EnableNotifications {
		go s.sendNotification("task_deleted", task)
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Task %s deleted", taskID),
	}, nil
}

// listTasks 列出所有任务
func (s *Scheduler) listTasks() (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// getTask 获取单个任务
func (s *Scheduler) getTask(params map[string]interface{}) (*Task, error) {
	taskID, ok := params["task_id"].(string)
	if !ok || taskID == "" {
		return nil, fmt.Errorf("task_id is required for get operation")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	return task, nil
}

// sendNotification 发送任务通知
func (s *Scheduler) sendNotification(eventType string, task *Task) {
	// TODO: 实现实际的通知逻辑
	// 这里可以集成邮件、WebSocket、webhook等通知方式
	fmt.Printf("Notification: %s - Task %s (%s)\n", eventType, task.ID, task.Title)
}
