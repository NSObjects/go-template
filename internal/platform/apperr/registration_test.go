package apperr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func TestAllErrorCodesHaveDefinitions(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to resolve current file path")
	}

	fset := token.NewFileSet()
	pkgDir := filepath.Dir(currentFile)
	files, err := parseProductionFiles(fset, pkgDir)
	if err != nil {
		t.Fatalf("parse package: %v", err)
	}

	defined := collectErrorCodeConstants(files)
	registered := collectDefinitionKeys(files)
	missing := make([]string, 0)
	for name := range defined {
		if _, ok := registered[name]; !ok {
			missing = append(missing, name)
		}
	}

	sort.Strings(missing)
	if len(missing) > 0 {
		t.Fatalf("unregistered error codes: %s", strings.Join(missing, ", "))
	}
}

func parseProductionFiles(fset *token.FileSet, pkgDir string) ([]*ast.File, error) {
	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return nil, fmt.Errorf("read package dir: %w", err)
	}

	files := make([]*ast.File, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !productionGoFile(name) {
			continue
		}
		path := filepath.Join(pkgDir, name)
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", name, err)
		}
		if file.Name.Name != "apperr" {
			return nil, fmt.Errorf("parse %s: package = %s, want apperr", name, file.Name.Name)
		}
		files = append(files, file)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("package apperr not found in %s", pkgDir)
	}
	return files, nil
}

func productionGoFile(name string) bool {
	return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
}

func collectErrorCodeConstants(files []*ast.File) map[string]struct{} {
	defined := map[string]struct{}{}
	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			decl, ok := n.(*ast.GenDecl)
			if !ok || decl.Tok != token.CONST {
				return true
			}
			for _, spec := range decl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range valueSpec.Names {
					if strings.HasPrefix(name.Name, "Err") {
						defined[name.Name] = struct{}{}
					}
				}
			}
			return false
		})
	}
	return defined
}

func collectDefinitionKeys(files []*ast.File) map[string]struct{} {
	registered := map[string]struct{}{}
	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			composite, ok := n.(*ast.CompositeLit)
			if !ok || !isDefinitionsLiteral(composite.Type) {
				return true
			}
			for _, element := range composite.Elts {
				keyValue, ok := element.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				if ident, ok := keyValue.Key.(*ast.Ident); ok {
					registered[ident.Name] = struct{}{}
				}
			}
			return false
		})
	}
	return registered
}

func isDefinitionsLiteral(expr ast.Expr) bool {
	mapType, ok := expr.(*ast.MapType)
	if !ok {
		return false
	}
	key, ok := mapType.Key.(*ast.Ident)
	if !ok || key.Name != "int" {
		return false
	}
	value, ok := mapType.Value.(*ast.Ident)
	return ok && value.Name == "Definition"
}
