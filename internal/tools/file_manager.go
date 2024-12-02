package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gay/plugintools/internal/config"
	"gay/plugintools/internal/core"
)

// FileManager 文件管理工具
type FileManager struct{}

// NewFileManager 创建新的文件管理工具实例
func NewFileManager() *FileManager {
	return &FileManager{}
}

// GetInfo 实现Tool接口
func (fm *FileManager) GetInfo() core.ToolInfo {
	return core.ToolInfo{
		ID:          "file-manager",
		Name:        "File Manager",
		Description: "Provides file system operations like list, copy, move, delete",
		Version:     "1.0.0",
		Category:    "System",
	}
}

// GetParams 实现Tool接口
func (fm *FileManager) GetParams() []core.ParamSpec {
	return []core.ParamSpec{
		{
			Name:        "operation",
			Type:        "string",
			Required:    true,
			Description: "Operation to perform (list, copy, move, delete)",
		},
		{
			Name:        "path",
			Type:        "string",
			Required:    true,
			Description: "File or directory path",
		},
		{
			Name:        "destination",
			Type:        "string",
			Required:    false,
			Description: "Destination path for copy/move operations",
		},
	}
}

// Execute 实现Tool接口
func (fm *FileManager) Execute(params map[string]interface{}) (interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter is required")
	}

	path, ok := params["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path parameter is required")
	}

	// 验证路径是否在允许的范围内
	if !fm.isPathAllowed(path) {
		return nil, fmt.Errorf("access to path %s is not allowed", path)
	}

	switch operation {
	case "list":
		return fm.list(path)
	case "delete":
		return nil, fm.delete(path)
	case "copy", "move":
		dest, ok := params["destination"].(string)
		if !ok {
			return nil, fmt.Errorf("destination parameter is required for copy/move operations")
		}
		if !fm.isPathAllowed(dest) {
			return nil, fmt.Errorf("access to destination path %s is not allowed", dest)
		}
		if operation == "copy" {
			return nil, fm.copy(path, dest)
		}
		return nil, fm.move(path, dest)
	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}
}

// isPathAllowed 检查路径是否在允许的范围内
func (fm *FileManager) isPathAllowed(path string) bool {
	cfg := config.Get()
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	for _, allowedPath := range cfg.Tools.FileManager.AllowedPaths {
		allowedAbs, err := filepath.Abs(allowedPath)
		if err != nil {
			continue
		}
		if isSubPath(allowedAbs, absPath) {
			return true
		}
	}
	return false
}

// isSubPath 检查childPath是否是parentPath的子路径
func isSubPath(parentPath, childPath string) bool {
	rel, err := filepath.Rel(parentPath, childPath)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}

// list 列出目录内容
func (fm *FileManager) list(path string) ([]map[string]interface{}, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"name":    entry.Name(),
			"size":    info.Size(),
			"mode":    info.Mode().String(),
			"modTime": info.ModTime(),
			"isDir":   entry.IsDir(),
		})
	}

	return result, nil
}

// delete 删除文件或目录
func (fm *FileManager) delete(path string) error {
	return os.RemoveAll(path)
}

// copy 复制文件或目录
func (fm *FileManager) copy(src, dst string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 检查文件大小限制
	if !sourceInfo.IsDir() && sourceInfo.Size() > config.Get().Tools.FileManager.MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", config.Get().Tools.FileManager.MaxFileSize)
	}

	if sourceInfo.IsDir() {
		return fm.copyDir(src, dst)
	}
	return fm.copyFile(src, dst)
}

// copyFile 复制单个文件
func (fm *FileManager) copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// copyDir 复制目录
func (fm *FileManager) copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := fm.copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := fm.copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// move 移动文件或目录
func (fm *FileManager) move(src, dst string) error {
	return os.Rename(src, dst)
}
