# golangci-lint 配置总结

## ✅ 完成的工作

### 1. 配置 golangci-lint
- 使用 golangci-lint v2.4.0
- 由于配置文件格式兼容性问题，采用命令行参数配置
- 跳过 vendor 目录和生成的代码文件

### 2. 集成到 Makefile
添加了完整的 lint 命令集：

```bash
# 基础检查
make lint                 # 运行所有 linter（允许失败）
make lint-strict          # 严格检查（失败时退出）
make lint-fast            # 快速检查（只运行快速 linter）

# 修复和报告
make lint-fix             # 自动修复可修复的问题
make lint-report          # 生成 checkstyle 格式报告

# 特定目录检查
make lint-dir DIR=./path  # 检查特定目录

# 安装工具
make install-lint         # 安装 golangci-lint
```

### 3. 自动修复效果
- 运行 `make lint-fix` 后，问题从 30 个减少到 25 个
- 自动修复了 5 个可以自动修复的问题
- 主要是格式化和简单的代码优化问题

## 📊 当前问题统计

| 类型 | 数量 | 描述 |
|------|------|------|
| errcheck | 6 | 未检查的错误返回值 |
| govet | 3 | 格式字符串问题 |
| staticcheck | 15 | 代码质量问题 |
| unused | 1 | 未使用的字段 |
| **总计** | **25** | **需要手动修复** |

## 🔧 需要手动修复的问题

### 1. errcheck (6个) - 高优先级
未检查的错误返回值，可能导致错误被忽略：

```go
// 需要修复的文件：
- internal/log/elasticsearch_sink.go:77
- internal/log/logger.go:55
- internal/log/logger_test.go:111,145
- internal/log/loki_sink.go:89
- muban/modgen/utils/file_utils.go:37
```

### 2. govet (3个) - 中优先级
格式字符串问题，可能导致运行时错误：

```go
// 需要修复的文件：
- muban/modgen/templates/openapi_templates.go:261
- muban/modgen/templates/test_templates.go:183,550
```

### 3. staticcheck (15个) - 中优先级
代码质量问题，包括：
- 废弃的 API 使用
- 空的代码分支
- 不安全的 context key 使用
- 不必要的代码

### 4. unused (1个) - 低优先级
未使用的字段：

```go
// internal/log/logger.go:64
mu   sync.RWMutex  // 未使用
```

## 🚀 使用建议

### 开发流程
1. **开发前**: `make lint-fast` - 快速检查
2. **开发中**: `make lint-fix` - 自动修复
3. **提交前**: `make lint-strict` - 严格检查
4. **定期**: `make lint-report` - 生成报告

### CI/CD 集成
```yaml
- name: Run golangci-lint
  run: make lint-strict
```

### 优先级修复顺序
1. **errcheck** - 修复错误处理
2. **govet** - 修复格式字符串问题
3. **staticcheck** - 修复代码质量问题
4. **unused** - 清理未使用的代码

## 📚 相关文档

- [lint.md](./lint.md) - 详细使用指南
- [golangci-lint 官方文档](https://golangci-lint.run/)
- [Go 代码审查指南](https://github.com/golang/go/wiki/CodeReviewComments)

## 🎯 下一步

1. 修复 errcheck 问题（错误处理）
2. 修复 govet 问题（格式字符串）
3. 逐步修复 staticcheck 问题
4. 清理未使用的代码
5. 考虑添加 pre-commit hook 自动运行 lint
