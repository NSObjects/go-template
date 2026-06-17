package archtest

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDocumentationDescribesConfiguredCapabilities(t *testing.T) {
	root := repositoryRoot(t)
	files := []string{
		"internal/boot/README.md",
		"internal/platform/server/README.md",
		"internal/platform/server/middlewares/README.md",
		"configs/README.md",
	}

	combined := strings.Builder{}
	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		combined.Write(data)
		combined.WriteByte('\n')
	}
	text := strings.ToLower(combined.String())

	required := []string{
		"configured infrastructure resources",
		"mysql",
		"redis",
		"mongodb",
		"jaeger",
		"adapter",
		"usecase-owned outbound interface",
	}
	for _, want := range required {
		if !strings.Contains(text, want) {
			t.Fatalf("documentation missing %q", want)
		}
	}

	forbidden := []string{"manifest entry", "placeholder resource"}
	for _, word := range forbidden {
		if strings.Contains(text, word) {
			t.Fatalf("documentation contains forbidden placeholder workflow %q", word)
		}
	}
}

func TestRepositoryDoesNotExposeOldTemplateConceptNames(t *testing.T) {
	root := repositoryRoot(t)
	files := []string{
		"internal/boot/README.md",
		"internal/platform/server/README.md",
		"internal/platform/server/middlewares/README.md",
		"configs/README.md",
		"configs/config.toml",
		"env.example",
		"docker-compose.yaml",
		"k8s/deployment.yaml",
	}
	forbidden := []string{
		"clean-lite",
		"cleanlite",
		"generated demo user",
		"default password",
		"默认密码",
	}

	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		text := strings.ToLower(string(data))
		for _, word := range forbidden {
			if strings.Contains(text, word) {
				t.Fatalf("%s contains retired template wording %q", file, word)
			}
		}
	}
}

func repositoryRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate test file")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
}
