<center>

# AkashGen API Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![API](https://img.shields.io/badge/API-REST-orange.svg)](https://restfulapi.net/)

一个基于 Go 语言的高性能 API 代理服务，用于访问 Akash Network 的图像生成功能。该服务提供了一个简洁的 REST API 接口，方便与 Akash Network 的去中心化 AI 图像生成基础设施进行交互。

[中文文档](README_CN.md) | [English](README.md)

</center>

## 特性

- 🚀 **高性能**: 使用 Go 和 Gin 框架构建，性能优异
- 🔄 **异步处理**: 非阻塞的任务提交和状态轮询机制
- 🛡️ **并发控制**: 内置速率限制防止资源耗尽
- 📊 **结构化日志**: 使用 Zap 记录详细日志，便于生产环境监控
- 🏗️ **清晰架构**: 采用模块化设计，遵循 Go 最佳实践
- 🌐 **REST API**: 简单直观的 HTTP 端点
- ⚡ **优雅关闭**: 正确的清理和连接处理
- 🎯 **GPU 偏好**: 自动从可用选项中选择 GPU

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- 能够访问 Akash Network 的网络连接

### 安装

1. **克隆仓库**
   ```bash
   git clone https://github.com/006lp/akashgen-api-go.git
   cd akashgen-api-go
   ```

2. **初始化 Go 模块**
   ```bash
   go mod init akashgen-api-go
   go mod tidy
   ```

3. **构建并运行**
   ```bash
   go build -o akashgen-api .
   ./akashgen-api
   ```

或直接运行:
```bash
go run main.go
```

服务默认会在端口 `6571` 上启动。

## API 文档

### 生成图像

使用 Akash Network 的去中心化基础设施生成 AI 图像。

**端点:** `POST /api/generate`

**请求体:**
```json
{
    "prompt": "山脉上美丽的日落",
    "negative": "模糊，低质量",
    "sampler": "dpmpp_2m",
    "scheduler": "karras"
}
```

> 注意：提示词建议使用英语描述，识别效果更佳。

**参数说明:**
- `prompt` (必需): 对期望图像的文本描述
- `negative` (可选): 描述图像中要避免的内容
- `sampler` (必需): 使用的采样方法（详见下方支持的采样器）
- `scheduler` (必需): 调度算法（详见下方支持的调度器）

**支持的采样器 (Samplers):**
```
euler, euler_cfg_pp, euler_ancestral, euler_ancestral_cfg_pp, heun, heunpp2, 
dpm_2, dpm_2_ancestral, lms, dpm_fast, dpm_adaptive, dpmpp_2s_ancestral, 
dpmpp_2s_ancestral_cfg_pp, dpmpp_sde, dpmpp_sde_gpu, dpmpp_2m, dpmpp_2m_cfg_pp, 
dpmpp_2m_sde, dpmpp_2m_sde_gpu, dpmpp_3m_sde, dpmpp_3m_sde_gpu, ddpm, lcm, 
ipndm, ipndm_v, deis, ddim, uni_pc, uni_pc_bh2
```

**支持的调度器 (Schedulers):**
```
normal, karras, exponential, sgm_uniform, simple, ddim_uniform, beta, linear_quadratic
```

**响应:**
- **成功 (200)**: 返回生成的图像二进制数据
- **错误 (400)**: 无效的请求格式
- **错误 (500)**: 生成失败或超时

**使用 curl 的示例:**
```bash
curl -X POST http://localhost:6571/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "宁静的湖泊，背景是山脉",
    "negative": "丑陋，模糊，低分辨率",
    "sampler": "dpmpp_2m",
    "scheduler": "karras"
  }' \
  --output generated_image.png
```

### 健康检查

检查服务是否正常运行。

**端点:** `GET /health`

**响应:**
```json
{
    "status": "ok"
}
```

## 配置

可以通过修改 `config/config.go` 中的常量来配置服务：

```go
const (
    // 超时配置
    GenerateTimeout    = 30 * time.Second  // 生成请求最大时间
    StatusTimeout      = 10 * time.Second  // 状态检查最大时间
    ImageFetchTimeout  = 30 * time.Second  // 图像下载最大时间
    PollingInterval    = 1 * time.Second   // 状态检查间隔
    MaxPollingDuration = 5 * time.Minute   // 等待完成的最大时间
    
    // 并发控制
    MaxConcurrentRequests = 10  // 最大同时请求数
    
    // 服务器配置
    ServerPort = ":6571"  // 监听端口
)
```

## 项目架构

```
akashgen-api-go/
├── main.go              # 应用程序入口点
├── config/
│   └── config.go       # 配置常量
├── handlers/
│   └── generate.go     # HTTP 请求处理器
├── models/
│   └── types.go        # 数据结构和类型
├── services/
│   ├── akash.go        # Akash Network API 客户端
│   └── image.go        # 图像获取服务
├── middleware/
│   └── logger.go       # HTTP 日志中间件
└── utils/
    └── http.go         # HTTP 工具函数
```

### 组件说明

- **Handlers**: 处理 HTTP 请求和响应
- **Services**: 与外部 API 交互的业务逻辑
- **Models**: 请求和响应的数据结构
- **Middleware**: 横切关注点，如日志记录
- **Config**: 集中的配置管理
- **Utils**: 可重用的工具函数

## 开发

### 项目结构

该项目遵循 Go 最佳实践，关注点清晰分离：

- 依赖注入的清晰架构
- 模块化设计，便于测试和维护
- 适当的错误处理和日志记录
- 基于上下文的超时管理

### 生产构建

```bash
# 构建优化的二进制文件
go build -ldflags="-w -s" -o akashgen-api .

# 或为不同平台构建
GOOS=linux GOARCH=amd64 go build -o akashgen-api-linux .
```

### Docker 支持

创建 `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o akashgen-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/akashgen-api .
EXPOSE 6571
CMD ["./akashgen-api"]
```

构建和运行:
```bash
docker build -t akashgen-api .
docker run -p 6571:6571 akashgen-api
```

## GPU 支持

服务会自动按顺序从首选 GPU 类型中选择：
1. RTX4090
2. A10
3. A100
4. V100-32Gi
5. H100

## 错误处理

API 提供详细的错误响应：

```json
{
    "error": "生成图像失败",
    "details": "具体错误描述"
}
```

常见错误场景：
- 无效的 JSON 格式
- 缺少必需字段
- 网络超时
- 上游服务失败
- 作业执行失败

## 性能

- **并发控制**: 可配置的请求限制防止过载
- **超时机制**: 多层超时防止请求挂起
- **连接池**: HTTP 客户端重用提高效率
- **内存管理**: 适当的资源清理
- **优雅关闭**: 干净的终止处理

## 监控

服务提供适合聚合的结构化 JSON 日志：

```json
{
    "level": "info",
    "ts": 1640995200.123456,
    "msg": "HTTP Request",
    "method": "POST",
    "path": "/api/generate",
    "status": 200,
    "duration": "2.5s",
    "client_ip": "192.168.1.100"
}
```

## 贡献

1. Fork 这个仓库
2. 创建你的功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

### 开发指南

- 遵循 Go 约定和最佳实践
- 为新功能添加测试
- 更新 API 变更的文档
- 使用有意义的提交消息
- 确保代码通过 `go fmt` 和 `go vet`

## 测试

```bash
# 运行测试
go test ./...

# 运行覆盖率测试
go test -cover ./...

# 运行竞态条件检测
go test -race ./...
```

## 故障排除

### 常见问题

1. **端口已被占用**: 在 `config/config.go` 中更改端口
2. **超时错误**: 对于较慢的网络，增加超时值
3. **内存问题**: 如果遇到内存压力，减少 `MaxConcurrentRequests`
4. **连接被拒绝**: 确保 Akash Network 端点可访问

### 调试模式

在开发环境中，可以通过修改 `main.go` 中的日志初始化来启用调试日志：

```go
// 将 zap.NewProduction() 替换为:
logger, err = zap.NewDevelopment()
```

## 许可证

该项目基于 AGPL v3 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 支持

- 为 bug 或功能请求创建 [issue](https://github.com/yourusername/akashgen-api-go/issues)
- 查看 [Akash Network 文档](https://docs.akash.network) 了解网络相关问题
- 参与 [Akash 社区](https://akash.network/community) 讨论

## 致谢

- [Akash Network](https://akash.network) 提供去中心化云基础设施
- [Gin Web Framework](https://github.com/gin-gonic/gin) 用于 HTTP 路由
- [Zap](https://github.com/uber-go/zap) 用于结构化日志

---

**注意**: 这是 Akash Network 图像生成 API 的非官方代理服务。请参考 Akash Network 官方文档获取有关其服务的最新信息。