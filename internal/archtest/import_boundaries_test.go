package archtest

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const modulePath = "github.com/NSObjects/go-template"

const (
	internalModulesDir  = "modules"
	internalPlatformDir = "platform"
	adaptersDir         = "adapters"
)

type goPackage struct {
	ImportPath string
	Imports    []string
}

func TestInternalRootOnlyExposesFrameworkModuleBuckets(t *testing.T) {
	root := repositoryRoot(t)
	entries, err := os.ReadDir(filepath.Join(root, "internal"))
	if err != nil {
		t.Fatalf("read internal root: %v", err)
	}

	allowed := map[string]struct{}{
		"archtest":          {},
		"boot":              {},
		internalModulesDir:  {},
		internalPlatformDir: {},
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			t.Fatalf("internal root contains file %q, want only archtest, boot, modules, platform directories", entry.Name())
		}
		if _, ok := allowed[entry.Name()]; !ok {
			t.Fatalf("internal root exposes %q, want only archtest, boot, modules, platform", entry.Name())
		}
	}
}

func TestBusinessImportBoundaries(t *testing.T) {
	packages := listInternalPackages(t)
	businessRoots := findBusinessRoots(packages)
	adapterPathsByRoot := findBusinessAdapterPathsByRoot(packages)

	for _, pkg := range packages {
		switch {
		case isLayerPackage(pkg.ImportPath, "domain"):
			assertDomainImportsStdlibOnly(t, pkg)
			assertNoOtherBusinessImports(t, pkg, businessRoots)
		case isLayerPackage(pkg.ImportPath, "usecase"):
			assertNoForbiddenImports(t, pkg, usecaseForbiddenImports(pkg.ImportPath, adapterPathsByRoot))
			assertNoOtherBusinessImports(t, pkg, businessRoots)
		case isLayerPackage(pkg.ImportPath, "http"):
			assertNoForbiddenImports(t, pkg, httpForbiddenImports(pkg.ImportPath, adapterPathsByRoot))
			assertNoOtherBusinessImports(t, pkg, businessRoots)
		case isBusinessAdapterPackage(pkg.ImportPath):
			assertNoForbiddenImports(t, pkg, adapterForbiddenImports(pkg.ImportPath))
			assertNoOtherBusinessImports(t, pkg, businessRoots)
		case isPlatformPackage(pkg.ImportPath):
			assertNoBusinessImports(t, pkg, businessRoots)
		}
	}
}

func TestFindBusinessRoots(t *testing.T) {
	packages := []goPackage{
		{ImportPath: modulePath + "/internal/modules/order/usecase"},
		{ImportPath: modulePath + "/internal/modules/order/http"},
		{ImportPath: modulePath + "/internal/modules/order/adapters/memory"},
		{ImportPath: modulePath + "/internal/modules/order/adapters/redis"},
		{ImportPath: modulePath + "/internal/platform/server/httpresp"},
		{ImportPath: modulePath + "/internal/mysqlinfra"},
		{ImportPath: modulePath + "/internal/platform/infrastructure/mysql"},
		{ImportPath: modulePath + "/internal/tools/mysql"},
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
	if _, ok := got["tools"]; ok {
		t.Fatal("findBusinessRoots() detected tools as a business module")
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
		{importPath: modulePath + "/internal/platform/apperr", want: false},
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
	if !isLayerPackage(modulePath+"/internal/modules/order/domain", "domain") {
		t.Fatal("isLayerPackage() did not detect business domain package")
	}
	if !isBusinessAdapterPackage(modulePath + "/internal/modules/order/adapters/redis") {
		t.Fatal("isBusinessAdapterPackage() did not detect business redis adapter")
	}
	if isLayerPackage(modulePath+"/internal/platform/infrastructure/mysql", "mysql") {
		t.Fatal("isLayerPackage() detected infrastructure mysql package as business layer")
	}
	if isLayerPackage(modulePath+"/internal/tools/mysql", "mysql") {
		t.Fatal("isLayerPackage() detected tools mysql package as business layer")
	}
}

func TestPlatformPackageDetectionIncludesPlatformRoot(t *testing.T) {
	tests := []struct {
		importPath string
		want       bool
	}{
		{importPath: modulePath + "/internal/platform", want: true},
		{importPath: modulePath + "/internal/platform/server", want: true},
		{importPath: modulePath + "/internal/modules/order", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.importPath, func(t *testing.T) {
			if got := isPlatformPackage(tt.importPath); got != tt.want {
				t.Fatalf("isPlatformPackage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForbiddenImportsCoverConcreteAdapters(t *testing.T) {
	packages := []goPackage{
		{ImportPath: modulePath + "/internal/modules/order/usecase"},
		{ImportPath: modulePath + "/internal/modules/order/http"},
		{ImportPath: modulePath + "/internal/modules/order/adapters/memory"},
		{ImportPath: modulePath + "/internal/modules/order/adapters/client"},
		{ImportPath: modulePath + "/internal/modules/order/adapters/queue"},
		{ImportPath: modulePath + "/internal/modules/order/adapters/redis"},
	}
	adapterPathsByRoot := findBusinessAdapterPathsByRoot(packages)

	usecaseForbidden := usecaseForbiddenImports(modulePath+"/internal/modules/order/usecase", adapterPathsByRoot)
	if !containsString(usecaseForbidden, modulePath+"/internal/modules/order/adapters/memory") {
		t.Fatal("usecase forbidden imports do not include memory adapter")
	}
	if !containsString(usecaseForbidden, modulePath+"/internal/modules/order/adapters/client") {
		t.Fatal("usecase forbidden imports do not include client adapter")
	}
	if !containsString(usecaseForbidden, modulePath+"/internal/modules/order/adapters/redis") {
		t.Fatal("usecase forbidden imports do not include redis adapter")
	}
	if !containsString(usecaseForbidden, modulePath+"/internal/platform/infrastructure") {
		t.Fatal("usecase forbidden imports do not include shared infrastructure")
	}

	httpForbidden := httpForbiddenImports(modulePath+"/internal/modules/order/http", adapterPathsByRoot)
	if !containsString(httpForbidden, modulePath+"/internal/modules/order/adapters/queue") {
		t.Fatal("http forbidden imports do not include queue adapter")
	}
	if !containsString(httpForbidden, modulePath+"/internal/platform/infrastructure") {
		t.Fatal("http forbidden imports do not include shared infrastructure")
	}

	adapterForbidden := adapterForbiddenImports(modulePath + "/internal/modules/order/adapters/client")
	if !containsString(adapterForbidden, modulePath+"/internal/platform/server/httpresp") {
		t.Fatal("adapter forbidden imports do not include server http response package")
	}
	if !containsString(adapterForbidden, modulePath+"/internal/modules/order/http") {
		t.Fatal("adapter forbidden imports do not include same-business HTTP adapter")
	}
}

func TestForbiddenImportsCoverMongoDBAndTracingDrivers(t *testing.T) {
	adapterPathsByRoot := map[string][]string{"order": nil}

	usecaseForbidden := usecaseForbiddenImports(modulePath+"/internal/modules/order/usecase", adapterPathsByRoot)
	if !containsString(usecaseForbidden, "go.mongodb.org/mongo-driver") {
		t.Fatal("usecase forbidden imports do not include MongoDB driver")
	}
	if !containsString(usecaseForbidden, "go.opentelemetry.io/otel") {
		t.Fatal("usecase forbidden imports do not include OpenTelemetry")
	}

	httpForbidden := httpForbiddenImports(modulePath+"/internal/modules/order/http", adapterPathsByRoot)
	if !containsString(httpForbidden, "go.mongodb.org/mongo-driver") {
		t.Fatal("http forbidden imports do not include MongoDB driver")
	}
	if !containsString(httpForbidden, "go.opentelemetry.io/otel") {
		t.Fatal("http forbidden imports do not include OpenTelemetry")
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

func isLayerPackage(importPath, layer string) bool {
	_, parts, ok := businessModuleParts(importPath)
	return ok && len(parts) >= 1 && parts[0] == layer
}

func isBusinessAdapterPackage(importPath string) bool {
	_, parts, ok := businessModuleParts(importPath)
	return ok && len(parts) >= 2 && parts[0] == adaptersDir && parts[1] != ""
}

func isPlatformPackage(importPath string) bool {
	parts := internalParts(importPath)
	return len(parts) >= 1 && parts[0] == internalPlatformDir
}

func internalParts(importPath string) []string {
	const marker = "/internal/"
	index := strings.Index(importPath, marker)
	if index < 0 {
		return nil
	}
	return strings.Split(importPath[index+len(marker):], "/")
}

func businessModuleParts(importPath string) (string, []string, bool) {
	parts := internalParts(importPath)
	if len(parts) < 3 || parts[0] != internalModulesDir || parts[1] == "" {
		return "", nil, false
	}
	return parts[1], parts[2:], true
}

func businessPath(root string) string {
	return modulePath + "/internal/modules/" + root
}

func findBusinessRoots(packages []goPackage) map[string]struct{} {
	roots := make(map[string]struct{})
	for _, pkg := range packages {
		root, parts, ok := businessModuleParts(pkg.ImportPath)
		if !ok || len(parts) < 1 {
			continue
		}
		if isBusinessLayer(parts[0]) {
			roots[root] = struct{}{}
		}
	}
	return roots
}

func findBusinessAdapterPathsByRoot(packages []goPackage) map[string][]string {
	pathsByRoot := make(map[string][]string)
	seen := make(map[string]struct{})
	for _, pkg := range packages {
		root, parts, ok := businessModuleParts(pkg.ImportPath)
		if !ok || len(parts) < 2 || parts[0] != adaptersDir || parts[1] == "" {
			continue
		}

		adapterPath := businessPath(root) + "/" + adaptersDir + "/" + parts[1]
		if _, ok := seen[adapterPath]; ok {
			continue
		}
		seen[adapterPath] = struct{}{}
		pathsByRoot[root] = append(pathsByRoot[root], adapterPath)
	}
	return pathsByRoot
}

func isBusinessLayer(layer string) bool {
	switch layer {
	case "domain", "usecase", "http":
		return true
	default:
		return isBusinessAdapterLayer(layer)
	}
}

func isBusinessAdapterLayer(layer string) bool {
	switch layer {
	case adaptersDir:
		return true
	default:
		return false
	}
}

func sharedInfrastructurePaths() []string {
	return []string{modulePath + "/internal/platform/infrastructure"}
}

func outerDetailForbiddenImports() []string {
	return []string{
		"github.com/labstack/echo",
		"gorm.io",
		"github.com/redis/go-redis",
		"go.mongodb.org/mongo-driver",
		"go.opentelemetry.io/otel",
		modulePath + "/internal/platform/configs",
		modulePath + "/internal/log",
		modulePath + "/internal/platform/server",
		modulePath + "/internal/platform/server/httpresp",
		modulePath + "/internal/code",
	}
}

func usecaseForbiddenImports(importPath string, adapterPathsByRoot map[string][]string) []string {
	forbidden := outerDetailForbiddenImports()
	forbidden = append(forbidden, sharedInfrastructurePaths()...)
	root, _, ok := businessModuleParts(importPath)
	if !ok {
		return forbidden
	}

	forbidden = append(forbidden,
		businessPath(root)+"/http",
	)
	forbidden = append(forbidden, adapterPathsByRoot[root]...)
	return forbidden
}

func httpForbiddenImports(importPath string, adapterPathsByRoot map[string][]string) []string {
	root, _, ok := businessModuleParts(importPath)
	forbidden := make([]string, 0, 7)
	forbidden = append(forbidden,
		"gorm.io",
		"github.com/redis/go-redis",
		"go.mongodb.org/mongo-driver",
		"go.opentelemetry.io/otel",
		modulePath+"/internal/platform/configs",
		modulePath+"/internal/log",
		modulePath+"/internal/code",
	)
	forbidden = append(forbidden, sharedInfrastructurePaths()...)
	if !ok {
		return forbidden
	}

	forbidden = append(forbidden,
		businessPath(root)+"/domain",
	)
	forbidden = append(forbidden, adapterPathsByRoot[root]...)
	return forbidden
}

func adapterForbiddenImports(importPath string) []string {
	root, _, ok := businessModuleParts(importPath)
	forbidden := make([]string, 0, 7)
	forbidden = append(forbidden,
		"github.com/labstack/echo",
		modulePath+"/internal/platform/configs",
		modulePath+"/internal/log",
		modulePath+"/internal/platform/server",
		modulePath+"/internal/platform/server/httpresp",
		modulePath+"/internal/code",
	)
	if !ok {
		return forbidden
	}

	forbidden = append(forbidden,
		businessPath(root)+"/http",
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
		assertNoForbiddenImports(t, pkg, []string{businessPath(root)})
	}
}

func assertNoBusinessImports(t *testing.T, pkg goPackage, businessRoots map[string]struct{}) {
	t.Helper()

	for root := range businessRoots {
		assertNoForbiddenImports(t, pkg, []string{businessPath(root)})
	}
}

func businessRoot(importPath string) (string, bool) {
	root, parts, ok := businessModuleParts(importPath)
	if !ok || len(parts) < 1 || !isBusinessLayer(parts[0]) {
		return "", false
	}
	return root, true
}

func isStandardLibraryImport(imported string) bool {
	first, _, _ := strings.Cut(imported, "/")
	return !strings.Contains(first, ".")
}

func matchesImport(imported, forbidden string) bool {
	return imported == forbidden || strings.HasPrefix(imported, forbidden+"/")
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
