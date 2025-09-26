package project

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

//go:embed templates/*
var templateFS embed.FS

// EmbedTemplateGenerator 基于嵌入文件的模板生成器
type EmbedTemplateGenerator struct {
	outputDir string
	data      templateData
	fs        fs.FS
	logger    func(string, ...interface{})
	modes     ModeResolver
}

type templateData struct {
	ModulePath  string
	ProjectName string
	DisplayName string
	PackageName string
	EnvName     string

}

type fileSpec struct {
	outputPath     string
	templatePath   string
	renderTemplate bool
}

// ModeResolver decides the output permission bits for generated files.
type ModeResolver func(outputPath string) fs.FileMode

// Option customizes the behaviour of the EmbedTemplateGenerator.
type Option func(*EmbedTemplateGenerator)

// NewEmbedTemplateGenerator 创建新的嵌入模板生成器
func NewEmbedTemplateGenerator(outputDir, modulePath, projectName string, opts ...Option) *EmbedTemplateGenerator {
	generator := &EmbedTemplateGenerator{
		outputDir: outputDir,
		data:      buildTemplateData(modulePath, projectName),
		fs:        templateFS,
		logger: func(format string, args ...interface{}) {
			fmt.Printf(format, args...)
		},
		modes: desiredFileMode,
	}

	for _, opt := range opts {
		opt(generator)
	}

	return generator
}

// WithFS allows injecting a custom file system for template discovery (useful in tests).
func WithFS(fsys fs.FS) Option {
	return func(g *EmbedTemplateGenerator) {
		if fsys != nil {
			g.fs = fsys
		}
	}
}

// WithLogger customises where progress output is written.
func WithLogger(logger func(string, ...interface{})) Option {
	return func(g *EmbedTemplateGenerator) {
		if logger != nil {
			g.logger = logger
		}
	}
}

// WithModeResolver overrides how output permissions are decided.
func WithModeResolver(resolver ModeResolver) Option {
	return func(g *EmbedTemplateGenerator) {
		if resolver != nil {
			g.modes = resolver
		}
	}
}

// Generate 生成项目文件
func (g *EmbedTemplateGenerator) Generate() error {
	filesToGenerate, err := g.filesToGenerate()
	if err != nil {
		return fmt.Errorf("读取模板文件列表失败: %w", err)
	}

	if err := g.createDirectories(filesToGenerate); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	if err := g.generateFiles(filesToGenerate); err != nil {
		return fmt.Errorf("生成文件失败: %w", err)
	}

	return nil
}

func (g *EmbedTemplateGenerator) filesToGenerate() ([]fileSpec, error) {
	files := make([]fileSpec, 0)
	outputs := make(map[string]string)

	err := fs.WalkDir(g.fs, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath := strings.TrimPrefix(path, "templates/")
		if relPath == path {
			return fmt.Errorf("无法解析模板路径: %s", path)
		}

		relPath = filepath.FromSlash(relPath)

		if filepath.Base(relPath) == ".gitkeep" {
			return nil
		}

		spec := fileSpec{
			outputPath:     relPath,
			templatePath:   path,
			renderTemplate: strings.HasSuffix(relPath, ".tmpl"),
		}

		if spec.renderTemplate {
			spec.outputPath = strings.TrimSuffix(spec.outputPath, ".tmpl")
		}


		if g.shouldSkip(spec.outputPath) {
			return nil
		}


		if prev, exists := outputs[spec.outputPath]; exists {
			return fmt.Errorf("检测到重复模板输出路径: %s (由 %s 和 %s 提供)", spec.outputPath, prev, path)
		}

		outputs[spec.outputPath] = path
		files = append(files, spec)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].outputPath < files[j].outputPath
	})

	return files, nil
}

// createDirectories 根据待生成的文件自动创建目录结构
func (g *EmbedTemplateGenerator) createDirectories(files []fileSpec) error {
	dirSet := make(map[string]struct{})
	for _, file := range files {
		dir := filepath.Dir(file.outputPath)
		if dir == "." {
			continue
		}
		dirSet[dir] = struct{}{}
	}

	dirs := make([]string, 0, len(dirSet))
	for dir := range dirSet {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	for _, dir := range dirs {
		fullPath := filepath.Join(g.outputDir, dir)
		if err := os.MkdirAll(fullPath, 0o755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
		g.printf("✅ 创建目录: %s\n", fullPath)
	}

	return nil
}

// generateFiles 生成所有文件
func (g *EmbedTemplateGenerator) generateFiles(files []fileSpec) error {
	for _, file := range files {
		if err := g.generateFile(file); err != nil {
			return fmt.Errorf("生成文件 %s 失败: %w", file.outputPath, err)
		}
		g.printf("✅ 生成文件: %s\n", filepath.Join(g.outputDir, file.outputPath))
	}

	return nil
}

// generateFile 生成单个文件
func (g *EmbedTemplateGenerator) generateFile(file fileSpec) error {
	// 从嵌入的文件系统读取模板文件
	templateContent, err := fs.ReadFile(g.fs, file.templatePath)
	if err != nil {
		return fmt.Errorf("读取模板文件失败: %w", err)
	}

	// 创建输出文件
	outputPath := filepath.Join(g.outputDir, file.outputPath)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	mode := g.modes(file.outputPath)

	if file.renderTemplate {
		tmpl, err := template.New(filepath.Base(file.templatePath)).Parse(string(templateContent))
		if err != nil {
			return fmt.Errorf("解析模板文件失败: %w", err)
		}

		outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
		if err != nil {
			return fmt.Errorf("创建输出文件失败: %w", err)
		}
		if err := tmpl.Execute(outputFile, g.data); err != nil {
			closeErr := outputFile.Close()
			if closeErr != nil {
				return fmt.Errorf("执行模板失败: %v (关闭输出文件失败: %w)", err, closeErr)
			}
			return fmt.Errorf("执行模板失败: %w", err)
		}

		if err := outputFile.Close(); err != nil {
			return fmt.Errorf("关闭输出文件失败: %w", err)
		}

	} else {
		if err := os.WriteFile(outputPath, templateContent, mode); err != nil {
			return fmt.Errorf("写入文件失败: %w", err)
		}
	}

	if err := os.Chmod(outputPath, mode); err != nil {
		return fmt.Errorf("设置文件权限失败: %w", err)
	}

	return nil
}

func (g *EmbedTemplateGenerator) printf(format string, args ...interface{}) {
	if g.logger == nil {
		return
	}
	g.logger(format, args...)
}

func (g *EmbedTemplateGenerator) shouldSkip(outputPath string) bool {
	normalized := filepath.ToSlash(outputPath)

	if normalized == "go.sum" {
		return true
	}

	if strings.HasPrefix(normalized, "docs/") && strings.HasSuffix(normalized, ".md") {
		return true
	}

	return false
}

func desiredFileMode(outputPath string) fs.FileMode {
	switch filepath.Ext(outputPath) {
	case ".sh":
		return 0o755
	default:
		return 0o644
	}
}

func buildTemplateData(modulePath, projectName string) templateData {
	sanitizedModule := strings.TrimSpace(modulePath)
	base := filepath.Base(sanitizedModule)
	pkgName := sanitizePackageName(base)

	trimmedProject := strings.TrimSpace(projectName)
	if trimmedProject == "" {
		trimmedProject = base
	}

	display := trimmedProject
	slug := toKebabCase(trimmedProject)
	if slug == "" {
		slug = pkgName
	}
	if slug == "" {
		slug = "app"
	}

	envName := toSnakeCase(slug)
	if envName == "" {
		envName = slug
	}

	if display == "" {
		display = slug
	}

	return templateData{
		ModulePath:  sanitizedModule,
		ProjectName: slug,
		DisplayName: display,
		PackageName: pkgName,
		EnvName:     envName,
	}
}

func sanitizePackageName(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	if lower == "" {
		return "app"
	}

	runes := make([]rune, 0, len(lower))
	prevUnderscore := false

	for _, r := range lower {
		switch {
		case unicode.IsLetter(r):
			runes = append(runes, r)
			prevUnderscore = false
		case unicode.IsDigit(r):
			if len(runes) == 0 {
				runes = append(runes, '_')
			}
			runes = append(runes, r)
			prevUnderscore = false
		case r == '_' || r == '.':
			if len(runes) > 0 && !prevUnderscore {
				runes = append(runes, '_')
				prevUnderscore = true
			}
		default:
			if len(runes) > 0 && !prevUnderscore {
				runes = append(runes, '_')
				prevUnderscore = true
			}
		}
	}

	if len(runes) == 0 {
		return "app"
	}
	if runes[len(runes)-1] == '_' {
		runes = runes[:len(runes)-1]
	}
	if len(runes) == 0 {
		return "app"
	}
	return string(runes)
}

func toKebabCase(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	if lower == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(lower))
	prevHyphen := false
	for _, r := range lower {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			prevHyphen = false
		case r == '-' || r == '_' || r == ' ':
			if !prevHyphen {
				b.WriteByte('-')
				prevHyphen = true
			}
		default:
			if !prevHyphen {
				b.WriteByte('-')
				prevHyphen = true
			}
		}
	}

	slug := strings.Trim(b.String(), "-")
	return slug
}

func toSnakeCase(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	if lower == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(lower))
	prevUnderscore := false
	for _, r := range lower {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			prevUnderscore = false
		case r == '-' || r == '_' || r == ' ':
			if !prevUnderscore {
				b.WriteByte('_')
				prevUnderscore = true
			}
		default:
			if !prevUnderscore {
				b.WriteByte('_')
				prevUnderscore = true
			}
		}

	}
	g.logger(format, args...)
}

	sanitized := strings.Trim(b.String(), "_")
	return sanitized

}
