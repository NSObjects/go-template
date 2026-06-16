# 配置说明

项目默认只加载一个静态配置文件：

```go
cfg, err := configs.Load("configs/config.toml")
```

配置文件支持 TOML、YAML、JSON，格式由文件后缀识别。环境变量可以覆盖同名配置项。

## 当前配置项

```toml
[system]
port = ":9322"
level = 1 # 1=debug, 2=online
env = "dev"

[jwt]
enabled = false
secret = ""
expire = 3600
skip_paths = ["/api/health", "/api/info"]
```

## 环境变量覆盖

```bash
export ECHOADMIN_SYSTEM_PORT=:8080
export ECHOADMIN_SYSTEM_LEVEL=1
export ECHOADMIN_SYSTEM_ENV=dev
export ECHOADMIN_JWT_ENABLED=false
export ECHOADMIN_JWT_SECRET=
export ECHOADMIN_JWT_EXPIRE=3600
```

JWT 默认不启用，且默认 secret 为空。业务启用 JWT 时必须显式提供 secret。
