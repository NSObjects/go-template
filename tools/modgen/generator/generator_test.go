/*
 * 代码生成器集成测试
 */

package generator

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestEnvironment 设置测试环境，切换到项目根目录
func setupTestEnvironment(t *testing.T) func() {
	originalDir, _ := os.Getwd()
	projectRoot := "/Users/lintao/Documents/code/go-template"
	os.Chdir(projectRoot)

	return func() {
		os.Chdir(originalDir)
	}
}

func TestNewGenerator(t *testing.T) {
	config := &Config{
		Name:        "test",
		Route:       "/tests",
		Force:       true,
		PackagePath: "github.com/test/project",
		RepoRoot:    "/tmp/test-repo",
	}

	generator := NewGenerator(config)
	if generator == nil {
		t.Fatal("生成器不应为 nil")
	}

	if generator.config != config {
		t.Error("配置未正确设置")
	}
}

func TestGenerateFromDefaultTemplate(t *testing.T) {
	// 设置测试环境
	defer setupTestEnvironment(t)()

	// 创建临时目录
	tempDir := t.TempDir()

	config := &Config{
		Name:        "test",
		Route:       "/tests",
		Force:       true,
		PackagePath: "github.com/test/project",
		RepoRoot:    tempDir,
	}

	generator := NewGenerator(config)

	// 创建必要的目录结构
	dirs := []string{
		"internal/api/biz",
		"internal/api/service/param",
		"internal/code",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}
	}

	// 执行生成
	err := generator.Generate()
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}

	// 验证生成的文件
	expectedFiles := []string{
		"internal/api/biz/test.go",
		"internal/api/service/test.go",
		"internal/api/service/param/test.go",
		"internal/code/test.go",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(tempDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("文件未生成: %s", file)
		}
	}
}

func TestGenerateWithTests(t *testing.T) {
	// 设置测试环境
	defer setupTestEnvironment(t)()

	// 创建临时目录
	tempDir := t.TempDir()

	config := &Config{
		Name:          "test",
		Route:         "/tests",
		Force:         true,
		PackagePath:   "github.com/test/project",
		RepoRoot:      tempDir,
		GenerateTests: true,
	}

	generator := NewGenerator(config)

	// 创建必要的目录结构
	dirs := []string{
		"internal/api/biz",
		"internal/api/service/param",
		"internal/code",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}
	}

	// 执行生成
	err := generator.Generate()
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}

	// 验证测试文件
	testFiles := []string{
		"internal/api/biz/test_test.go",
		"internal/api/service/test_test.go",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("测试文件未生成: %s", file)
		}
	}
}

func TestGenerateWithSpecialNames(t *testing.T) {
	testCases := []struct {
		name     string
		route    string
		expected string
	}{
		{"user_profile", "/user-profiles", "UserProfile"},
		{"order_item", "/order-items", "OrderItem"},
		{"product", "/products", "Product"},
		{"test_module", "/test-modules", "TestModule"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置测试环境
			defer setupTestEnvironment(t)()

			// 创建临时目录
			tempDir := t.TempDir()

			config := &Config{
				Name:        tc.name,
				Route:       tc.route,
				Force:       true,
				PackagePath: "github.com/test/project",
				RepoRoot:    tempDir,
			}

			generator := NewGenerator(config)

			// 创建必要的目录结构
			dirs := []string{
				"internal/api/biz",
				"internal/api/service/param",
				"internal/code",
			}

			for _, dir := range dirs {
				err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
				if err != nil {
					t.Fatalf("创建目录失败: %v", err)
				}
			}

			// 执行生成
			err := generator.Generate()
			if err != nil {
				t.Fatalf("生成失败: %v", err)
			}

			// 验证生成的文件内容
			bizFile := filepath.Join(tempDir, "internal/api/biz", tc.name+".go")
			content, err := os.ReadFile(bizFile)
			if err != nil {
				t.Fatalf("读取文件失败: %v", err)
			}

			contentStr := string(content)
			if !contains(contentStr, tc.expected+"UseCase") {
				t.Errorf("缺少 %sUseCase", tc.expected)
			}
			if !contains(contentStr, tc.expected+"Handler") {
				t.Errorf("缺少 %sHandler", tc.expected)
			}
		})
	}
}

func TestGenerateErrorHandling(t *testing.T) {
	// 设置测试环境
	defer setupTestEnvironment(t)()

	// 测试无效的模块名
	config := &Config{
		Name:        "", // 空模块名
		Route:       "/tests",
		Force:       true,
		PackagePath: "github.com/test/project",
		RepoRoot:    "/tmp/test-repo",
	}

	generator := NewGenerator(config)
	err := generator.Generate()
	if err == nil {
		t.Error("应该返回错误")
	}
}

func TestGenerateForceOverwrite(t *testing.T) {
	// 设置测试环境
	defer setupTestEnvironment(t)()

	// 创建临时目录
	tempDir := t.TempDir()

	config := &Config{
		Name:        "test",
		Route:       "/tests",
		Force:       true,
		PackagePath: "github.com/test/project",
		RepoRoot:    tempDir,
	}

	generator := NewGenerator(config)

	// 创建必要的目录结构
	dirs := []string{
		"internal/api/biz",
		"internal/api/service/param",
		"internal/code",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}
	}

	// 第一次生成
	err := generator.Generate()
	if err != nil {
		t.Fatalf("第一次生成失败: %v", err)
	}

	// 第二次生成（应该覆盖）
	err = generator.Generate()
	if err != nil {
		t.Fatalf("第二次生成失败: %v", err)
	}

	// 验证文件存在
	bizFile := filepath.Join(tempDir, "internal/api/biz/test.go")
	if _, err := os.Stat(bizFile); os.IsNotExist(err) {
		t.Error("文件未生成")
	}
}

func TestGenerateWithoutForce(t *testing.T) {
	// 设置测试环境
	defer setupTestEnvironment(t)()

	// 创建临时目录
	tempDir := t.TempDir()

	config := &Config{
		Name:        "test",
		Route:       "/tests",
		Force:       false, // 不强制覆盖
		PackagePath: "github.com/test/project",
		RepoRoot:    tempDir,
	}

	generator := NewGenerator(config)

	// 创建必要的目录结构
	dirs := []string{
		"internal/api/biz",
		"internal/api/service/param",
		"internal/code",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}
	}

	// 第一次生成
	err := generator.Generate()
	if err != nil {
		t.Fatalf("第一次生成失败: %v", err)
	}

	// 第二次生成（应该跳过已存在的文件）
	err = generator.Generate()
	if err != nil {
		t.Fatalf("第二次生成失败: %v", err)
	}

	// 验证文件存在
	bizFile := filepath.Join(tempDir, "internal/api/biz/test.go")
	if _, err := os.Stat(bizFile); os.IsNotExist(err) {
		t.Error("文件未生成")
	}
}

func TestGenerateDefaultRoute(t *testing.T) {
	// 设置测试环境
	defer setupTestEnvironment(t)()

	// 创建临时目录
	tempDir := t.TempDir()

	config := &Config{
		Name:        "test",
		Route:       "", // 空路由，应该使用默认值
		Force:       true,
		PackagePath: "github.com/test/project",
		RepoRoot:    tempDir,
	}

	generator := NewGenerator(config)

	// 创建必要的目录结构
	dirs := []string{
		"internal/api/biz",
		"internal/api/service/param",
		"internal/code",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}
	}

	// 执行生成
	err := generator.Generate()
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}

	// 验证服务文件中的路由
	svcFile := filepath.Join(tempDir, "internal/api/service/test.go")
	content, err := os.ReadFile(svcFile)
	if err != nil {
		t.Fatalf("读取文件失败: %v", err)
	}

	contentStr := string(content)
	// 应该使用默认路由 /tests (复数形式)
	if !contains(contentStr, "g.GET(\"/tests\"") {
		t.Error("应该使用默认路由 /tests")
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
