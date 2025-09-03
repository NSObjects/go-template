# 配置文件说明

## 环境配置

项目使用单一配置文件 `config.toml`，通过 `env` 参数来区分不同环境：

### 环境类型

- **dev**: 开发环境
- **test**: 测试环境  
- **prod**: 生产环境

### 环境配置示例

```toml
[system]
env = "dev"    # 开发环境
# env = "test" # 测试环境
# env = "prod" # 生产环境
```

## 日志配置

### 开发环境 (env = "dev")
- **控制台输出**: 彩色格式，便于开发调试
- **文件输出**: 禁用，只输出到终端
- **日志级别**: debug，显示所有日志

### 测试环境 (env = "test")
- **控制台输出**: 彩色格式
- **文件输出**: 启用，JSON格式
- **日志级别**: info，显示重要日志

### 生产环境 (env = "prod")
- **控制台输出**: JSON格式
- **文件输出**: 启用，JSON格式
- **日志级别**: info，只显示重要日志

## 部署配置

### 开发环境部署
```bash
# 修改配置文件
[system]
env = "dev"

# 启动服务
make run
```

### 生产环境部署
```bash
# 修改配置文件
[system]
env = "prod"

# 启动服务
make run
```

## 配置参数说明

### 系统配置
- `port`: 服务端口
- `level`: 系统级别 (1=debug, 2=online)
- `env`: 环境类型 (dev/test/prod)

### 日志配置
- `level`: 日志级别 (debug/info/warn/error)
- `format`: 日志格式 (color/json/text)
- `console.format`: 控制台输出格式
- `file.filename`: 文件输出路径
- `file.max_size`: 文件最大大小(MB)
- `file.max_backups`: 最大备份数量
- `file.max_age`: 文件保留天数

### 数据库配置
- `mysql.host`: MySQL主机地址
- `mysql.port`: MySQL端口
- `mysql.database`: 数据库名称
- `mysql.user`: 用户名
- `mysql.password`: 密码

### 安全配置
- `jwt.secret`: JWT签名密钥
- `jwt.expire`: JWT过期时间(秒)
- `jwt.skip_paths`: 跳过JWT验证的路径

## 注意事项

1. **生产环境**: 请修改默认密码和密钥
2. **日志文件**: 生产环境建议输出到 `/var/log/` 目录
3. **数据库**: 生产环境请使用强密码
4. **JWT密钥**: 生产环境请使用复杂的随机密钥
5. **CORS**: 生产环境请限制允许的域名

## 环境变量覆盖

可以通过环境变量覆盖配置文件中的值：

```bash
# 设置环境
export ECHO_ADMIN_ENV=prod

# 设置端口
export ECHO_ADMIN_PORT=8080

# 启动服务
make run
```
