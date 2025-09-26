package generator

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGenerateDefaultModuleProducesTestsByDefault(t *testing.T) {
	restore := chdirToProjectRoot(t)
	defer restore()

	repoRoot := t.TempDir()
	createRepoLayout(t, repoRoot)

	cfg := &Config{
		Name:          "book",
		Force:         true,
		PackagePath:   "github.com/example/project",
		RepoRoot:      repoRoot,
		GenerateTests: true,
	}

	gen := NewGenerator(cfg)
	if err := gen.Generate(); err != nil {
		t.Fatalf("generate default module: %v", err)
	}

	assertFileExists(t, filepath.Join(repoRoot, "internal/api/biz/book.go"))
	assertFileExists(t, filepath.Join(repoRoot, "internal/api/service/book.go"))
	assertFileExists(t, filepath.Join(repoRoot, "internal/api/service/param/book.go"))
	assertFileExists(t, filepath.Join(repoRoot, "internal/code/book.go"))
	assertFileExists(t, filepath.Join(repoRoot, "internal/api/biz/book_test.go"))
	assertFileExists(t, filepath.Join(repoRoot, "internal/api/service/book_test.go"))
}

func TestGenerateSkipsTestsWhenDisabled(t *testing.T) {
	restore := chdirToProjectRoot(t)
	defer restore()

	repoRoot := t.TempDir()
	createRepoLayout(t, repoRoot)

	cfg := &Config{
		Name:          "order",
		Force:         true,
		PackagePath:   "github.com/example/project",
		RepoRoot:      repoRoot,
		GenerateTests: false,
	}

	gen := NewGenerator(cfg)
	if err := gen.Generate(); err != nil {
		t.Fatalf("generate module without tests: %v", err)
	}

	ensureFileMissing(t, filepath.Join(repoRoot, "internal/api/biz/order_test.go"))
	ensureFileMissing(t, filepath.Join(repoRoot, "internal/api/service/order_test.go"))
}

func TestGenerateFromOpenAPISpec(t *testing.T) {
	restore := chdirToProjectRoot(t)
	defer restore()

	repoRoot := t.TempDir()
	createRepoLayout(t, repoRoot)

	specPath := writeOpenAPISpec(t, simpleOpenAPISpec())

	cfg := &Config{
		Name:          "user",
		Force:         true,
		PackagePath:   "github.com/example/project",
		RepoRoot:      repoRoot,
		OpenAPIFile:   specPath,
		GenerateTests: true,
	}

	gen := NewGenerator(cfg)
	if err := gen.Generate(); err != nil {
		t.Fatalf("generate module from openapi: %v", err)
	}

	serviceContent := readFile(t, filepath.Join(repoRoot, "internal/api/service/user.go"))
	if !strings.Contains(serviceContent, "g.GET(\"/users\"") {
		t.Fatalf("expected service to register /users route, got: %s", serviceContent)
	}

	bizTests := readFile(t, filepath.Join(repoRoot, "internal/api/biz/user_test.go"))
	if !strings.Contains(bizTests, "table-driven") {
		t.Fatalf("expected generated tests to mention table-driven style")
	}
}

func TestGenerateAllModulesFromOpenAPISpec(t *testing.T) {
	restore := chdirToProjectRoot(t)
	defer restore()

	repoRoot := t.TempDir()
	createRepoLayout(t, repoRoot)

	specPath := writeOpenAPISpec(t, multiModuleOpenAPISpec())

	cfg := &Config{
		Force:         true,
		PackagePath:   "github.com/example/project",
		RepoRoot:      repoRoot,
		OpenAPIFile:   specPath,
		GenerateTests: true,
		GenerateAll:   true,
	}

	gen := NewGenerator(cfg)
	if err := gen.generateAllModules(); err != nil {
		t.Fatalf("generate all modules: %v", err)
	}

	assertFileExists(t, filepath.Join(repoRoot, "internal/api/service/user.go"))
	assertFileExists(t, filepath.Join(repoRoot, "internal/api/service/article.go"))
	assertFileExists(t, filepath.Join(repoRoot, "internal/api/service/article_test.go"))
}

func TestGenerateFailsWhenNameMissing(t *testing.T) {
	restore := chdirToProjectRoot(t)
	defer restore()

	repoRoot := t.TempDir()
	createRepoLayout(t, repoRoot)

	cfg := &Config{
		Force:       true,
		PackagePath: "github.com/example/project",
		RepoRoot:    repoRoot,
	}

	gen := NewGenerator(cfg)
	if err := gen.Generate(); err == nil {
		t.Fatalf("expected error when module name is empty")
	}
}

func chdirToProjectRoot(t *testing.T) func() {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to determine caller path")
	}

	projectRoot, err := filepath.Abs(filepath.Join(filepath.Dir(file), "..", "..", ".."))
	if err != nil {
		t.Fatalf("resolve project root: %v", err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("determine working directory: %v", err)
	}

	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("change directory to project root: %v", err)
	}

	return func() {
		if err := os.Chdir(original); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	}
}

func createRepoLayout(t *testing.T, root string) {
	t.Helper()

	dirs := []string{
		"internal/api/biz",
		"internal/api/service/param",
		"internal/code",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatalf("create directory %s: %v", dir, err)
		}
	}

	writeFile(t, filepath.Join(root, "internal/api/biz/biz.go"), "package biz\n")
	writeFile(t, filepath.Join(root, "internal/api/service/service.go"), "package service\n")
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s to exist: %v", path, err)
	}
}

func ensureFileMissing(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	if err == nil {
		t.Fatalf("expected file %s to be absent", path)
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected error for %s: %v", path, err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %s: %v", path, err)
	}
	return string(data)
}

func writeOpenAPISpec(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "openapi.yaml")
	writeFile(t, path, content)
	return path
}

func simpleOpenAPISpec() string {
	return `openapi: 3.0.0
info:
  title: Simple
  version: 1.0.0
paths:
  /users:
    get:
      tags: [user]
      summary: List users
      responses:
        "200":
          description: ok
`
}

func multiModuleOpenAPISpec() string {
	return `openapi: 3.0.0
info:
  title: Multi
  version: 1.0.0
tags:
  - name: user
  - name: article
paths:
  /users:
    get:
      tags: [user]
      summary: List users
      responses:
        "200":
          description: ok
  /articles:
    post:
      tags: [article]
      summary: Create article
      requestBody:
        content:
          application/json:
            schema:
              type: object
      responses:
        "201":
          description: created
`
}
