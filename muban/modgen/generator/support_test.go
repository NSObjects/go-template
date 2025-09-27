package generator

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnsureContextSupportCreatesSupportFileWhenMissing(t *testing.T) {
	tmpDir := t.TempDir()
	utilsDir := filepath.Join(tmpDir, "internal", "utils")
	require.NoError(t, os.MkdirAll(utilsDir, 0o755))

	// 旧项目中可能存在的占位文件
	legacyPath := filepath.Join(utilsDir, "context_trace.go")
	require.NoError(t, os.WriteFile(legacyPath, []byte("package utils\n"), 0o644))

	g := &Generator{config: &Config{RepoRoot: tmpDir}}
	require.NoError(t, g.ensureContextSupport())

	supportPath := filepath.Join(utilsDir, "context_support.go")
	data, err := os.ReadFile(supportPath)
	require.NoError(t, err)
	content := string(data)

	require.Contains(t, content, "func BuildContext", "should render context support helpers")
	require.Contains(t, content, "func ExtractTraceContext")

	// 旧文件保持为空，实现由 context_support.go 提供
	legacyData, err := os.ReadFile(legacyPath)
	require.NoError(t, err)
	require.Equal(t, strings.TrimSpace(legacyContextTraceStub), strings.TrimSpace(string(legacyData)))
}

func TestEnsureContextSupportRewritesLegacyTraceFile(t *testing.T) {
	tmpDir := t.TempDir()
	utilsDir := filepath.Join(tmpDir, "internal", "utils")
	require.NoError(t, os.MkdirAll(utilsDir, 0o755))

	// 已有的 context.go 定义
	contextPath := filepath.Join(utilsDir, "context.go")
	contextSource := `package utils

import "context"

func BuildContext() context.Context { return context.TODO() }
`
	require.NoError(t, os.WriteFile(contextPath, []byte(contextSource), 0o644))

	// 遗留的 context_trace.go 内容
	legacyContent := `package utils

func GetTraceID() {}
func BuildContext(c echo.Context) {}
`
	legacyPath := filepath.Join(utilsDir, "context_trace.go")
	require.NoError(t, os.WriteFile(legacyPath, []byte(legacyContent), 0o644))

	g := &Generator{config: &Config{RepoRoot: tmpDir}}
	require.NoError(t, g.ensureContextSupport())

	_, err := os.Stat(filepath.Join(utilsDir, "context_support.go"))
	require.True(t, errors.Is(err, os.ErrNotExist), "should not create support file when helpers exist")

	legacyData, err := os.ReadFile(legacyPath)
	require.NoError(t, err)
	require.Equal(t, strings.TrimSpace(legacyContextTraceStub), strings.TrimSpace(string(legacyData)))
}


