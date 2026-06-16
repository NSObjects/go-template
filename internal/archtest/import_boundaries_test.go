package archtest

import (
	"encoding/json"
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

	for _, pkg := range packages {
		switch {
		case isLayerPackage(pkg.ImportPath, "domain"):
			assertNoForbiddenImports(t, pkg, domainForbiddenImports())
		case isLayerPackage(pkg.ImportPath, "usecase"):
			assertNoForbiddenImports(t, pkg, usecaseForbiddenImports(pkg.ImportPath))
		}
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
		if exitErr, ok := err.(*exec.ExitError); ok {
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
	return len(parts) >= 2 && parts[1] == layer
}

func internalParts(importPath string) []string {
	const marker = "/internal/"
	index := strings.Index(importPath, marker)
	if index < 0 {
		return nil
	}
	return strings.Split(importPath[index+len(marker):], "/")
}

func domainForbiddenImports() []string {
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
	forbidden := domainForbiddenImports()
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

func matchesImport(imported string, forbidden string) bool {
	return imported == forbidden || strings.HasPrefix(imported, forbidden+"/")
}
