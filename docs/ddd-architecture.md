# DDD/六边形架构设计文档

## 概述

本文档详细描述了项目中DDD（领域驱动设计）/六边形架构的实现方案。该架构将业务逻辑与基础设施分离，通过依赖倒置原则实现高度可测试和可维护的代码结构。

## 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Controller                          │
│                  (Inbound Adapter)                         │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                Application Service                          │
│                 (Use Cases)                                │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                  Domain Layer                              │
│              (Entities, Value Objects,                     │
│               Domain Services, Events)                     │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                Repository Interface                        │
│                 (Outbound Port)                            │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              Repository Implementation                     │
│                (Outbound Adapter)                          │
└─────────────────────────────────────────────────────────────┘
```

## 核心概念

### 1. 限界上下文（Bounded Context）

每个业务领域形成一个独立的限界上下文，包含完整的DDD架构层次：

```
internal/
├── user/                    # User限界上下文
├── order/                   # Order限界上下文
├── product/                 # Product限界上下文
└── payment/                 # Payment限界上下文
```

### 2. 领域层（Domain Layer）

**职责**：包含核心业务逻辑，不依赖任何外部技术

**组件**：
- **聚合根（Aggregate Root）**：业务实体的根对象
- **值对象（Value Object）**：不可变的对象，通过值比较相等性
- **领域服务（Domain Service）**：不属于任何实体的业务逻辑
- **领域事件（Domain Events）**：业务中发生的重要事情
- **仓储接口（Repository Interface）**：数据访问的抽象

**示例**：
```go
// internal/user/domain/user.go
type User struct {
    kernel.BaseEntity
    kernel.BaseAggregateRoot
    
    username  string
    email     Email           // 值对象
    birthDate BirthDate       // 值对象
    status    UserStatus
}

func (u *User) ChangeEmail(newEmailStr string) error {
    // 业务规则验证
    if u.status != UserStatusActive {
        return errors.New("cannot change email for inactive user")
    }
    
    newEmail, err := NewEmail(newEmailStr)
    if err != nil {
        return err
    }
    
    if u.email.Equals(newEmail) {
        return errors.New("new email is the same as current email")
    }
    
    oldEmail := u.email.Value()
    u.email = newEmail
    u.MarkUpdated()
    
    // 发布领域事件
    emailChanged := &UserEmailChanged{
        BaseEvent: *event.NewBaseEvent("UserEmailChanged", u.ID().Value()),
        UserID:    u.ID().Value(),
        OldEmail:  oldEmail,
        NewEmail:  newEmail.Value(),
    }
    u.AddDomainEvent(emailChanged)
    
    return nil
}
```

### 3. 应用层（Application Layer）

**职责**：用例编排，协调领域对象完成业务操作

**组件**：
- **应用服务（Application Service）**：实现用例
- **DTO（Data Transfer Object）**：数据传输对象
- **Assembler**：DTO与领域对象的转换器
- **入站端口（Inbound Port）**：应用服务接口

**示例**：
```go
// internal/user/app/user_service.go
type UserService interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (*UserDTO, error)
    GetUser(ctx context.Context, userID string) (*UserDTO, error)
    UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) error
    ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error)
    DeleteUser(ctx context.Context, userID string) error
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req CreateUserRequest) (*UserDTO, error) {
    // 验证邮箱是否已存在
    email, err := domain.NewEmail(req.Email)
    if err != nil {
        return nil, kernel.NewBusinessRuleError("INVALID_EMAIL", err.Error())
    }
    
    existingUser, err := s.userRepo.FindByEmail(ctx, email)
    if err == nil && existingUser != nil {
        return nil, kernel.NewBusinessRuleError("EMAIL_EXISTS", "email already exists")
    }
    
    // 创建用户聚合根
    userID, err := domain.NewUserID(req.UserID)
    if err != nil {
        return nil, kernel.NewBusinessRuleError("INVALID_USER_ID", err.Error())
    }
    
    birthDate, err := time.Parse("2006-01-02", req.BirthDate)
    if err != nil {
        return nil, kernel.NewBusinessRuleError("INVALID_BIRTH_DATE", "invalid birth date format")
    }
    
    user, err := domain.NewUser(userID, req.Username, req.Email, birthDate)
    if err != nil {
        return nil, kernel.NewBusinessRuleError("CREATE_USER_FAILED", err.Error())
    }
    
    // 在事务中保存用户
    err = s.txManager.ExecuteInTransaction(ctx, func(ctx context.Context) error {
        if err := s.userRepo.Save(ctx, user); err != nil {
            return err
        }
        
        // 发布领域事件
        for _, event := range user.GetUncommittedEvents() {
            if err := s.eventBus.Publish(ctx, event); err != nil {
                return err
            }
        }
        user.MarkEventsAsCommitted()
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    return UserAssembler{}.ToDTO(user), nil
}
```

### 4. 适配器层（Adapter Layer）

**职责**：实现端口接口，连接外部世界与内部业务逻辑

**组件**：
- **入站适配器（Inbound Adapters）**：HTTP控制器、消息队列消费者等
- **出站适配器（Outbound Adapters）**：数据库仓储实现、外部API客户端等

**示例**：
```go
// internal/user/adapters/http_handler.go
type UserController struct {
    userService app.UserService
}

func (c *UserController) CreateUser(ctx echo.Context) error {
    var req app.CreateUserRequest
    if err := ctx.Bind(&req); err != nil {
        return code.WrapValidationError(err, "bind request failed")
    }
    
    if err := ctx.Validate(&req); err != nil {
        return code.WrapValidationError(err, "validation failed")
    }
    
    bizCtx := utils.BuildContext(ctx)
    user, err := c.userService.CreateUser(bizCtx, req)
    if err != nil {
        return err
    }
    
    return resp.OneDataResponse(ctx, user)
}

// internal/user/adapters/repository.go
type userRepositoryImpl struct {
    db *db.DataManager
    mapper UserMapper
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *domain.User) error {
    po := r.mapper.ToPO(user)
    
    // 检查用户是否已存在
    userID, err := domain.NewUserID(user.ID().Value())
    if err != nil {
        return code.WrapDatabaseError(err, "invalid user ID")
    }
    
    exists, err := r.Exists(ctx, userID)
    if err != nil {
        return code.WrapDatabaseError(err, "check user exists failed")
    }
    
    if exists {
        // 更新用户
        err = r.db.MySQLWithContext(ctx).Model(&UserPO{}).Where("id = ?", po.ID).Updates(po).Error
    } else {
        // 创建用户
        err = r.db.MySQLWithContext(ctx).Create(po).Error
    }
    
    if err != nil {
        return code.WrapDatabaseError(err, "save user failed")
    }
    
    return nil
}
```

## 依赖规则

### 1. 层级依赖规则

- `adapters/**` 只能依赖 `application/**`, `domain/**`, `shared/**`；不得被反向依赖
- `application/**` 只能依赖 `domain/**`, `shared/**`；不得依赖 `adapters/**`（除了构造期wiring）
- `domain/**` 仅允许依赖 `shared/**` 的极小抽象（如 `errors/kernel`）
- `port/**` 由内层定义，外层实现（依赖倒置）

### 2. 上下文隔离规则

- 不同限界上下文之间不得直接依赖对方的 `domain/**` 或 `application/**`
- 跨上下文通信只能通过：
  1. **集成事件**（Integration Events）- 异步通信
  2. **防腐层ACL**（Anti-Corruption Layer）- 同步调用时翻译外部模型
  3. **共享内核**（Shared Kernel）- 仅限极小的通用抽象

### 3. 架构约束测试

```go
// internal/user/arch_test.go
func TestArchitectureConstraints(t *testing.T) {
    tests := []struct {
        name        string
        dir         string
        allowedDeps []string
        prohibitedDeps []string
    }{
        {
            name: "domain layer should not depend on adapters or application",
            dir:  "domain",
            allowedDeps: []string{
                "github.com/NSObjects/go-template/internal/shared",
            },
            prohibitedDeps: []string{
                "github.com/NSObjects/go-template/internal/user/adapters",
                "github.com/NSObjects/go-template/internal/user/app",
            },
        },
        {
            name: "application layer should not depend on adapters",
            dir:  "app",
            allowedDeps: []string{
                "github.com/NSObjects/go-template/internal/user/domain",
                "github.com/NSObjects/go-template/internal/shared",
            },
            prohibitedDeps: []string{
                "github.com/NSObjects/go-template/internal/user/adapters",
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            checkDependencies(t, tt.dir, tt.allowedDeps, tt.prohibitedDeps)
        })
    }
}
```

## 事件驱动架构

### 1. 领域事件

领域事件表示业务中发生的重要事情，由聚合根发布：

```go
// internal/user/domain/events.go
type UserCreated struct {
    event.BaseEvent
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// 在聚合根中发布事件
func (u *User) ChangeEmail(newEmailStr string) error {
    // ... 业务逻辑 ...
    
    // 发布领域事件
    emailChanged := &UserEmailChanged{
        BaseEvent: *event.NewBaseEvent("UserEmailChanged", u.ID().Value()),
        UserID:    u.ID().Value(),
        OldEmail:  oldEmail,
        NewEmail:  newEmail.Value(),
    }
    u.AddDomainEvent(emailChanged)
    
    return nil
}
```

### 2. 事件总线

事件总线负责发布和订阅领域事件：

```go
// internal/shared/event/bus.go
type EventBus interface {
    Publish(ctx context.Context, event Event) error
    Subscribe(handler EventHandler)
    Start(ctx context.Context) error
    Stop() error
}

type EventHandler interface {
    Handle(ctx context.Context, event Event) error
    EventType() string
}
```

## 测试策略

### 1. 单元测试

- **领域层**：测试业务规则和领域逻辑
- **应用层**：测试用例编排，使用Mock对象
- **适配器层**：测试与外部系统的集成

### 2. 集成测试

- 测试完整的用例流程
- 使用测试数据库
- 验证事件发布

### 3. 架构测试

- 验证依赖关系正确性
- 确保架构约束得到遵守

## 代码生成支持

项目支持从OpenAPI规范生成DDD架构代码：

```bash
# 生成DDD架构代码
go run muban/main.go modgen openapi --spec=api.yaml --output=internal --architecture=ddd-hex
```

## 迁移指南

### 从三层架构到DDD架构

1. **创建新的限界上下文目录结构**
2. **迁移领域逻辑**：将业务规则从biz层移到domain层
3. **重构数据访问**：将data层改为adapters层，实现domain层的仓储接口
4. **更新应用服务**：将biz层改为app层，专注于用例编排
5. **更新控制器**：将service层改为adapters层，调用app层服务
6. **配置依赖注入**：更新fx模块配置
7. **添加架构测试**：确保依赖关系正确

### 渐进式迁移

- 新旧架构可以共存
- 新功能使用DDD架构
- 现有功能逐步迁移
- 通过测试确保迁移正确性

## 最佳实践

1. **保持聚合根小**：一个聚合根只管理一个业务概念
2. **使用值对象**：封装验证逻辑和业务规则
3. **发布领域事件**：用于跨聚合通信
4. **依赖注入**：通过构造函数注入依赖
5. **测试驱动**：先写测试，再写实现
6. **架构约束**：使用测试确保依赖关系正确
7. **渐进式迁移**：新旧架构可以共存，逐步迁移

## 总结

DDD/六边形架构通过清晰的层次划分和依赖规则，实现了：

- **高内聚低耦合**：每层职责明确，依赖关系清晰
- **可测试性**：通过依赖注入和接口抽象，易于单元测试
- **可维护性**：业务逻辑与技术实现分离，易于修改和扩展
- **可扩展性**：支持微服务架构的演进路径

这种架构特别适合复杂业务场景，能够有效管理业务复杂度和技术复杂度。
