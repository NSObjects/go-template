### Api 项目实践流程

---

#### 代码编写从最外层向最里层实现 路由->用例->数据层->数据结构，具体流程如下

* 路由层文件 - delivery
* 用例层文件 - usecase
* 数据层文件 - repository
* 数据结构层文件 - domain

```mermaid
graph TD;
api设计-->api文档编写;
api文档编写-->与前端/接口使用者讨论接口数据;

编写路由测试-->编写路由代码;
编写路由代码-->运行测试;
运行测试-->|通过|提交;
运行测试-->|不通过|检查;

定义用例层接口-->测试用例层接口;
测试用例层接口-->|通过|编写接口实现;

定义数据层接口-->测试数据层接口;
测试数据层接口-->|通过|编写接口实现;
  
定义数据结构层接口/结构-->测试数据结构层接口/结构;
测试数据结构层接口/结构-->|通过|编写接口实现;

编写接口实现-->|将Mock对象替换为实现的对象|提交;
```

### mock 文档

https://github.com/golang/mock

### sqlmock 文档

https://github.com/DATA-DOG/go-sqlmock

### echo 对http接口测试文档

https://echo.labstack.com/guide/testing

```
.
├── Dockerfile
├── Makefile
├── README.md
├── app                         ----程序入口
│   └── main.go
├── config.toml                 ----配置文件
├── delivery                    ----展示层 (路由)
│   ├── middlewares             ----中间件
│   │   ├── middleware.go
│   │   └── middleware_test.go
│   ├── server                 
│   │   ├── api.go
│   │   ├── api_test.go
│   │   └── echo_server.go
│   ├── user.go
│   └── user_test.go
├── doc                         ----文档
│   └── openapi.yaml
├── docker-compose.yaml
├── domain                      ----数据结构
│   ├── model.go
│   └── user.go
├── go.mod
├── go.sum
├── repository                  ----数据操作
│   ├── Makefile
│   ├── mock
│   │   └── user.go
│   ├── repository.go
│   ├── user.go
│   └── user_test.go
├── tools                       ----工具函数封装
│   ├── api_error.go
│   ├── api_error_test.go
│   ├── configs                 ----配置
│   │   ├── config.go
│   │   └── config_test.go
│   ├── db                      ----数据库封装
│   │   ├── db.go
│   │   └── db_test.go
│   ├── log                     ----日志
│   │   ├── log.go
│   │   └── log_test.go
│   ├── response.go
│   ├── response_test.go
│   ├── status_code.go
│   ├── tools.go
│   └── tools_test.go
└── usecase                     ----用例层
    ├── Makefile
    ├── mock
    │   └── user.go
    ├── user.go
    ├── user_test.go
    └── usercase.go
```

## 第三库选择

### goroutine pool

https://github.com/panjf2000/ants/

### 同步mysql数据到elastic

https://github.com/siddontang/go-mysql-elasticsearch/

### 命令行工具

https://github.com/spf13/cobra/

### 爬虫工具

https://github.com/gocolly/colly/

### kingshard是一个由Go开发高性能MySQL Proxy
#### 项目，kingshard在满足基本的读写分离的功能上，致力于简化MySQL分库分表操作；能够让DBA通过kingshard轻松平滑地实现MySQL数据库扩容。 kingshard的性能是直连MySQL性能的80%以上

https://github.com/flike/kingshard/blob/master/README_ZH.md

