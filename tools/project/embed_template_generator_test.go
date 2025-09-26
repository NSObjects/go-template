package project

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateProjectCreatesExecutableScripts(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	generator := NewEmbedTemplateGenerator(tempDir, "github.com/example/demo", "demo")
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
}

func TestGenerateProjectWithCustomLoggerAndModes(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
