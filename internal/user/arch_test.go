package user

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

// TestArchitectureConstraints 测试架构约束
func TestArchitectureConstraints(t *testing.T) {
	tests := []struct {
		name           string
		dir            string
		allowedDeps    []string
		prohibitedDeps []string
	}{
		{
			name: "domain layer should not depend on adapters or application",
			dir:  "domain",
			allowedDeps: []string{
				"github.com/NSObjects/go-template/internal/shared",
			},
			prohibitedDeps: []string{
				"github.com/NSObjects/go-template/internal/user/adapters",
				"github.com/NSObjects/go-template/internal/user/app",
			},
		},
		{
			name: "application layer should not depend on adapters",
			dir:  "app",
			allowedDeps: []string{
				"github.com/NSObjects/go-template/internal/user/domain",
				"github.com/NSObjects/go-template/internal/shared",
			},
			prohibitedDeps: []string{
				"github.com/NSObjects/go-template/internal/user/adapters",
			},
		},
		{
			name: "adapters can depend on domain, application, and shared",
			dir:  "adapters",
			allowedDeps: []string{
				"github.com/NSObjects/go-template/internal/user/domain",
				"github.com/NSObjects/go-template/internal/user/app",
				"github.com/NSObjects/go-template/internal/shared",
				"github.com/NSObjects/go-template/internal/shared/ports/resp",
				"github.com/NSObjects/go-template/internal/pkg/code",
				"github.com/NSObjects/go-template/internal/pkg/utils",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkDependencies(t, tt.dir, tt.allowedDeps, tt.prohibitedDeps)
		})
	}
}

// checkDependencies 检查依赖关系
func checkDependencies(t *testing.T, dir string, allowedDeps, prohibitedDeps []string) {
	// 获取目录下的所有Go文件
	files, err := filepath.Glob(filepath.Join("internal/user", dir, "*.go"))
	if err != nil {
		t.Fatalf("Failed to glob files: %v", err)
	}

	for _, file := range files {
		// 跳过测试文件
		if strings.HasSuffix(file, "_test.go") {
			continue
		}

		// 解析Go文件
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			t.Errorf("Failed to parse file %s: %v", file, err)
			continue
		}

		// 检查导入
		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, "\"")

			// 检查是否包含禁止的依赖
			for _, prohibited := range prohibitedDeps {
				if strings.Contains(importPath, prohibited) {
					t.Errorf("File %s imports prohibited dependency: %s", file, importPath)
				}
			}

			// 检查是否只包含允许的依赖（如果指定了允许的依赖）
			if len(allowedDeps) > 0 {
				allowed := false
				for _, allowedDep := range allowedDeps {
					if strings.Contains(importPath, allowedDep) {
						allowed = true
						break
					}
				}

				// 允许标准库和第三方库
				if !allowed && !isStandardLibrary(importPath) && !isThirdPartyLibrary(importPath) {
					t.Errorf("File %s imports unexpected dependency: %s", file, importPath)
				}
			}
		}
	}
}

// isStandardLibrary 检查是否为标准库
func isStandardLibrary(importPath string) bool {
	standardLibs := []string{
		"context", "fmt", "time", "strings", "strconv", "errors",
		"regexp", "encoding/json", "net/http", "os", "io", "log",
		"math", "sort", "sync", "testing", "reflect", "runtime",
	}

	for _, lib := range standardLibs {
		if importPath == lib {
			return true
		}
	}

	return false
}

// isThirdPartyLibrary 检查是否为第三方库
func isThirdPartyLibrary(importPath string) bool {
	thirdPartyPrefixes := []string{
		"github.com/", "golang.org/", "gopkg.in/", "go.uber.org/",
		"gorm.io/", "github.com/labstack/", "github.com/go-redis/",
	}

	for _, prefix := range thirdPartyPrefixes {
		if strings.HasPrefix(importPath, prefix) {
			return true
		}
	}

	return false
}
