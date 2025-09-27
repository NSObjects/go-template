package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	modutils "github.com/NSObjects/go-template/muban/modgen/utils"
)

func (g *Generator) ensureSupportFiles() error {
	if err := g.ensureContextSupport(); err != nil {
		return err
	}
	if err := g.ensureRespSupport(); err != nil {
		return err
	}
	return nil
}

func (g *Generator) ensureContextSupport() error {
	utilsDir := filepath.Join(g.config.RepoRoot, "internal", "utils")

	has, err := packageHasSymbol(utilsDir, "BuildContext(")
	if err != nil {
		return fmt.Errorf("检查 utils 上下文支持函数失败: %w", err)
	}

	supportTarget := filepath.Join(g.config.RepoRoot, "internal", "utils", "context_support.go")
	legacyTarget := filepath.Join(g.config.RepoRoot, "internal", "utils", "context_trace.go")

	if has {
		return ensureLegacyContextTraceStub(legacyTarget)
	}

	renderer, err := g.templateRenderer()
	if err != nil {
		return fmt.Errorf("创建模板渲染器失败: %w", err)
	}

	content, err := renderer.RenderContextSupport()
	if err != nil {
		return fmt.Errorf("渲染 context 支持模板失败: %w", err)
	}

	modutils.MustWrite(supportTarget, content, true)
	return ensureLegacyContextTraceStub(legacyTarget)

}

const legacyContextTraceStub = `package utils


// Deprecated: context_trace.go 已被迁移至 context_support.go。
// 本文件仅用于兼容老版本引用，同时避免重复定义。
`
func ensureLegacyContextTraceStub(target string) error {

	data, err := os.ReadFile(target)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("创建兼容性文件目录失败: %w", err)
			}
			if err := os.WriteFile(target, []byte(legacyContextTraceStub), 0o644); err != nil {
				return fmt.Errorf("写入兼容性 context_trace.go 失败: %w", err)
			}
			fmt.Printf("写入兼容性文件: %s\n", target)
			return nil
		}
		return fmt.Errorf("读取遗留 context_trace.go 失败: %w", err)
	}

	content := string(data)

	trimmed := strings.TrimSpace(content)
	if trimmed == strings.TrimSpace(legacyContextTraceStub) {
		return nil
	}

	if trimmed == "" || trimmed == "package utils" {
		if err := os.WriteFile(target, []byte(legacyContextTraceStub), 0o644); err != nil {
			return fmt.Errorf("写入兼容性 context_trace.go 失败: %w", err)
		}
		fmt.Printf("重写兼容性文件: %s\n", target)
		return nil
	}

	if strings.Contains(content, "func BuildContext") || strings.Contains(content, "ExtractTraceContext") {
		if err := os.WriteFile(target, []byte(legacyContextTraceStub), 0o644); err != nil {
			return fmt.Errorf("写入兼容性 context_trace.go 失败: %w", err)
		}

		fmt.Printf("重写兼容性文件: %s\n", target)
		return nil
	}

	return nil
}

func (g *Generator) ensureRespSupport() error {
	respDir := filepath.Join(g.config.RepoRoot, "internal", "resp")

	hasOne, err := packageHasSymbol(respDir, "OneDataResponse(")
	if err != nil {
		return fmt.Errorf("检查 resp 数据响应函数失败: %w", err)
	}
	hasOperate, err := packageHasSymbol(respDir, "OperateSuccess(")
	if err != nil {
		return fmt.Errorf("检查 resp 操作响应函数失败: %w", err)
	}
	if hasOne && hasOperate {
		return nil
	}

	target := filepath.Join(respDir, "response_helpers.go")
	if _, err := os.Stat(target); err == nil {
		return nil
	}

	renderer, err := g.templateRenderer()
	if err != nil {
		return fmt.Errorf("创建模板渲染器失败: %w", err)
	}

	content, err := renderer.RenderRespSupport()
	if err != nil {
		return fmt.Errorf("渲染 resp 支持模板失败: %w", err)
	}

	modutils.MustWrite(target, content, false)
	return nil
}

func packageHasSymbol(dir, symbol string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return false, err
		}
		if strings.Contains(string(data), symbol) {
			return true, nil
		}
	}

	return false, nil
}
