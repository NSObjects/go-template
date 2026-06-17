# Server

`internal/platform/server` 只负责 Echo HTTP runtime：

- 创建 Echo 实例；
- 注册基础 middleware；
- 注册 `/api/health`、`/api/info`；
- 提供业务路由入口 `API()`；
- 处理 HTTP 错误响应；
- 管理启动和优雅关闭。

`/api/info` 只返回静态配置里的 `app.name`、`app.version` 和当前时间。

它不 import 业务模块。业务路由只在 `internal/boot` 里显式组装。

## API

```go
srv, err := server.New(cfg)
if err != nil {
	return err
}

if err := srv.Run(ctx); err != nil {
	return err
}
```

`server.New` 只接收静态配置值，不接收动态配置 store。业务路由不要在这里手写注册；
在 `internal/boot` 里用 `boot.NewModule`、`boot.Provide` 和 `boot.Route` 声明
adapter、usecase、handler 和 route，让 boot 持有同一个 `do` dependency graph。
