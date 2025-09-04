# 代码质量检查 (Linting)

本项目使用 `golangci-lint` 进行代码质量检查，确保代码符合最佳实践和编码规范。

## 安装

```bash
# 安装 golangci-lint
make install-lint

# 或者手动安装
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## 可用的 Lint 命令

### 基础检查

```bash
# 运行所有 linter（允许失败）
make lint

# 严格检查（发现问题时退出）
make lint-strict

# 快速检查（只运行快速 linter）
make lint-fast
```

### 修复和报告

```bash
# 自动修复可修复的问题
make lint-fix

# 生成 checkstyle 格式的报告
make lint-report
```

### 特定目录检查

```bash
# 检查特定目录
make lint-dir DIR=./internal/api
make lint-dir DIR=./internal/code
```

## 当前发现的问题

运行 `make lint` 后，项目中发现以下类型的问题：

- **errcheck (6个)**: 未检查的错误返回值
- **govet (4个)**: 格式字符串问题
- **staticcheck (19个)**: 代码质量问题
- **unused (1个)**: 未使用的字段

## 常见问题修复

### 1. 未检查的错误返回值 (errcheck)

```go
// ❌ 错误
defer resp.Body.Close()

// ✅ 正确
defer func() {
    if err := resp.Body.Close(); err != nil {
        log.Printf("Error closing response body: %v", err)
    }
}()
```

### 2. 格式字符串问题 (govet)

```go
// ❌ 错误
g.Printf(buf.String())

// ✅ 正确
g.Printf("%s", buf.String())
```

### 3. 未使用的字段 (unused)

```go
// ❌ 错误
type Logger struct {
    mu   sync.RWMutex  // 未使用
    sink Sink
}

// ✅ 正确 - 删除未使用的字段或使用它
type Logger struct {
    sink Sink
}
```

## 集成到开发流程

### 1. 开发前检查

```bash
# 运行快速检查
make lint-fast

# 如果有问题，尝试自动修复
make lint-fix
```

### 2. 提交前检查

```bash
# 运行完整检查
make lint-strict

# 生成报告
make lint-report
```

### 3. CI/CD 集成

在 CI 流程中添加 lint 检查：

```yaml
- name: Run golangci-lint
  run: make lint-strict
```

## 配置说明

当前配置跳过了以下目录和文件：

- `vendor/` - 第三方依赖
- `tools/vendor/` - 工具依赖
- `internal/api/data/query/` - 生成的查询代码
- `.*\.gen\.go$` - 所有生成的 Go 文件

## 最佳实践

1. **定期运行 lint 检查**：在开发过程中定期运行 `make lint-fast`
2. **修复问题**：及时修复发现的问题，保持代码质量
3. **使用自动修复**：对于可以自动修复的问题，使用 `make lint-fix`
4. **严格模式**：在 CI/CD 中使用 `make lint-strict` 确保代码质量
5. **生成报告**：定期生成 lint 报告，跟踪代码质量趋势

## 参考资源

- [golangci-lint 官方文档](https://golangci-lint.run/)
- [Go 代码审查指南](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
