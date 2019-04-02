### Api 项目实践流程

---

#### 代码编写从最外层向最里层实现 路由->用例->数据层->数据结构，具体流程如下

* 路由层文件 - api
* 用例层文件 - ucase
* 数据层文件 - repository
* 数据结构层文件 - models

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


<<<<<<< Updated upstream
=======
```
.
├── Makefile
├── README.md
├── apis                            ---- 展示层 (路由)
│   ├── api.go
│   ├── api_helper                  ---- api辅助包
│   │   ├── api_error.go            ---- error封装
│   │   ├── response.go             ---- 返回值封装
│   │   └── status_code.go          ---- 状态码封装
│   └── api_test.go
├── cmd
│   └── api                         ---- api 入口
│       ├── Dockerfile
│       └── main.go
├── config.toml                     ---- 配置文件
├── configs
│   └── config.go
├── doc
│   └── openapi.yaml
├── docker-compose.yaml
├── go.mod
├── go.sum
├── init
│   └── init.go
├── models                          ---- 模型层（ Models ） 
│   └── model.go
├── repository                      ---- 仓库层（ Repository )
│   └── repository.go
├── tools
│   ├── db
│   │   └── db.go                   ---- 数据库
│   ├── log
│   │   └── log.go                  ---- 日志管理
│   └── tools.go
├── ucase                           ---- 用例层 ( Usecase )
    └── usercase.go
```

## 第三库选择

### goroutine pool

https://github.com/panjf2000/ants/

### 常用数据结构go实现参考

https://github.com/emirpasic/gods

https://github.com/Workiva/go-datastructures/

### 同步mysql数据到elastic

https://github.com/siddontang/go-mysql-elasticsearch/

### 命令行工具

https://github.com/spf13/cobra/

### 爬虫工具

https://github.com/gocolly/colly/

### kingshard是一个由Go开发高性能MySQL Proxy
#### 项目，kingshard在满足基本的读写分离的功能上，致力于简化MySQL分库分表操作；能够让DBA通过kingshard轻松平滑地实现MySQL数据库扩容。 kingshard的性能是直连MySQL性能的80%以上

https://github.com/flike/kingshard/blob/master/README_ZH.md
>>>>>>> Stashed changes
