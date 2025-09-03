# Middlewares Package 说明

## 概述

Middlewares包提供了完整的中间件系统，包括JWT认证、Casbin权限控制、错误处理、日志记录等功能。

## 中间件列表

### 1. JWT中间件 (`jwt.go`)

**功能**: JWT令牌认证和验证

**配置**:
```go
type JWTConfig struct {
    SigningKey []byte   // 签名密钥
    SkipPaths  []string // 跳过路径
    Enabled    bool     // 是否启用
}
```

**特性**:
- 支持路径跳过（精确匹配和通配符匹配）
- 可配置的签名密钥
- 可启用/禁用
- 自动错误处理

**使用示例**:
```go
config := &JWTConfig{
    SigningKey: []byte("your-secret"),
    SkipPaths:  []string{"/api/health", "/api/login"},
    Enabled:    true,
}
e.Use(JWT(config))
```

### 2. Casbin中间件 (`casbin.go`)

**功能**: 基于Casbin的权限控制

**配置**:
```go
type CasbinConfig struct {
    Enabled    bool     // 是否启用
    SkipPaths  []string // 跳过路径
    AdminUsers []string // 管理员用户
}
```

**特性**:
- 支持管理员用户绕过权限检查
- 可配置的跳过路径
- 基于路径和HTTP方法的权限控制
- 可启用/禁用

**使用示例**:
```go
config := &CasbinConfig{
    Enabled:    true,
    SkipPaths:  []string{"/api/health", "/api/info"},
    AdminUsers: []string{"root", "admin"},
}
e.Use(Casbin(enforcer, config))
```

### 3. 错误处理中间件 (`error.go`)

**功能**: 统一的错误处理和恢复

**特性**:
- 自动panic恢复
- 统一的错误响应格式
- 详细的错误日志记录
- 支持不同类型的错误处理

### 4. 中间件配置 (`config.go`)

**功能**: 统一的中间件配置管理

**配置**:
```go
type MiddlewareConfig struct {
    EnableRecovery bool
    EnableLogger   bool
    EnableGzip     bool
    EnableCORS     bool
    EnableJWT      bool
    EnableCasbin   bool
    LoggerFormat   string
    JWT            *JWTConfig
    Casbin         *CasbinConfig
}
```

**特性**:
- 统一的配置管理
- 支持启用/禁用特定中间件
- 可配置的日志格式
- 自动应用中间件

## 配置示例

### 完整配置
```go
config := &MiddlewareConfig{
    EnableRecovery: true,
    EnableLogger:   true,
    EnableGzip:     true,
    EnableCORS:     true,
    EnableJWT:      true,
    EnableCasbin:   true,
    LoggerFormat:   "method=${method}, uri=${uri}, status=${status}\n",
    JWT: &JWTConfig{
        SigningKey: []byte("your-secret"),
        SkipPaths:  []string{"/api/health", "/api/login"},
        Enabled:    true,
    },
    Casbin: &CasbinConfig{
        Enabled:    true,
        SkipPaths:  []string{"/api/health", "/api/info"},
        AdminUsers: []string{"root", "admin"},
    },
}

// 应用中间件
ApplyMiddlewares(e, config)
ApplyCasbinMiddleware(e, enforcer, config.Casbin)
```

### 最小配置
```go
// 使用默认配置
ApplyMiddlewares(e, nil)
```

## 安全特性

1. **JWT安全**:
   - 强密钥验证
   - 路径跳过控制
   - 自动错误处理

2. **权限控制**:
   - 基于角色的访问控制
   - 管理员用户支持
   - 细粒度权限管理

3. **错误处理**:
   - 不泄露敏感信息
   - 统一的错误格式
   - 详细的日志记录

## 性能优化

1. **中间件按需加载**: 只加载启用的中间件
2. **路径跳过优化**: 支持通配符匹配
3. **配置缓存**: 避免重复配置
4. **错误恢复**: 防止panic导致服务崩溃

## 扩展性

1. **模块化设计**: 每个中间件独立
2. **配置驱动**: 通过配置控制行为
3. **可插拔**: 支持自定义中间件
4. **向后兼容**: 保持API稳定性
