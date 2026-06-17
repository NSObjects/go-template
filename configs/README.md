# 配置说明

项目默认只加载一个静态配置文件：

```go
cfg, err := configs.Load("configs/config.toml")
```

配置文件支持 TOML、YAML、JSON，格式由文件后缀识别。无后缀时按 TOML 解析，未知后缀会在启动时失败。环境变量可以覆盖同名配置项。

## 当前配置项

```toml
[system]
port = ":9322"
level = 1 # 1=debug, 2=online
env = "dev"

[jwt]
enabled = false
secret = ""
skip_paths = ["/api/health", "/api/info"]
```

## 环境变量覆盖

```bash
export GO_TEMPLATE_SYSTEM_PORT=:9322
export GO_TEMPLATE_SYSTEM_LEVEL=1
export GO_TEMPLATE_SYSTEM_ENV=dev
export GO_TEMPLATE_JWT_ENABLED=false
export GO_TEMPLATE_JWT_SECRET=
```

JWT 默认不启用，且默认 secret 为空。业务启用 JWT 时必须显式提供 secret。
