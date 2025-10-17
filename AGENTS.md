# Repository Development Notes

This project scaffolds API layers (service, biz, data) from OpenAPI descriptions. When working under the repository, follow the
principles below so regenerated code and manual extensions continue to compose cleanly.

## Layer Responsibilities

### Service Layer (`internal/api/service`)
**Responsibility**: Transport layer handlers and request/response translation
- Perform syntactic validation, bind request parameters, and convert business responses into the unified `resp` envelope format
- **Prohibited**: Embed business rules or persistence logic. Service functions should remain thin, primarily orchestrating calls to biz interfaces and forwarding typed errors
- Use `resp.SuccessJSON`, `resp.ListDataResponse`, `resp.OneDataResponse`, `resp.OperateSuccess` to ensure all success responses conform to `{ "msg": "string", "code": 0, "data": { ... } }` format

**Correct Example**:
```go
func (c *UserController) CreateUser(ctx echo.Context) error {
    var req param.UserCreateRequest
    if err := BindAndValidate(ctx, &req); err != nil {
        return err  // Return directly, handled by ErrorHandler middleware
    }
    
    bizCtx := utils.BuildContext(ctx)
    err := c.user.Create(bizCtx, req)
    if err != nil {
        return err  // Return directly, handled by ErrorHandler middleware
    }
    
    return resp.OperateSuccess(ctx)  // Use unified response format
}

func (c *UserController) ListUsers(ctx echo.Context) error {
    var req param.UserListUsersRequest
    if err := BindAndValidate(ctx, &req); err != nil {
        return err
    }
    
    bizCtx := utils.BuildContext(ctx)
    list, total, err := c.user.ListUsers(bizCtx, req)
    if err != nil {
        return err
    }
    
    return resp.ListDataResponse(ctx, list, total)  // Use list response format
}
```

**Incorrect Example**:
```go
// ❌ Wrong: Direct database operations in Service layer
func (c *UserController) CreateUser(ctx echo.Context) error {
    var user model.User
    if err := c.db.Create(&user).Error; err != nil {
        return err
    }
    return ctx.JSON(200, user)  // ❌ Wrong: Not using unified response format
}

// ❌ Wrong: Business logic processing in Service layer
func (c *UserController) CreateUser(ctx echo.Context) error {
    var req param.UserCreateRequest
    if err := BindAndValidate(ctx, &req); err != nil {
        return err
    }
    
    // ❌ Wrong: Business logic should be in biz layer
    if req.Age < 18 {
        return errors.New("age must be greater than 18")
    }
    
    return c.user.Create(ctx, req)
}
```

### Biz Layer (`internal/api/biz`)
**Responsibility**: Encapsulate domain workflows and business invariants while staying storage-agnostic
- Coordinate repository interfaces, handle branching logic (e.g., rate limits, token rotation), and return rich error values from the `code` package
- Keep side-effects limited to calling repositories or emitting domain events; prefer pure functions for validation helpers
- Inject dependencies through constructors, avoid direct dependency on `DataManager`

**Correct Example**:
```go
type UserHandler struct {
    repo UserRepository  // Inject through interface, not direct DataManager dependency
}

func NewUserHandler(repo UserRepository) UserUseCase {
    return &UserHandler{repo: repo}
}

func (h *UserHandler) CreateUser(ctx context.Context, req param.UserCreateRequest) error {
    // Business rule validation
    if err := h.validateUserData(req); err != nil {
        return code.WrapValidationError(err, "user data validation failed")
    }
    
    // Call Repository layer
    err := h.repo.Create(ctx, req)
    if err != nil {
        return code.WrapDatabaseError(err, "create user failed")
    }
    
    return nil
}

// Pure function validation helper
func (h *UserHandler) validateUserData(req param.UserCreateRequest) error {
    if req.Age < 18 {
        return errors.New("age must be greater than 18")
    }
    return nil
}
```

**Incorrect Example**:
```go
// ❌ Wrong: Direct dependency on DataManager
type UserHandler struct {
    dm *data.DataManager  // ❌ Wrong: Should inject through Repository interface
}

func (h *UserHandler) CreateUser(ctx context.Context, req param.UserCreateRequest) error {
    // ❌ Wrong: Direct database operations in Biz layer
    var user model.User
    if err := h.dm.MySQLWithContext(ctx).Create(&user).Error; err != nil {
        return err  // ❌ Wrong: Not using code package to wrap errors
    }
    return nil
}
```

### Data Layer (`internal/api/data`)
**Responsibility**: Implement concrete persistence adapters and external integrations
- Expose clear interfaces that can be mocked; avoid leaking transport concepts (HTTP status, envelopes) into this layer
- Implement Repository interfaces, focus on data access logic
- Use `code` package to wrap database errors, provide meaningful error messages

**Correct Example**:
```go
type userRepository struct {
    d *db.DataManager
}

func NewUserRepository(d *db.DataManager) biz.UserRepository {
    return userRepository{d: d}
}

func (u userRepository) Create(ctx context.Context, req param.UserCreateRequest) error {
    user := model.User{
        Username:  req.Username,
        Email:     req.Email,
        Age:       int32(req.Age),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    err := u.d.Query.User.WithContext(ctx).Create(&user)
    if err != nil {
        return code.WrapDatabaseError(err, "create user failed")
    }
    
    return nil
}

func (u userRepository) ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {
    users, err := u.d.Query.User.WithContext(ctx).Offset(req.Offset()).Limit(req.Limit()).Find()
    if err != nil {
        return nil, 0, code.WrapDatabaseError(err, "query user list failed")
    }
    
    var list []param.UserListItem
    for _, user := range users {
        list = append(list, param.UserListItem{
            Id:        user.ID,
            Username:  user.Username,
            Email:     user.Email,
            Age:       int(user.Age),
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
        })
    }
    
    count, err := u.d.Query.User.Count()
    if err != nil {
        return nil, 0, code.WrapDatabaseError(err, "query user count failed")
    }
    
    return list, count, nil
}
```

**Incorrect Example**:
```go
// ❌ Wrong: Business logic processing in Data layer
func (u userRepository) Create(ctx context.Context, req param.UserCreateRequest) error {
    // ❌ Wrong: Business logic should be in biz layer
    if req.Age < 18 {
        return errors.New("age must be greater than 18")
    }
    
    // Database operations...
}

// ❌ Wrong: Not wrapping database errors
func (u userRepository) Create(ctx context.Context, req param.UserCreateRequest) error {
    err := u.d.Query.User.WithContext(ctx).Create(&user)
    if err != nil {
        return err  // ❌ Wrong: Should use code package to wrap
    }
    return nil
}
```

## Coding Conventions

### Dependency Injection Principles
- Favor composition over shared state. Inject dependencies via constructors to keep components testable
- Prefer Go interfaces that describe behaviors needed by the caller rather than concrete types

**Correct Example**:
```go
// Define Repository interface
type UserRepository interface {
    Create(ctx context.Context, req param.UserCreateRequest) error
    GetByID(ctx context.Context, id int64) (param.UserData, error)
    ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error)
    Update(ctx context.Context, id int64, req param.UserUpdateRequest) error
    Delete(ctx context.Context, id int64) error
}

// Inject through interface in Biz layer
type UserHandler struct {
    repo UserRepository
}
```

### Error Handling Standards
- Wrap domain failures with helpers defined in `internal/code` package. These errors will be translated by the global `ErrorHandler` middleware
- Success responses must go through `resp` package methods to ensure `{ "msg": "string", "code": 0, "data": { ... } }` contract

**Error Code Usage Example**:
```go
// Define error codes in code package
const (
    ErrUserNotFound = 10001
    ErrUserAlreadyExists = 10002
)

// Use in Data layer
func (u userRepository) GetByID(ctx context.Context, id int64) (param.UserData, error) {
    user, err := u.d.Query.User.WithContext(ctx).GetByID(uint(id))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return param.UserData{}, code.WrapNotFoundError(err, "user not found")
        }
        return param.UserData{}, code.WrapDatabaseError(err, "query user failed")
    }
    return convertToUserData(user), nil
}
```

### Response Format Standards
Use the unified response methods provided by `resp` package:

```go
// List data response
return resp.ListDataResponse(ctx, list, total)

// Single data response  
return resp.OneDataResponse(ctx, data)

// Operation success response
return resp.OperateSuccess(ctx)

// Custom success response
return resp.SuccessJSON(ctx, map[string]interface{}{
    "message": "operation successful",
    "data": result,
})
```

### File Organization
- Organize files by concern: new business logic belongs beside the generated scaffolding, while shared helpers should live in dedicated packages (`internal/resp`, `internal/code`, etc.)
- Use context-aware functions (`ctx context.Context`) for operations that may need cancellation or tracing
- Handle secrets and credentials via configuration structs—avoid hard-coding values in source files

## Code Generation Guidelines

### Template Generation Principles
1. **Service layer templates** must:
   - Use `resp` package methods for all success responses
   - Return errors directly, let ErrorHandler middleware handle them
   - Not contain any business logic or database operations
   - Use `BindAndValidate` for parameter binding and validation

2. **Biz layer templates** must:
   - Define Repository interfaces
   - Inject Repository dependencies through constructors
   - Use `code` package to wrap all errors
   - Not directly depend on `DataManager`
   - Include business rule validation

3. **Data layer templates** must:
   - Implement Repository interfaces
   - Focus on data access logic
   - Use `code` package to wrap database errors
   - Not contain business rules

### Response Format Requirements
All API responses must use the following format:
```json
{
  "msg": "operation successful",
  "code": 0,
  "data": {
    // actual data
  }
}
```

### Error Handling Requirements
- Use `internal/code` package to define error codes
- Use `code.WrapXXXError` in Data layer to wrap errors
- Use `code.WrapXXXError` in Biz layer to wrap Repository errors
- Service layer returns errors directly, no additional processing

## Testing Expectations

- Write table-driven tests for new logic and prefer standard library assertions (`if`, `t.Helper`, etc.)
- Service layer tests should exercise request binding and response envelopes
- Biz layer tests should mock repositories to cover edge cases (lockouts, rotations, etc.)
- Data layer tests should use test databases for integration testing
- Run `go test ./...` (or a targeted subset) before submitting changes; include any non-standard flags used in the test output

## Working With Generated Code

- The OpenAPI specification is the source of truth for generated artifacts. Avoid editing generated files directly—extend behavior in protected regions or create new files to survive regeneration
- When updating templates or tooling that emits these layers, regenerate from the spec and ensure the repository still builds and passes tests prior to merging

## Code Review Checklist

When reviewing generated code, ensure:

- [ ] Service layer only contains request binding, parameter validation, and response formatting
- [ ] Biz layer accesses data through Repository interfaces, not directly depending on DataManager
- [ ] Data layer implements Repository interfaces, focusing on data access
- [ ] All errors are wrapped using `code` package
- [ ] All success responses use `resp` package methods
- [ ] Dependencies are injected through constructors, not global variables
- [ ] Interface definitions are near callers, implementations away from callers
- [ ] Context is properly passed in all async operations
- [ ] Business rule validation is in Biz layer
- [ ] Data conversion is in Data layer

## DDD/六边形架构规范

本项目同时支持传统三层架构和DDD/六边形架构。新项目推荐使用DDD架构，现有项目可渐进式迁移。

### DDD架构目录结构

```
internal/
├── <bounded-context>/              # 限界上下文（如 user, order）
│   ├── domain/                     # 领域层（核心业务逻辑）
│   │   ├── <entity>.go            # 聚合根 + 值对象
│   │   ├── <entity>_repository.go # 仓储接口（出站端口）
│   │   └── events.go              # 领域事件
│   │
│   ├── app/                        # 应用层（用例编排）
│   │   ├── <entity>_service.go    # 应用服务（包含入站端口接口）
│   │   ├── dto.go                 # DTO + Assembler
│   │   └── ports.go               # 其他出站端口（TxManager/EventBus等）
│   │
│   ├── adapters/                   # 适配器层
│   │   ├── http_handler.go        # HTTP入站适配器
│   │   ├── repository.go          # 仓储出站适配器实现
│   │   ├── po.go                  # 持久化对象 + Mapper
│   │   └── tx_manager.go          # 事务管理实现
│   │
│   ├── wire.go                     # 依赖注入（fx）
│   └── arch_test.go                # 架构约束测试
│
├── infra/                          # 基础设施层（跨上下文共享）
│   ├── persistence/               # 持久化基础设施
│   │   ├── mysql.go               # MySQL连接池 + GORM初始化
│   │   ├── redis.go               # Redis客户端初始化
│   │   └── data_manager.go        # DataManager（包装DB连接）
│   └── wire.go                    # 基础设施依赖注入
│
└── shared/                         # 共享内核
    ├── kernel/                     # 基础抽象
    │   ├── entity.go              # 基础实体（ID/时间戳）
    │   ├── aggregate.go           # 聚合根标记
    │   └── errors.go              # 领域错误
    └── event/                      # 事件基础设施
        ├── event.go               # 事件接口/信封
        └── bus.go                 # 事件总线接口
```

### DDD架构原则

#### 1. 依赖规则（关键约束）

**层级依赖规则**：
- `adapters/**` 只能依赖 `application/**`, `domain/**`, `shared/**`；不得被反向依赖
- `application/**` 只能依赖 `domain/**`, `shared/**`；不得依赖 `adapters/**`（除了构造期wiring）
- `domain/**` 仅允许依赖 `shared/**` 的极小抽象（如 `errors/kernel`）
- `port/**` 由内层定义，外层实现（依赖倒置）

**上下文隔离规则**：
- 不同限界上下文（如 `user/`, `order/`）之间不得直接依赖对方的 `domain/**` 或 `application/**`
- 跨上下文通信只能通过：
  1. **集成事件**（Integration Events）- 异步通信
  2. **防腐层ACL**（Anti-Corruption Layer）- 同步调用时翻译外部模型
  3. **共享内核**（Shared Kernel）- 仅限极小的通用抽象

#### 2. 领域层职责

**聚合根（Aggregate Root）**：
- 封装业务规则和不变性
- 管理聚合内的一致性边界
- 发布领域事件
- 无ORM标签，保持纯净

**值对象（Value Object）**：
- 不可变
- 包含验证逻辑
- 通过值比较相等性

**仓储接口（Repository Interface）**：
- 定义在领域层
- 由适配器层实现
- 使用领域语言

**领域事件（Domain Events）**：
- 表示业务中发生的重要事情
- 由聚合根发布
- 用于跨聚合通信

#### 3. 应用层职责

**应用服务（Application Service）**：
- 用例编排
- 事务边界管理
- 调用领域服务/实体方法
- 调用出站端口（Repository）
- 发布领域事件
- DTO与Entity的转换

**禁止**：
- 包含业务规则（属于领域层）
- 直接操作数据库（通过Repository端口）

#### 4. 适配器层职责

**入站适配器（Inbound Adapters）**：
- HTTP控制器
- 消息队列消费者
- 定时任务

**出站适配器（Outbound Adapters）**：
- 数据库仓储实现
- 外部API客户端
- 消息队列生产者
- 缓存实现

### DDD开发规范

#### 1. 领域实体示例

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

#### 2. 应用服务示例

```go
// internal/user/app/user_service.go
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

#### 3. 仓储实现示例

```go
// internal/user/adapters/repository.go
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

#### 4. 架构约束测试

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
        // ... 更多测试
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            checkDependencies(t, tt.dir, tt.allowedDeps, tt.prohibitedDeps)
        })
    }
}
```

### 迁移指南

#### 从三层架构到DDD架构

1. **创建新的限界上下文目录结构**
2. **迁移领域逻辑**：将业务规则从biz层移到domain层
3. **重构数据访问**：将data层改为adapters层，实现domain层的仓储接口
4. **更新应用服务**：将biz层改为app层，专注于用例编排
5. **更新控制器**：将service层改为adapters层，调用app层服务
6. **配置依赖注入**：更新fx模块配置
7. **添加架构测试**：确保依赖关系正确

#### 代码生成支持

项目支持两种架构模式的代码生成：

```bash
# 传统三层架构（默认）
go run muban/main.go modgen openapi --spec=api.yaml --output=internal/api

# DDD/六边形架构
go run muban/main.go modgen openapi --spec=api.yaml --output=internal --architecture=ddd-hex
```

### 最佳实践

1. **保持聚合根小**：一个聚合根只管理一个业务概念
2. **使用值对象**：封装验证逻辑和业务规则
3. **发布领域事件**：用于跨聚合通信
4. **依赖注入**：通过构造函数注入依赖
5. **测试驱动**：先写测试，再写实现
6. **架构约束**：使用测试确保依赖关系正确
7. **渐进式迁移**：新旧架构可以共存，逐步迁移