package project

import (
	"fmt"
	"os"
	"path/filepath"
)

// TestEmbedGeneration 测试嵌入模板生成
func TestEmbedGeneration() {
	// 创建测试目录
	testDir := "/tmp/local-embed-test"
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0755)

	fmt.Printf("🚀 开始测试嵌入模板生成器\n")
	fmt.Printf("📁 测试目录: %s\n", testDir)

	// 创建嵌入模板生成器
	generator := NewEmbedTemplateGenerator(testDir, "github.com/test/local-embed", "local-embed-test")

	// 生成项目
	if err := generator.Generate(); err != nil {
		fmt.Printf("❌ 生成失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 项目生成完成！\n")

	// 检查生成的文件
	keyFiles := []string{
		"main.go",
		"go.mod",
		"README.md",
	}

	successCount := 0
	for _, file := range keyFiles {
		fullPath := filepath.Join(testDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fmt.Printf("❌ 文件不存在: %s\n", file)
		} else {
			fmt.Printf("✅ 文件存在: %s\n", file)
			successCount++
		}
	}

	fmt.Printf("\n📊 生成结果: %d/%d 文件成功生成\n", successCount, len(keyFiles))

	// 检查main.go的内容
	mainPath := filepath.Join(testDir, "main.go")
	if content, err := os.ReadFile(mainPath); err == nil {
		fmt.Printf("\n📄 main.go 内容:\n%s\n", string(content))
	}
}
