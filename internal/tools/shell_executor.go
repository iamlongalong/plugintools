package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"gay/plugintools/internal/config"
	"gay/plugintools/internal/core"
)

// ShellExecutor Shell命令执行工具
type ShellExecutor struct{}

// NewShellExecutor 创建新的Shell执行工具实例
func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{}
}

// GetInfo 实现Tool接口
func (se *ShellExecutor) GetInfo() core.ToolInfo {
	return core.ToolInfo{
		ID:          "shell-executor",
		Name:        "Shell Executor",
		Description: "Execute shell commands with timeout and output capture",
		Version:     "1.0.0",
		Category:    "System",
	}
}

// GetParams 实现Tool接口
func (se *ShellExecutor) GetParams() []core.ParamSpec {
	return []core.ParamSpec{
		{
			Name:        "command",
			Type:        "string",
			Required:    true,
			Description: "Shell command to execute",
		},
		{
			Name:        "timeout",
			Type:        "integer",
			Required:    false,
			Default:     30,
			Description: "Command execution timeout in seconds",
		},
		{
			Name:        "working_dir",
			Type:        "string",
			Required:    false,
			Description: "Working directory for command execution",
		},
	}
}

// Execute 实现Tool接口
func (se *ShellExecutor) Execute(params map[string]interface{}) (interface{}, error) {
	command, ok := params["command"].(string)
	if !ok || command == "" {
		return nil, fmt.Errorf("command parameter is required")
	}

	// 验证命令是否在允许列表中
	if !se.isCommandAllowed(command) {
		return nil, fmt.Errorf("command not allowed: %s", command)
	}

	timeout := 30
	if t, ok := params["timeout"].(float64); ok {
		timeout = int(t)
	}

	// 检查超时限制
	cfg := config.Get()
	if timeout > cfg.Tools.ShellExecutor.MaxTimeout {
		return nil, fmt.Errorf("timeout exceeds maximum allowed value of %d seconds", cfg.Tools.ShellExecutor.MaxTimeout)
	}

	workingDir, _ := params["working_dir"].(string)

	// 创建命令
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	// 捕获输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 启动命令
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}

	// 设置超时
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// 等待命令完成或超时
	var err error
	select {
	case err = <-done:
	case <-time.After(time.Duration(timeout) * time.Second):
		if err = cmd.Process.Kill(); err != nil {
			return nil, fmt.Errorf("failed to kill process: %v", err)
		}
		return nil, fmt.Errorf("command timed out after %d seconds", timeout)
	}

	// 返回结果
	result := map[string]interface{}{
		"exit_code": cmd.ProcessState.ExitCode(),
		"stdout":    stdout.String(),
		"stderr":    stderr.String(),
		"success":   err == nil,
	}

	return result, nil
}

// isCommandAllowed 检查命令是否在允许列表中
func (se *ShellExecutor) isCommandAllowed(command string) bool {
	cfg := config.Get()
	cmdName := strings.Fields(command)[0]

	for _, allowed := range cfg.Tools.ShellExecutor.AllowedCommands {
		if cmdName == allowed {
			return true
		}
	}
	return false
}
