package main

import (
	"context"
	"fmt"

	"github.com/NSObjects/echo-admin/internal/api/data"
	"github.com/NSObjects/echo-admin/internal/configs"
	"go.uber.org/fx"
)

// 这是一个模板示例，展示如何使用框架
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

func run(dm *data.DataManager) {
	ctx := context.Background()

	fmt.Println("=== 框架模板示例 ===")

	// 检查组件状态
	fmt.Println("\n1. 检查组件状态:")
	health := dm.Health(ctx)
	for component, err := range health {
		if err != nil {
			fmt.Printf("❌ %s: %v\n", component, err)
		} else {
			fmt.Printf("✅ %s: healthy\n", component)
		}
	}

	// 检查组件启用状态
	fmt.Println("\n2. 组件启用状态:")
	components := []string{"mysql", "redis", "kafka", "mongodb"}
	for _, component := range components {
		enabled := dm.IsComponentEnabled(component)
		status := "❌ 未启用"
		if enabled {
			status = "✅ 已启用"
		}
		fmt.Printf("%s: %s\n", component, status)
	}

	// 使用MySQL示例
	if dm.IsComponentEnabled("mysql") {
		fmt.Println("\n3. MySQL使用示例:")
		var count int64
		if err := dm.MySQLWithContext(ctx).Table("users").Count(&count).Error; err != nil {
			fmt.Printf("MySQL查询失败: %v\n", err)
		} else {
			fmt.Printf("用户总数: %d\n", count)
		}
	}

	// 使用Redis示例
	if dm.IsComponentEnabled("redis") {
		fmt.Println("\n4. Redis使用示例:")
		redisClient := dm.RedisWithContext(ctx)
		if err := redisClient.Set(ctx, "test_key", "test_value", 0).Err(); err != nil {
			fmt.Printf("Redis设置失败: %v\n", err)
		} else {
			val, err := redisClient.Get(ctx, "test_key").Result()
			if err != nil {
				fmt.Printf("Redis获取失败: %v\n", err)
			} else {
				fmt.Printf("Redis值: %s\n", val)
			}
		}
	}

	// 使用Kafka示例
	if dm.IsComponentEnabled("kafka") {
		fmt.Println("\n5. Kafka使用示例:")
		if err := dm.SendKafkaMessage("test-topic", []byte("key"), []byte("value")); err != nil {
			fmt.Printf("Kafka发送失败: %v\n", err)
		} else {
			fmt.Println("Kafka消息发送成功")
		}
	}

	fmt.Println("\n=== 模板示例完成 ===")
	fmt.Println("\n使用说明:")
	fmt.Println("1. 使用 'make gen-module NAME=your_module' 生成新模块")
	fmt.Println("2. 使用 'make db-gen' 生成数据库模型")
	fmt.Println("3. 使用 'make gen-code' 生成错误码")
	fmt.Println("4. 在 biz.go 和 service.go 中注册新模块")
}
