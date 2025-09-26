package project

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateProjectCreatesExecutableScripts(t *testing.T) {

	tempDir := t.TempDir()

	generator := NewEmbedTemplateGenerator(
		tempDir,
		"github.com/example/demo",
		"demo",
		WithLogger(func(string, ...interface{}) {}),
	)
	if err := generator.Generate(); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	scriptPath := filepath.Join(tempDir, "scripts", "dev.sh")
	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatalf("stat dev.sh: %v", err)
	}

	if mode := info.Mode().Perm(); mode != 0o755 {
		t.Fatalf("unexpected script mode: %v", mode)
	}

	goModBytes, err := os.ReadFile(filepath.Join(tempDir, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if !strings.Contains(string(goModBytes), "module github.com/example/demo") {
		t.Fatalf("go.mod does not contain module path: %s", goModBytes)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "go.sum")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected go.sum to be skipped, got err=%v", err)
	}

	docsDir := filepath.Join(tempDir, "docs")
	if entries, err := os.ReadDir(docsDir); err == nil {
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".md") {
				t.Fatalf("expected markdown docs to be skipped, found %s", entry.Name())
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected error reading docs dir: %v", err)
	}
}

func TestGenerateProjectSanitizesNames(t *testing.T) {

	tempDir := t.TempDir()

	generator := NewEmbedTemplateGenerator(
		tempDir,
		"github.com/example/my-service",
		"",
		WithLogger(func(string, ...interface{}) {}),
	)

	if err := generator.Generate(); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	envBytes, err := os.ReadFile(filepath.Join(tempDir, "env.example"))
	if err != nil {
		t.Fatalf("read env.example: %v", err)
	}
	envContent := string(envBytes)
	if !strings.Contains(envContent, "MYSQL_DATABASE=my_service") {
		t.Fatalf("expected env.example to use snake_case project name, got %q", envContent)
	}

	composeBytes, err := os.ReadFile(filepath.Join(tempDir, "docker-compose.yaml"))
	if err != nil {
		t.Fatalf("read docker-compose.yaml: %v", err)
	}
	compose := string(composeBytes)
	if !strings.Contains(compose, "  mysql:\n    image") {
		t.Fatalf("docker-compose content unexpected: %q", compose)
	}
	if !strings.Contains(compose, "MYSQL_DATABASE: my_service") {
		t.Fatalf("expected docker-compose to use snake_case for database name, got %q", compose)
	}
	if !strings.Contains(compose, "services:\n  my-service:\n") {
		t.Fatalf("expected docker-compose to keep kebab-case service name, got %q", compose)
	}

	readmeBytes, err := os.ReadFile(filepath.Join(tempDir, "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	if !strings.Contains(string(readmeBytes), "# my-service") {
		t.Fatalf("expected README heading to preserve display name, got %q", readmeBytes)
	}
}

func TestGenerateProjectWithCustomLoggerAndModes(t *testing.T) {

	tempDir := t.TempDir()
	var logBuilder strings.Builder

	generator := NewEmbedTemplateGenerator(
		tempDir,
		"github.com/example/custom",
		"custom",
		WithLogger(func(format string, args ...interface{}) {
			fmt.Fprintf(&logBuilder, format, args...)
		}),
		WithModeResolver(func(path string) fs.FileMode {
			if strings.HasSuffix(path, ".sh") {
				return 0o700
			}
			return 0o640
		}),
	)

	if err := generator.Generate(); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	logs := logBuilder.String()
	if !strings.Contains(logs, "✅ 生成文件") {
		t.Fatalf("expected generator to log progress, got %q", logs)
	}

	scriptInfo, err := os.Stat(filepath.Join(tempDir, "scripts", "dev.sh"))
	if err != nil {
		t.Fatalf("stat dev.sh: %v", err)
	}
	if mode := scriptInfo.Mode().Perm(); mode != 0o700 {
		t.Fatalf("unexpected script mode: %v", mode)
	}

	goModInfo, err := os.Stat(filepath.Join(tempDir, "go.mod"))
	if err != nil {
		t.Fatalf("stat go.mod: %v", err)
	}
	if mode := goModInfo.Mode().Perm(); mode != 0o640 {
		t.Fatalf("unexpected go.mod mode: %v", mode)
	}
}

func TestRunWithWriterUsesProvidedOutput(t *testing.T) {

	tempDir := t.TempDir()
	targetDir := filepath.Join(tempDir, "myapp")

	opts := Options{
		ModulePath: "github.com/example/runwriter",
		OutputDir:  targetDir,
	}

	var output strings.Builder
	if err := RunWithWriter(opts, &output); err != nil {
		t.Fatalf("RunWithWriter() error = %v", err)
	}

	if !strings.Contains(output.String(), "✅ 项目生成完成") {
		t.Fatalf("expected completion message in output, got %q", output.String())
	}

	if _, err := os.Stat(filepath.Join(targetDir, "go.mod")); err != nil {
		t.Fatalf("expected go.mod to be generated: %v", err)
	}
}

func TestRunWithWriterRejectsInvalidModulePath(t *testing.T) {
	tempDir := t.TempDir()

	opts := Options{
		ModulePath: "github.com/example/with space",
		OutputDir:  filepath.Join(tempDir, "out"),
	}

	err := RunWithWriter(opts, io.Discard)
	if err == nil {
		t.Fatalf("expected invalid module path to return error")
	}
	if !strings.Contains(err.Error(), "无效的 Go Module 路径") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestRunWithWriterRejectsDangerousForce(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取工作目录失败: %v", err)
	}

	t.Cleanup(func() {
		if chdirErr := os.Chdir(originalWD); chdirErr != nil {
			t.Fatalf("恢复工作目录失败: %v", chdirErr)
		}
	})

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("切换目录失败: %v", err)
	}

	opts := Options{
		ModulePath: "github.com/example/danger",
		OutputDir:  tempDir,
		Force:      true,
	}

	var output strings.Builder
	err = RunWithWriter(opts, &output)
	if err == nil {
		t.Fatalf("期望在危险目录上报错，但未返回错误")
	}
	if !strings.Contains(err.Error(), "拒绝清理") {
		t.Fatalf("返回的错误信息不包含预期提示: %v", err)
	}
}
