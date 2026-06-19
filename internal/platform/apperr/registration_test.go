package apperr

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
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
	pkgs, err := parser.ParseDir(fset, pkgDir, productionGoFile, 0)
	if err != nil {
		t.Fatalf("parse package: %v", err)
	}

	pkg, ok := pkgs["apperr"]
	if !ok {
		t.Fatalf("package apperr not found in %s", pkgDir)
	}

	defined := collectErrorCodeConstants(pkg)
	registered := collectDefinitionKeys(pkg)
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

func productionGoFile(info fs.FileInfo) bool {
	name := info.Name()
	return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
}

func collectErrorCodeConstants(pkg *ast.Package) map[string]struct{} {
	defined := map[string]struct{}{}
	for _, file := range pkg.Files {
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

func collectDefinitionKeys(pkg *ast.Package) map[string]struct{} {
	registered := map[string]struct{}{}
	for _, file := range pkg.Files {
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
