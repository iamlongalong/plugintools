package core

import (
	"fmt"
	"sync"
)

// DefaultRegistry 默认的工具注册表实现
type DefaultRegistry struct {
	tools map[string]Tool
	mu    sync.RWMutex
}

// NewRegistry 创建一个新的工具注册表
func NewRegistry() *DefaultRegistry {
	return &DefaultRegistry{
		tools: make(map[string]Tool),
	}
}

// Register 注册一个新工具
func (r *DefaultRegistry) Register(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info := tool.GetInfo()
	if info.ID == "" {
		return fmt.Errorf("tool ID cannot be empty")
	}

	if _, exists := r.tools[info.ID]; exists {
		return fmt.Errorf("tool with ID %s already exists", info.ID)
	}

	r.tools[info.ID] = tool
	return nil
}

// Get 获取指定ID的工具
func (r *DefaultRegistry) Get(id string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[id]
	if !exists {
		return nil, fmt.Errorf("tool with ID %s not found", id)
	}

	return tool, nil
}

// List 列出所有已注册的工具
func (r *DefaultRegistry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}

	return tools
}

// Unregister 注销一个工具
func (r *DefaultRegistry) Unregister(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[id]; !exists {
		return fmt.Errorf("tool with ID %s not found", id)
	}

	delete(r.tools, id)
	return nil
}
