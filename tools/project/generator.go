package project

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// GeneratorConfig contains the configuration required to scaffold a project copy.
type GeneratorConfig struct {
	SourceRoot          string
	OutputDir           string
	TemplateModulePath  string
	TargetModulePath    string
	TemplateProjectName string
	TargetProjectName   string
	Force               bool
}

// Generator copies the template repository into a new project directory.
type Generator struct {
	cfg             GeneratorConfig
	replacer        *strings.Replacer
	excludePrefixes []string
}

// NewGenerator validates configuration and prepares a generator instance.
func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	if cfg.SourceRoot == "" {
		return nil, errors.New("未提供模板仓库路径")
	}
	if cfg.OutputDir == "" {
		return nil, errors.New("未提供生成目标目录")
	}
	if cfg.TemplateModulePath == "" {
		return nil, errors.New("未提供模板 Module 路径")
	}
	if cfg.TargetModulePath == "" {
		return nil, errors.New("未提供目标 Module 路径")
	}

	if cfg.TemplateProjectName == "" {
		cfg.TemplateProjectName = filepath.Base(cfg.TemplateModulePath)
	}
	if cfg.TargetProjectName == "" {
		cfg.TargetProjectName = filepath.Base(cfg.TargetModulePath)
	}

	pairs := []string{
		cfg.TemplateModulePath, cfg.TargetModulePath,
	}
	if cfg.TemplateProjectName != cfg.TargetProjectName {
		pairs = append(pairs,
			cfg.TemplateProjectName, cfg.TargetProjectName,
			strings.ToUpper(cfg.TemplateProjectName), strings.ToUpper(cfg.TargetProjectName),
		)
	}

	replacer := strings.NewReplacer(pairs...)

	return &Generator{
		cfg:      cfg,
		replacer: replacer,
		excludePrefixes: []string{
			"tools",
			".git",
			".github",
			".idea",
			".vscode",
		},
	}, nil
}

// Generate performs the copy and token replacement process.
func (g *Generator) Generate() error {
	info, err := os.Stat(g.cfg.SourceRoot)
	if err != nil {
		return fmt.Errorf("模板仓库不存在: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("模板路径不是目录: %s", g.cfg.SourceRoot)
	}

	outputDir := g.cfg.OutputDir
	if err := g.prepareOutputDir(outputDir); err != nil {
		return err
	}

	return filepath.WalkDir(g.cfg.SourceRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(g.cfg.SourceRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		if g.shouldSkip(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		targetPath := filepath.Join(outputDir, rel)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0o755)
		}

		return g.copyFile(path, targetPath, d)
	})
}

func (g *Generator) prepareOutputDir(outputDir string) error {
	if stat, err := os.Stat(outputDir); err == nil {
		if !stat.IsDir() {
			return fmt.Errorf("目标路径已存在且不是目录: %s", outputDir)
		}
		if !g.cfg.Force {
			return fmt.Errorf("目标目录已存在: %s (使用 --force 可覆盖)", outputDir)
		}
		if err := os.RemoveAll(outputDir); err != nil {
			return fmt.Errorf("清理目标目录失败: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("检查目标目录失败: %w", err)
	}

	return os.MkdirAll(outputDir, 0o755)
}

func (g *Generator) shouldSkip(rel string) bool {
	rel = filepath.ToSlash(rel)

	// 检查是否匹配任何排除前缀（支持多级路径）
	for _, prefix := range g.excludePrefixes {
		if strings.HasPrefix(rel, prefix+"/") || rel == prefix {
			return true
		}
	}

	// 检查是否是输出目录本身（避免递归复制）
	outputBase := filepath.Base(g.cfg.OutputDir)
	if rel == outputBase || strings.HasPrefix(rel, outputBase+"/") {
		return true
	}

	if strings.HasSuffix(rel, ".DS_Store") {
		return true
	}

	return false
}

func (g *Generator) copyFile(src, dest string, entry fs.DirEntry) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("读取文件失败 %s: %w", src, err)
	}

	out := data
	if utf8.Valid(data) {
		replaced := g.replacer.Replace(string(data))
		out = []byte(replaced)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("创建目录失败 %s: %w", filepath.Dir(dest), err)
	}

	mode := fs.FileMode(0o644)
	if info, err := entry.Info(); err == nil {
		mode = info.Mode()
	}

	if err := os.WriteFile(dest, out, mode.Perm()); err != nil {
		return fmt.Errorf("写入文件失败 %s: %w", dest, err)
	}

	return nil
}
