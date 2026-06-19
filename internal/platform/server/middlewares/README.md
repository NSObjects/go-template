# Middlewares Package 说明

## 概述

Middlewares 包提供 API 模板的通用中间件能力，包括请求元数据、JWT 验证、错误恢复、请求日志、压缩和可选 CORS。

## 中间件列表

### 1. 请求元数据中间件 (`request_context.go`)

**功能**: 把 request id、trace id、span id 等请求元数据写入标准 `context.Context`。

**边界约定**:
- 该中间件只搬运请求元数据，不做认证授权。
- request id、trace id、span id 只接受长度受限的可见 ASCII 值；非法 request id 会重新生成，非法 trace/span id 会丢弃。
- 该中间件不读取 `X-User-ID` 这类客户端身份 header；真实用户身份必须由认证边界验证后写入 context。
- usecase 层通过 `internal/platform/requestctx` 读取元数据，不依赖 Echo context。

### 2. JWT 中间件 (`jwt.go`)

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
- 验证成功后把标准 `sub` claim 写入 `requestctx.UserID`；业务需要其他身份语义时，应由真实 auth 模块显式转换

### 3. 错误处理中间件 (`error.go`)

**功能**: 统一的错误处理和恢复

**特性**:
- 自动 panic 恢复
- 统一的错误响应格式
- 错误日志记录默认只记录 URL path，不记录 query string
- 支持业务错误码转换

**边界约定**:
- `internal/platform/apperr` 是错误码、错误类别和对外安全消息的唯一来源，并且不依赖 HTTP。
- `ErrorHandler` 是 HTTP 错误边界，负责把 Echo 错误、panic 和未知错误归一化为应用错误，并记录结构化日志。
- `internal/platform/server/httpresp.APIError` 负责把 `apperr.Info` 映射为 HTTP 状态码和 JSON 响应。
- Usecase 层直接返回错误，Adapter 层包装外部系统错误，业务模块只包装业务语义错误或透传已编码错误。
- 业务错误可以返回具体、安全的 `message`；内部错误对外始终返回安全文案，原始错误和诊断上下文只进入日志 `detail`。

### 4. 中间件配置 (`config.go`)

**功能**: 统一的中间件配置管理。请求日志使用 boot 阶段安装的 zerolog logger，并把 request id、trace id、span id、user id 等字段挂到当前请求 context。

**配置**:
```go
type MiddlewareConfig struct {
    EnableRecovery       bool
    EnableRequestContext bool
    EnableLogger         bool
    EnableGzip           bool
    EnableCORS           bool
    CORS                 middleware.CORSConfig
    EnableJWT            bool
    JWT                  *JWTConfig
}
```

**使用示例**:
```go
config := &MiddlewareConfig{
    EnableRecovery:       true,
    EnableRequestContext: true,
    EnableLogger:         true,
    EnableGzip:           true,
    EnableCORS:           true,
    CORS: middleware.CORSConfig{
        AllowOrigins: []string{"https://app.example.com"},
    },
    EnableJWT:    true,
    JWT: &JWTConfig{
        SigningKey: []byte(secretFromConfig),
        SkipPaths:  []string{"/api/health", "/api/info", "/api/ready"},
        Enabled:    true,
    },
}

ApplyMiddlewares(e, config)
```

`server.New` 使用保守默认值：不启用 CORS。项目真的需要跨域时，在静态配置或环境变量里显式打开 `http.cors.enabled`，并确认允许的 origin，避免模板默认放大浏览器访问面。

## 扩展性

权限系统、租户边界、审计等业务相关中间件应由具体项目按需接入，避免 API 模板默认绑定特定授权实现。
