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