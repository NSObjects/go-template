/*
 * Created by generator on 2025/9/3
 * Enhanced with better UX and detailed templates
 */

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/NSObjects/go-template/tools/modgen/generator"
	"github.com/NSObjects/go-template/tools/modgen/utils"
)

func main() {
	var name string
	var route string
	var force bool
	var openapiFile string
	var generateTests bool
	var generateAll bool

	flag.StringVar(&name, "name", "", "模块名，例如: user, article")
	flag.StringVar(&route, "route", "", "基础路由前缀，例如: /articles (默认使用 name 的复数形式)")
	flag.BoolVar(&force, "force", false, "若目标文件已存在则覆盖")
	flag.StringVar(&openapiFile, "openapi", "", "OpenAPI3文档路径，例如: doc/openapi.yaml")
	flag.BoolVar(&generateTests, "tests", false, "同时生成测试用例（Table-driven测试）")
	flag.BoolVar(&generateAll, "all", false, "生成所有API模块（需要指定OpenAPI文档）")
	flag.Parse()

	if !generateAll && name == "" {
		utils.PrintError("请使用 --name 指定模块名，或使用 --all 生成所有模块")
		fmt.Println("用法: go run tools/modgen/main.go --name=user")
		fmt.Println("或者: go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml")
		fmt.Println("或者: go run tools/modgen/main.go --name=user --openapi=doc/openapi.yaml --tests")
		fmt.Println("或者: go run tools/modgen/main.go --all --openapi=doc/openapi.yaml")
		fmt.Println("或者: go run tools/modgen/main.go --all --openapi=doc/openapi.yaml --tests")
		os.Exit(1)
	}

	if generateAll && openapiFile == "" {
		utils.PrintError("使用 --all 时必须指定 --openapi 文档路径")
		os.Exit(1)
	}

	// 工作空间根目录（tools/modgen 相对）
	cwd, _ := os.Getwd()
	repoRoot := utils.FindRepoRoot(cwd)
	if repoRoot == "" {
		utils.ExitWith("未找到仓库根目录，请在项目内运行")
	}

	// 获取项目包路径
	packagePath, err := utils.GetPackagePath(repoRoot)
	if err != nil {
		utils.ExitWith(fmt.Sprintf("获取项目包路径失败: %v", err))
	}

	// 创建生成器配置
	config := &generator.Config{
		Name:          name,
		Route:         route,
		Force:         force,
		OpenAPIFile:   openapiFile,
		GenerateTests: generateTests,
		PackagePath:   packagePath,
		RepoRoot:      repoRoot,
		GenerateAll:   generateAll,
	}

	// 创建生成器并执行生成
	gen := generator.NewGenerator(config)
	if err := gen.Generate(); err != nil {
		utils.ExitWith(fmt.Sprintf("生成失败: %v", err))
	}
}
