package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratorGenerateCopiesTemplate(t *testing.T) {
	templateDir := t.TempDir()

	// Create minimal template structure
	require.NoError(t, os.MkdirAll(filepath.Join(templateDir, "cmd"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "go.mod"), []byte("module github.com/NSObjects/go-template\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "cmd", "root.go"), []byte("package cmd\n\nimport \"github.com/NSObjects/go-template/internal/api\"\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "README.md"), []byte("# go-template\n"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(templateDir, "tools"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "tools", "ignore.txt"), []byte("should be ignored"), 0o644))

	outputDir := filepath.Join(t.TempDir(), "myapp")

	gen, err := NewGenerator(GeneratorConfig{
		SourceRoot:          templateDir,
		OutputDir:           outputDir,
		TemplateModulePath:  "github.com/NSObjects/go-template",
		TargetModulePath:    "github.com/example/myapp",
		TemplateProjectName: "go-template",
		TargetProjectName:   "myapp",
		Force:               true,
	})
	require.NoError(t, err)

	require.NoError(t, gen.Generate())

	// Ensure files copied with replacements
	rootContent, err := os.ReadFile(filepath.Join(outputDir, "cmd", "root.go"))
	require.NoError(t, err)
	content := string(rootContent)
	require.Contains(t, content, "github.com/example/myapp/internal/api")
	require.NotContains(t, content, "github.com/NSObjects/go-template/internal/api")

	readme, err := os.ReadFile(filepath.Join(outputDir, "README.md"))
	require.NoError(t, err)
	require.True(t, strings.Contains(string(readme), "myapp"))

	// Ensure excluded directory skipped
	_, err = os.Stat(filepath.Join(outputDir, "tools", "ignore.txt"))
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}

func TestGeneratorGenerateFailsWhenTargetExists(t *testing.T) {
	templateDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "go.mod"), []byte("module github.com/NSObjects/go-template\n"), 0o644))

	targetDir := filepath.Join(t.TempDir(), "exists")
	require.NoError(t, os.MkdirAll(targetDir, 0o755))

	gen, err := NewGenerator(GeneratorConfig{
		SourceRoot:          templateDir,
		OutputDir:           targetDir,
		TemplateModulePath:  "github.com/NSObjects/go-template",
		TargetModulePath:    "github.com/example/myapp",
		TemplateProjectName: "go-template",
		TargetProjectName:   "myapp",
		Force:               false,
	})
	require.NoError(t, err)

	err = gen.Generate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "目标目录已存在")
}
