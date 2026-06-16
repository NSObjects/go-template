# Server

`internal/server` 只负责 Echo HTTP runtime：

- 创建 Echo 实例；
- 注册基础 middleware；
- 注册 `/api/health`、`/api/info`；
- 提供业务路由入口 `API()`；
- 处理 HTTP 错误响应；
- 管理启动和优雅关闭。

它不 import 业务模块。业务路由只在 `internal/boot` 里显式组装。

## API

```go
srv := server.New(cfg)
orders := srv.API().Group("/orders")
orderhttp.Register(orders, handler)

if err := srv.Run(ctx); err != nil {
	return err
}
```

`server.New` 只接收静态配置值，不接收动态配置 store。业务需要自己的 DB、cache、queue 或外部客户端时，在 boot 中创建 concrete adapter，再传给对应 usecase。
