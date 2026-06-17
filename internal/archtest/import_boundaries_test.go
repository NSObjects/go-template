package archtest

import (
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const modulePath = "github.com/NSObjects/go-template"

type goPackage struct {
	ImportPath string
	Imports    []string
}

func TestCleanLiteImportBoundaries(t *testing.T) {
	packages := listInternalPackages(t)
	businessRoots := findBusinessRoots(packages)

	for _, pkg := range packages {
		switch {
		case isLayerPackage(pkg.ImportPath, "domain"):
			assertDomainImportsStdlibOnly(t, pkg)
			assertNoOtherBusinessImports(t, pkg, businessRoots)
		case isLayerPackage(pkg.ImportPath, "usecase"):
			assertNoForbiddenImports(t, pkg, usecaseForbiddenImports(pkg.ImportPath))
			assertNoOtherBusinessImports(t, pkg, businessRoots)
		case isLayerPackage(pkg.ImportPath, "http"):
			assertNoForbiddenImports(t, pkg, httpForbiddenImports(pkg.ImportPath))
			assertNoOtherBusinessImports(t, pkg, businessRoots)
		case isLayerPackage(pkg.ImportPath, "mysql"):
			assertNoForbiddenImports(t, pkg, mysqlForbiddenImports(pkg.ImportPath))
			assertNoOtherBusinessImports(t, pkg, businessRoots)
		case isInternalRootPackage(pkg.ImportPath, "server"), isInternalRootPackage(pkg.ImportPath, "configs"):
			assertNoBusinessImports(t, pkg, businessRoots)
		}
	}
}

func TestFindBusinessRoots(t *testing.T) {
	packages := []goPackage{
		{ImportPath: modulePath + "/internal/order/usecase"},
		{ImportPath: modulePath + "/internal/order/http"},
		{ImportPath: modulePath + "/internal/server/httpresp"},
		{ImportPath: modulePath + "/internal/mysqlinfra"},
		{ImportPath: modulePath + "/internal/infrastructure/mysql"},
		{ImportPath: modulePath + "/internal/infra/mysql"},
	}

	got := findBusinessRoots(packages)
	if _, ok := got["order"]; !ok {
		t.Fatal("findBusinessRoots() did not detect order module")
	}
	if _, ok := got["server"]; ok {
		t.Fatal("findBusinessRoots() detected server as a business module")
	}
	if _, ok := got["mysqlinfra"]; ok {
		t.Fatal("findBusinessRoots() detected infrastructure package as a business module")
	}
	if _, ok := got["infrastructure"]; ok {
		t.Fatal("findBusinessRoots() detected infrastructure as a business module")
	}
	if _, ok := got["infra"]; ok {
		t.Fatal("findBusinessRoots() detected infra as a business module")
	}
}

func TestStandardLibraryImportDetection(t *testing.T) {
	tests := []struct {
		importPath string
		want       bool
	}{
		{importPath: "context", want: true},
		{importPath: "net/http", want: true},
		{importPath: "github.com/labstack/echo/v4", want: false},
		{importPath: modulePath + "/internal/apperr", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.importPath, func(t *testing.T) {
			if got := isStandardLibraryImport(tt.importPath); got != tt.want {
				t.Fatalf("isStandardLibraryImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLayerPackageDetectionIgnoresInfrastructureRoots(t *testing.T) {
	if !isLayerPackage(modulePath+"/internal/order/mysql", "mysql") {
		t.Fatal("isLayerPackage() did not detect business mysql package")
	}
	if isLayerPackage(modulePath+"/internal/infrastructure/mysql", "mysql") {
		t.Fatal("isLayerPackage() detected infrastructure mysql package as business layer")
	}
	if isLayerPackage(modulePath+"/internal/infra/mysql", "mysql") {
		t.Fatal("isLayerPackage() detected infra mysql package as business layer")
	}
}

func listInternalPackages(t *testing.T) []goPackage {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))

	cmd := exec.Command("go", "list", "-json", "./internal/...")
	cmd.Dir = root

	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			t.Fatalf("go list ./internal/... failed: %v\n%s", err, exitErr.Stderr)
		}
		t.Fatalf("go list ./internal/... failed: %v", err)
	}

	dec := json.NewDecoder(strings.NewReader(string(out)))
	var packages []goPackage
	for {
		var pkg goPackage
		if err := dec.Decode(&pkg); err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("decode go list output: %v", err)
		}
		packages = append(packages, pkg)
	}
	return packages
}

func isLayerPackage(importPath string, layer string) bool {
	parts := internalParts(importPath)
	return len(parts) >= 2 && isBusinessRoot(parts[0]) && parts[1] == layer
}

func isInternalRootPackage(importPath string, root string) bool {
	parts := internalParts(importPath)
	return len(parts) >= 1 && parts[0] == root
}

func internalParts(importPath string) []string {
	const marker = "/internal/"
	index := strings.Index(importPath, marker)
	if index < 0 {
		return nil
	}
	return strings.Split(importPath[index+len(marker):], "/")
}

func findBusinessRoots(packages []goPackage) map[string]struct{} {
	roots := make(map[string]struct{})
	for _, pkg := range packages {
		parts := internalParts(pkg.ImportPath)
		if len(parts) < 2 {
			continue
		}
		if isBusinessRoot(parts[0]) && isBusinessLayer(parts[1]) {
			roots[parts[0]] = struct{}{}
		}
	}
	return roots
}

func isBusinessRoot(root string) bool {
	switch root {
	case "apperr", "archtest", "boot", "configs", "infra", "infrastructure", "requestctx", "server":
		return false
	default:
		return true
	}
}

func isBusinessLayer(layer string) bool {
	switch layer {
	case "domain", "usecase", "http", "mysql":
		return true
	default:
		return false
	}
}

func outerDetailForbiddenImports() []string {
	return []string{
		"github.com/labstack/echo",
		"gorm.io",
		"github.com/redis/go-redis",
		modulePath + "/internal/configs",
		modulePath + "/internal/log",
		modulePath + "/internal/server",
		modulePath + "/internal/server/httpresp",
		modulePath + "/internal/code",
	}
}

func usecaseForbiddenImports(importPath string) []string {
	forbidden := outerDetailForbiddenImports()
	parts := internalParts(importPath)
	if len(parts) == 0 {
		return forbidden
	}

	businessPath := modulePath + "/internal/" + parts[0]
	forbidden = append(forbidden,
		businessPath+"/http",
		businessPath+"/mysql",
	)
	return forbidden
}

func httpForbiddenImports(importPath string) []string {
	parts := internalParts(importPath)
	forbidden := make([]string, 0, 7)
	forbidden = append(forbidden,
		"gorm.io",
		"github.com/redis/go-redis",
		modulePath+"/internal/configs",
		modulePath+"/internal/log",
		modulePath+"/internal/code",
	)
	if len(parts) == 0 {
		return forbidden
	}

	businessPath := modulePath + "/internal/" + parts[0]
	forbidden = append(forbidden,
		businessPath+"/domain",
		businessPath+"/mysql",
	)
	return forbidden
}

func mysqlForbiddenImports(importPath string) []string {
	parts := internalParts(importPath)
	forbidden := make([]string, 0, 7)
	forbidden = append(forbidden,
		"github.com/labstack/echo",
		modulePath+"/internal/configs",
		modulePath+"/internal/log",
		modulePath+"/internal/server",
		modulePath+"/internal/server/httpresp",
		modulePath+"/internal/code",
	)
	if len(parts) == 0 {
		return forbidden
	}

	businessPath := modulePath + "/internal/" + parts[0]
	forbidden = append(forbidden,
		businessPath+"/http",
	)
	return forbidden
}

func assertNoForbiddenImports(t *testing.T, pkg goPackage, forbidden []string) {
	t.Helper()

	for _, imported := range pkg.Imports {
		for _, prefix := range forbidden {
			if matchesImport(imported, prefix) {
				t.Fatalf("%s imports forbidden package %s", pkg.ImportPath, imported)
			}
		}
	}
}

func assertDomainImportsStdlibOnly(t *testing.T, pkg goPackage) {
	t.Helper()

	for _, imported := range pkg.Imports {
		if !isStandardLibraryImport(imported) {
			t.Fatalf("%s imports non-standard-library package %s", pkg.ImportPath, imported)
		}
	}
}

func assertNoOtherBusinessImports(t *testing.T, pkg goPackage, businessRoots map[string]struct{}) {
	t.Helper()

	currentRoot, ok := businessRoot(pkg.ImportPath)
	if !ok {
		return
	}
	for root := range businessRoots {
		if root == currentRoot {
			continue
		}
		prefix := modulePath + "/internal/" + root
		assertNoForbiddenImports(t, pkg, []string{prefix})
	}
}

func assertNoBusinessImports(t *testing.T, pkg goPackage, businessRoots map[string]struct{}) {
	t.Helper()

	for root := range businessRoots {
		prefix := modulePath + "/internal/" + root
		assertNoForbiddenImports(t, pkg, []string{prefix})
	}
}

func businessRoot(importPath string) (string, bool) {
	parts := internalParts(importPath)
	if len(parts) < 2 || !isBusinessRoot(parts[0]) || !isBusinessLayer(parts[1]) {
		return "", false
	}
	return parts[0], true
}

func isStandardLibraryImport(imported string) bool {
	first, _, _ := strings.Cut(imported, "/")
	return !strings.Contains(first, ".")
}

func matchesImport(imported string, forbidden string) bool {
	return imported == forbidden || strings.HasPrefix(imported, forbidden+"/")
}
