package config

import (
	"encoding/json"
	"os"
	"sync"
)

// Config 系统配置结构
type Config struct {
	Server struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		ReadTimeout  int    `json:"read_timeout"`
		WriteTimeout int    `json:"write_timeout"`
	} `json:"server"`

	Security struct {
		APIKeys    []string `json:"api_keys"`
		EnableAuth bool     `json:"enable_auth"`
	} `json:"security"`

	Tools struct {
		FileManager struct {
			AllowedPaths []string `json:"allowed_paths"`
			MaxFileSize  int64    `json:"max_file_size"`
		} `json:"file_manager"`

		ShellExecutor struct {
			AllowedCommands []string `json:"allowed_commands"`
			MaxTimeout      int      `json:"max_timeout"`
		} `json:"shell_executor"`

		Scheduler struct {
			MaxTasks            int  `json:"max_tasks"`
			EnableNotifications bool `json:"enable_notifications"`
		} `json:"scheduler"`
	} `json:"tools"`
}

var (
	config *Config
	once   sync.Once
)

// Load 加载配置文件
func Load(path string) (*Config, error) {
	var err error
	once.Do(func() {
		config = &Config{}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			err = readErr
			return
		}
		if jsonErr := json.Unmarshal(data, config); jsonErr != nil {
			err = jsonErr
			return
		}
	})
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Get 获取配置实例
func Get() *Config {
	return config
}
