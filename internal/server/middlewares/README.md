# Middlewares Package 说明

## 概述

Middlewares 包提供 API 模板默认需要的通用中间件能力，包括 JWT 认证、错误恢复、请求日志、压缩和 CORS。

## 中间件列表

### 1. JWT 中间件 (`jwt.go`)

**功能**: JWT 令牌认证和验证

**配置**:
```go
type JWTConfig struct {
    SigningKey []byte
    SkipPaths  []string
    Enabled    bool
}
```

**特性**:
- 支持路径跳过
- 可配置签名密钥
- 可启用或禁用
- 统一错误处理

### 2. 错误处理中间件 (`error.go`)

**功能**: 统一的错误处理和恢复

**特性**:
- 自动 panic 恢复
- 统一的错误响应格式
- 错误日志记录
- 支持业务错误码转换

**边界约定**:
- `internal/code` 是错误码、HTTP 状态、错误分类和对外安全消息的唯一来源。
- `ErrorHandler` 是 HTTP 错误边界，负责把 Echo 错误、验证错误、panic 和未知错误归一化为已注册错误码，并记录结构化日志。
- `internal/server/httpresp.APIError` 只负责根据 `code.ErrorInfo` 渲染 JSON，不做错误归类或日志记录。
- Service 层直接返回错误，Data 层包装存储错误，Biz 层只包装业务语义错误或透传 Repository 已编码错误。
- HTTP 状态码刻意收敛为少量通用状态，业务差异通过响应 `code` 表达，避免模板使用者在细粒度 HTTP 状态之间做重复决策。
- 业务错误可以返回具体、安全的 `message`；内部错误对外始终返回注册表中的安全文案，原始错误和诊断上下文只进入日志 `details`。
- 项目内需要读取错误码时使用 `code.ParseRegisteredCoder` 或 `code.NewErrorInfo`，不要直接依赖第三方错误库的默认解析行为。

### 3. 中间件配置 (`config.go`)

**功能**: 统一的中间件配置管理

**配置**:
```go
type MiddlewareConfig struct {
    EnableRecovery bool
    EnableLogger   bool
    EnableGzip     bool
    EnableCORS     bool
    EnableJWT      bool
    LoggerFormat   string
    JWT            *JWTConfig
}
```

**使用示例**:
```go
config := &MiddlewareConfig{
    EnableRecovery: true,
    EnableLogger:   true,
    EnableGzip:     true,
    EnableCORS:     true,
    EnableJWT:      true,
    LoggerFormat:   "method=${method}, uri=${uri}, status=${status}\n",
    JWT: &JWTConfig{
        SigningKey: []byte("your-secret"),
        SkipPaths:  []string{"/api/health", "/api/login"},
        Enabled:    true,
    },
}

ApplyMiddlewares(e, config)
```

## 扩展性

权限系统、租户边界、审计等业务相关中间件应由具体项目按需接入，避免 API 模板默认绑定特定授权实现。
