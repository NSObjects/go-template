package main

import (
	"fmt"

	"github.com/NSObjects/echo-admin/internal/api/data"
	"github.com/NSObjects/echo-admin/internal/api/data/query"
	"github.com/NSObjects/echo-admin/internal/configs"
	"go.uber.org/fx"
)

func main() {
	// 创建FX应用
	app := fx.New(
		// 配置模块
		fx.Provide(func() configs.Config {
			return configs.NewCfg("configs/config.toml")
		}),

		// 数据模块
		data.Model,

		// 启动函数
		fx.Invoke(run),
	)

	// 运行应用
	app.Run()
}

func run(dm *data.DataManager, q *query.Query) {
	fmt.Println("=== Dynamic SQL 使用示例 ===")

	// 检查组件状态
	if !dm.IsComponentEnabled("mysql") {
		fmt.Println("MySQL组件未启用，无法演示Dynamic SQL")
		return
	}

	fmt.Println("\n1. 基础CRUD操作示例（通用方法）:")

	// 使用通用查询接口 - 这些方法对所有表都可用
	if q.Available() {
		// 获取用户总数
		count, err := q.User.CountRecords()
		if err != nil {
			fmt.Printf("获取用户总数失败: %v\n", err)
		} else {
			fmt.Printf("用户总数: %d\n", count)
		}

		// 检查用户是否存在
		exists, err := q.User.Exists(1)
		if err != nil {
			fmt.Printf("检查用户存在性失败: %v\n", err)
		} else {
			fmt.Printf("用户ID=1是否存在: %v\n", exists)
		}

		// 获取用户列表
		users, err := q.User.GetByIDs([]uint{1, 2, 3})
		if err != nil {
			fmt.Printf("获取用户列表失败: %v\n", err)
		} else {
			fmt.Printf("获取到 %d 个用户\n", len(users))
		}

		// 演示通用方法：casbin_rule表也有相同的方法
		casbinCount, err := q.CasbinRule.CountRecords()
		if err != nil {
			fmt.Printf("获取Casbin规则总数失败: %v\n", err)
		} else {
			fmt.Printf("Casbin规则总数: %d\n", casbinCount)
		}
	}

	fmt.Println("\n2. 分页和搜索示例（通用方法）:")

	// 使用通用分页和搜索接口
	if q.Available() {
		// 分页获取用户
		users, err := q.User.GetPage(0, 10, "id ASC")
		if err != nil {
			fmt.Printf("分页获取用户失败: %v\n", err)
		} else {
			fmt.Printf("分页获取到 %d 个用户\n", len(users))
		}

		// 搜索用户
		searchResults, err := q.User.Search("username", "admin")
		if err != nil {
			fmt.Printf("搜索用户失败: %v\n", err)
		} else {
			fmt.Printf("搜索到 %d 个用户\n", len(searchResults))
		}

		// 演示通用方法：casbin_rule表也有相同的分页和搜索方法
		casbinRules, err := q.CasbinRule.GetPage(0, 5, "id ASC")
		if err != nil {
			fmt.Printf("分页获取Casbin规则失败: %v\n", err)
		} else {
			fmt.Printf("分页获取到 %d 个Casbin规则\n", len(casbinRules))
		}
	}

	fmt.Println("\n3. 高级查询示例（模板表达式）:")

	// 使用高级查询接口，支持模板表达式
	if q.Available() {
		// 使用条件过滤
		filteredUsers, err := q.User.FilterWithCondition("status = 1")
		if err != nil {
			fmt.Printf("条件过滤用户失败: %v\n", err)
		} else {
			fmt.Printf("条件过滤到 %d 个用户\n", len(filteredUsers))
		}

		// 使用通用字段查询
		userByField, err := q.User.GetByField("username", "admin")
		if err != nil {
			fmt.Printf("根据字段查询用户失败: %v\n", err)
		} else {
			fmt.Printf("根据字段查询到用户: %+v\n", userByField)
		}

		// 演示通用方法：casbin_rule表也有相同的高级查询方法
		casbinByField, err := q.CasbinRule.GetByField("ptype", "p")
		if err != nil {
			fmt.Printf("根据字段查询Casbin规则失败: %v\n", err)
		} else {
			fmt.Printf("根据字段查询到Casbin规则: %+v\n", casbinByField)
		}
	}

	fmt.Println("\n=== Dynamic SQL 示例完成 ===")
	fmt.Println("\n核心特性:")
	fmt.Println("1. 通用方法：所有表都生成相同的查询方法")
	fmt.Println("2. 类型安全：所有生成的代码都是类型安全的，编译时检查")
	fmt.Println("3. 模板表达式：支持 if/else, where, set, for 等高级功能")
	fmt.Println("4. 占位符：@@table 自动替换为表名，@param 绑定参数")
	fmt.Println("\n使用说明:")
	fmt.Println("1. 使用 'make db-gen-dynamic' 生成Dynamic SQL查询")
	fmt.Println("2. 使用 'make db-gen-full' 生成完整的数据库代码")
	fmt.Println("3. 在业务代码中通过 query.Q 或注入的 *query.Query 使用")
	fmt.Println("4. 参考官方文档：https://gorm.io/zh_CN/gen/dynamic_sql.html")
}
