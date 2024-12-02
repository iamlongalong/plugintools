package core

// Tool 定义了统一的工具接口
type Tool interface {
	// GetInfo 返回工具的基本信息
	GetInfo() ToolInfo
	// Execute 执行工具，接收参数并返回结果
	Execute(params map[string]interface{}) (interface{}, error)
	// GetParams 返回工具的参数定义
	GetParams() []ParamSpec
}

// ToolInfo 包含工具的基本信息
type ToolInfo struct {
	ID          string `json:"id"`          // 工具唯一标识
	Name        string `json:"name"`        // 工具名称
	Description string `json:"description"` // 工具描述
	Version     string `json:"version"`     // 工具版本
	Category    string `json:"category"`    // 工具分类
}

// ParamSpec 定义了工具参数的规格
type ParamSpec struct {
	Name        string      `json:"name"`        // 参数名
	Type        string      `json:"type"`        // 参数类型
	Required    bool        `json:"required"`    // 是否必需
	Default     interface{} `json:"default"`     // 默认值
	Description string      `json:"description"` // 参数描述
}

// ToolRegistry 工具注册表接口
type ToolRegistry interface {
	// Register 注册一个新工具
	Register(tool Tool) error
	// Get 获取指定ID的工具
	Get(id string) (Tool, error)
	// List 列出所有已注册的工具
	List() []Tool
	// Unregister 注销一个工具
	Unregister(id string) error
}
