# Tools Platform

一个基于 OpenAPI 的工具平台，提供统一的工具调用接口。

## 功能特点

- 统一的工具接口抽象
- RESTful API
- 内置多种工具
- 可扩展的插件系统
- 认证和日志支持

## 内置工具

1. 文件管理工具 (file-manager)
   - 列出目录内容
   - 复制文件/目录
   - 移动文件/目录
   - 删除文件/目录

2. Shell命令执行工具 (shell-executor)
   - 执行shell命令
   - 超时控制
   - 输出捕获
   - 工作目录设置

3. 日程管理工具 (scheduler)
   - 创建/更新任务
   - 删除任务
   - 列出任务
   - 获取任务详情

## 快速开始

1. 安装
```bash
go mod download
```

2. 运行服务器
```bash
go run cmd/server/main.go
```

## API 使用示例

1. 获取所有工具列表
```bash
curl -H "X-API-Key: test-api-key" http://localhost:8080/api/v1/tools
```

2. 获取工具信息
```bash
curl -H "X-API-Key: test-api-key" http://localhost:8080/api/v1/tools/file-manager
```

3. 获取工具参数定义
```bash
curl -H "X-API-Key: test-api-key" http://localhost:8080/api/v1/tools/file-manager?params=true
```

4. 执行工具
```bash
# 列出目录内容
curl -X POST -H "X-API-Key: test-api-key" -H "Content-Type: application/json" \
     -d '{"operation":"list","path":"/tmp"}' \
     http://localhost:8080/api/v1/tools/file-manager

# 执行Shell命令
curl -X POST -H "X-API-Key: test-api-key" -H "Content-Type: application/json" \
     -d '{"command":"ls -l","timeout":30}' \
     http://localhost:8080/api/v1/tools/shell-executor

# 创建任务
curl -X POST -H "X-API-Key: test-api-key" -H "Content-Type: application/json" \
     -d '{"operation":"create","title":"测试任务","description":"这是一个测试任务","due_time":"2024-12-31T23:59:59Z"}' \
     http://localhost:8080/api/v1/tools/scheduler
```

## 配置说明

配置文件位于 `configs/config.json`，包含以下主要配置项：

- 服务器配置（地址、端口、超时等）
- 安全配置（API密钥、认证开关）
- 工具配置（各工具的特定配置）

## 添加新工具

1. 在 `internal/tools` 目录下创建新的工具实现
2. 实现 `Tool` 接口的所有方法：
   - `GetInfo()`
   - `GetParams()`
   - `Execute()`
3. 在 `cmd/server/main.go` 中注册新工具

## 安全性说明

- 所有API调用需要提供有效的API密钥
- 文件操作限制在允许的路径内
- Shell命令限制在允许的命令列表内
- 所有操作都有日志记录

## 开发计划

- [ ] 添加更多工具
- [ ] 实现工具版本管理
- [ ] 添加WebSocket支持
- [ ] 实现异步任务
- [ ] 添加更多安全特性
- [ ] 实现工具市场

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License 