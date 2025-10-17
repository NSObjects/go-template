# 项目目录结构说明

## 整体架构

本项目采用DDD/六边形架构，目录结构清晰分层，职责明确。

## 目录结构

```
internal/
├── configs/                    # 配置管理
│   ├── bootstrap.go           # 配置启动
│   ├── config.go              # 配置结构定义
│   ├── store.go               # 配置存储
│   └── ...                    # 其他配置相关文件
│
├── infra/                      # 基础设施层（业务无关）
│   ├── persistence/           # 持久化基础设施
│   │   ├── mysql.go          # MySQL连接
│   │   ├── redis.go          # Redis连接
│   │   └── data_manager.go   # 数据管理器
│   └── wire.go               # 基础设施依赖注入
│
├── pkg/                        # 通用工具包（可复用）
│   ├── code/                  # 错误码管理
│   │   ├── code.go           # 错误码定义
│   │   ├── errors.go         # 错误处理
│   │   └── ...               # 其他错误相关文件
│   ├── utils/                 # 工具函数
│   │   ├── context.go        # 上下文工具
│   │   ├── jwt.go            # JWT工具
│   │   ├── validator.go      # 验证器
│   │   └── ...               # 其他工具函数
│   └── validator/             # 自定义验证器
│       └── custom.go         # 自定义验证规则
│
├── shared/                     # 共享内核
│   ├── kernel/                # 基础抽象
│   │   ├── entity.go         # 基础实体
│   │   ├── aggregate.go      # 聚合根标记
│   │   └── errors.go         # 领域错误
│   ├── event/                 # 事件基础设施
│   │   ├── event.go          # 事件接口
│   │   └── bus.go            # 事件总线接口
│   ├── infra/                 # 共享基础设施
│   │   ├── cache/            # 缓存实现
│   │   │   └── redis.go      # Redis缓存
│   │   ├── health/           # 健康检查
│   │   │   └── checker.go    # 健康检查器
│   │   ├── log/              # 日志系统
│   │   │   ├── logger.go     # 日志接口
│   │   │   ├── factory.go    # 日志工厂
│   │   │   └── ...           # 其他日志相关文件
│   │   ├── metrics/          # 指标监控
│   │   │   └── prometheus.go # Prometheus指标
│   │   ├── middleware/       # 中间件
│   │   │   └── rate_limit.go # 限流中间件
│   │   └── server/           # 服务器
│   │       ├── echo_server.go # Echo服务器
│   │       ├── config.go     # 服务器配置
│   │       └── middlewares/  # 服务器中间件
│   │           ├── error.go  # 错误处理中间件
│   │           ├── jwt.go    # JWT中间件
│   │           ├── casbin.go # 权限中间件
│   │           └── ...       # 其他中间件
│   └── ports/                 # 共享端口
│       └── resp/             # 响应处理
│           ├── response.go   # 统一响应格式
│           └── response_test.go # 响应测试
│
└── user/                      # User限界上下文
    ├── domain/                # 领域层（核心业务逻辑）
    │   ├── user.go           # User聚合根 + 值对象
    │   ├── user_repository.go # 仓储接口（出站端口）
    │   ├── events.go         # 领域事件
    │   └── user_test.go      # 领域层测试
    │
    ├── app/                   # 应用层（用例编排）
    │   ├── user_service.go   # 应用服务
    │   ├── dto.go            # DTO + Assembler
    │   └── ports.go          # 其他出站端口接口
    │
    ├── adapters/              # 适配器层
    │   ├── http_handler.go   # HTTP入站适配器
    │   ├── repository.go     # 仓储出站适配器实现
    │   ├── po.go             # 持久化对象 + Mapper
    │   └── tx_manager.go     # 事务管理实现
    │
    ├── wire.go                # 依赖注入
    └── arch_test.go           # 架构约束测试
```

## 目录职责说明

### 1. configs/ - 配置管理
- 负责应用配置的加载、解析和管理
- 支持多种配置源（文件、环境变量、远程配置等）
- 提供配置热重载功能

### 2. infra/ - 基础设施层
- 提供业务无关的基础设施服务
- 包含数据库连接、缓存等持久化基础设施
- 通过依赖注入为业务层提供服务

### 3. pkg/ - 通用工具包
- 包含可复用的通用工具和函数
- 不依赖业务逻辑，可在不同项目中复用
- 包含错误码、工具函数、验证器等

### 4. shared/ - 共享内核
- **kernel/**: 提供DDD基础抽象（实体、聚合根、错误等）
- **event/**: 提供事件基础设施（事件接口、事件总线等）
- **infra/**: 提供跨上下文共享的基础设施（日志、缓存、服务器等）
- **ports/**: 提供共享端口（响应处理等）

### 5. user/ - 限界上下文
- **domain/**: 领域层，包含核心业务逻辑
- **app/**: 应用层，负责用例编排
- **adapters/**: 适配器层，实现技术细节
- **wire.go**: 依赖注入配置
- **arch_test.go**: 架构约束测试

## 设计原则

### 1. 分层原则
- **领域层**：只能依赖shared，包含核心业务逻辑
- **应用层**：只能依赖domain和shared，负责用例编排
- **适配器层**：可以依赖domain、application和shared，实现技术细节

### 2. 依赖倒置
- 内层定义接口，外层实现
- 通过依赖注入管理组件生命周期

### 3. 上下文隔离
- 不同限界上下文之间不能直接依赖
- 跨上下文通信通过事件或防腐层

### 4. 职责单一
- 每个目录都有明确的职责
- 避免职责混乱和循环依赖

## 扩展指南

### 添加新的限界上下文
1. 在 `internal/` 下创建新的上下文目录（如 `order/`）
2. 按照 `user/` 的结构创建 `domain/`、`app/`、`adapters/` 目录
3. 实现相应的领域逻辑、应用服务和适配器
4. 添加依赖注入配置和架构测试

### 添加新的共享基础设施
1. 在 `internal/shared/infra/` 下创建新的基础设施目录
2. 实现相应的接口和实现
3. 在 `internal/shared/infra/wire.go` 中添加依赖注入配置

### 添加新的通用工具
1. 在 `internal/pkg/` 下创建新的工具目录
2. 实现可复用的工具函数
3. 添加相应的测试文件

## 最佳实践

1. **保持目录结构清晰**：每个目录都有明确的职责
2. **遵循依赖规则**：确保依赖关系符合架构原则
3. **添加架构测试**：通过测试确保架构约束得到遵守
4. **文档化**：为每个目录和重要文件添加说明文档
5. **持续重构**：随着业务发展持续优化目录结构
