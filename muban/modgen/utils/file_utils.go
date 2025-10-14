/*
 * 文件处理工具函数
 */

package utils

import (
	"bufio"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"
)

// FindRepoRoot 查找仓库根目录
func FindRepoRoot(start string) string {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// GetPackagePath 从go.mod文件获取项目包路径
func GetPackagePath(repoRoot string) (string, error) {
	goModPath := filepath.Join(repoRoot, "go.mod")
	file, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		// 解析 module github.com/NSObjects/go-template
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[0] == "module" {
			return parts[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("无法从go.mod解析模块路径")
}

// MustWrite 写入文件
func MustWrite(path, content string, force bool) {
	if !force {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("跳过已存在文件: %s (使用 --force 可覆盖)\n", path)
			return
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		ExitWith(err.Error())
	}
	processed := []byte(content)
	if strings.HasSuffix(strings.ToLower(path), ".go") && os.Getenv("DISABLE_GEN_FORMAT") != "1" {
		if out, err := imports.Process(path, processed, &imports.Options{Comments: true, TabIndent: true, TabWidth: 8}); err == nil {
			processed = out
		}
		if out, err := format.Source(processed); err == nil {
			processed = out
		}
	}
	if err := os.WriteFile(path, processed, 0o644); err != nil {
		ExitWith(err.Error())
	}
	fmt.Printf("写入: %s\n", path)
}

// TryInject 尝试注入到fx.Options
func TryInject(filePath, item string) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	s := string(b)
	if strings.Contains(s, item) {
		return nil
	}

	// 查找 fx.Provide 的位置
	fxProvideIndex := strings.Index(s, "fx.Provide(")
	if fxProvideIndex == -1 {
		// 如果没有找到 fx.Provide，尝试查找 fx.Options
		fxOptionsIndex := strings.Index(s, "fx.Options(")
		if fxOptionsIndex == -1 {
			return nil
		}
		// 在 fx.Options 中添加 fx.Provide 调用
		start := fxOptionsIndex + len("fx.Options(")
		// 找到 fx.Options 的结束位置
		parenCount := 1
		end := start
		for i := start; i < len(s); i++ {
			if s[i] == '(' {
				parenCount++
			} else if s[i] == ')' {
				parenCount--
				if parenCount == 0 {
					end = i
					break
				}
			}
		}
		if end < start {
			return nil
		}
		before := s[:start]
		inside := s[start:end]
		after := s[end:]

		// 添加 fx.Provide 调用
		trimmed := strings.TrimSpace(inside)
		if trimmed != "" && !strings.HasSuffix(strings.TrimSpace(trimmed), ",") {
			inside = inside + ",\n"
		}
		inside = inside + "\t" + "fx.Provide(" + item + "),\n"
		out := before + inside + after
		return os.WriteFile(filePath, []byte(out), 0o644)
	}

	// 检查 fx.Provide 是否为空
	start := fxProvideIndex + len("fx.Provide(")
	// 找到 fx.Provide 的结束位置
	parenCount := 1
	end := start
	for i := start; i < len(s); i++ {
		if s[i] == '(' {
			parenCount++
		} else if s[i] == ')' {
			parenCount--
			if parenCount == 0 {
				end = i
				break
			}
		}
	}

	if end < start {
		fmt.Printf("Invalid fx.Provide structure: start=%d, end=%d\n", start, end)
		return nil
	}

	// 检查 fx.Provide 内部是否为空
	inside := s[start:end]
	trimmed := strings.TrimSpace(inside)
	if trimmed == "" || trimmed == "\n" {
		// fx.Provide 为空，直接添加项目
		before := s[:start]
		after := s[end:]
		inside = item
		out := before + inside + after
		return os.WriteFile(filePath, []byte(out), 0o644)
	}

	// fx.Provide 不为空，添加新的项目
	before := s[:start]
	after := s[end:]

	// 尝试获取上一行缩进
	indent := "\t"
	if li := strings.LastIndex(inside, "\n"); li >= 0 {
		line := inside[li+1:]
		indent = leadingWhitespace(line)
		if indent == "" {
			indent = "\t"
		}
	}

	// 添加新的项目
	trimmed = strings.TrimSpace(inside)
	if trimmed != "" && !strings.HasSuffix(strings.TrimSpace(trimmed), ",") {
		inside = inside + ",\n"
	}
	inside = inside + indent + item + ",\n"
	out := before + inside + after

	return os.WriteFile(filePath, []byte(out), 0o644)
}

// LeadingWhitespace 获取字符串开头的空白字符
func leadingWhitespace(s string) string {
	for i, r := range s {
		if r != ' ' && r != '\t' {
			return s[:i]
		}
	}
	return s
}

// ExitWith 退出程序
func ExitWith(msg string) {
	PrintError(msg)
	os.Exit(1)
}
