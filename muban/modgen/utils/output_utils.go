/*
 * 输出工具函数
 */

package utils

import "fmt"

// PrintInfo 打印信息
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf(BLUE+"[INFO]"+NC+" "+format+"\n", args...)
}

// PrintSuccess 打印成功信息
func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf(GREEN+"[SUCCESS]"+NC+" "+format+"\n", args...)
}

// ... existing code ...

// PrintError 打印错误信息
func PrintError(format string, args ...interface{}) {
	fmt.Printf(RED+"[ERROR]"+NC+" "+format+"\n", args...)
}

// PrintUsageInstructions 打印使用说明
func PrintUsageInstructions(name, pascal string) {
	fmt.Printf("\n📖 %s 模块使用说明:\n", name)
	fmt.Println("1. 参数结构: internal/api/service/param/" + name + ".go")
	fmt.Println("2. 业务逻辑: internal/api/biz/" + name + ".go")
	fmt.Println("3. 控制器: internal/api/service/" + name + ".go")
	fmt.Println("4. 数据模型: internal/api/data/model/" + name + ".go")
	fmt.Println("5. 错误码: internal/code/" + name + ".go")
	fmt.Println("\n🔧 下一步操作:")
	fmt.Println("1. 根据业务需求修改参数结构和数据模型")
	fmt.Println("2. 实现具体的业务逻辑")
	fmt.Println("3. 配置路由和中间件")
	fmt.Println("4. 运行 'make gen-code' 生成错误码文档")
	fmt.Println("5. 运行 'make run' 启动服务")
	fmt.Println("\n💡 提示:")
	fmt.Printf("- 如未自动注册，请手动将 New%[1]sHandler 和 AsRoute(New%[1]sController) 加入 fx.Options\n", pascal)
	fmt.Println("- 使用 'make db-gen' 生成数据库模型")
	fmt.Println("- 使用 'make gen-code' 生成错误码文档")
}
