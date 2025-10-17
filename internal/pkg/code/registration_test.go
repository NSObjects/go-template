package code

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestAllErrorCodesRegistered(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("unable to resolve current file path")
	}
	pkgDir := filepath.Dir(currentFile)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, func(info fs.FileInfo) bool {
		name := info.Name()
		if strings.HasSuffix(name, "_test.go") {
			return false
		}
		return strings.HasSuffix(name, ".go")
	}, 0)
	if err != nil {
		t.Fatalf("parse dir: %v", err)
	}

	pkg, ok := pkgs["code"]
	if !ok {
		t.Fatalf("package code not found in %s", pkgDir)
	}

	defined := map[string]struct{}{}
	for filename, file := range pkg.Files {
		if strings.HasSuffix(filename, "code_generated.go") {
			continue
		}
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

	genFile := filepath.Join(pkgDir, "code_generated.go")
	generated, err := parser.ParseFile(fset, genFile, nil, 0)
	if err != nil {
		t.Fatalf("parse generated file: %v", err)
	}

	registered := map[string]struct{}{}
	ast.Inspect(generated, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		ident, ok := call.Fun.(*ast.Ident)
		if !ok || ident.Name != "register" {
			return true
		}
		if len(call.Args) == 0 {
			return true
		}
		if argIdent, ok := call.Args[0].(*ast.Ident); ok {
			registered[argIdent.Name] = struct{}{}
		}
		return true
	})

	missing := make([]string, 0)
	for name := range defined {
		if _, ok := registered[name]; !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		t.Fatalf("unregistered error codes: %s", strings.Join(missing, ", "))
	}
}
