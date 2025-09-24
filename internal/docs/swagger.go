package docs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// SwaggerGenerator Swagger文档生成器
type SwaggerGenerator struct {
	server *echo.Echo
}

// NewSwaggerGenerator 创建Swagger文档生成器
func NewSwaggerGenerator(server *echo.Echo) *SwaggerGenerator {
	return &SwaggerGenerator{server: server}
}

// Generate 生成Swagger文档
func (sg *SwaggerGenerator) Generate(outputPath string) error {
	swagger := sg.buildSwaggerSpec()

	// 确保输出目录存在
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// 生成JSON格式的Swagger文档
	jsonData, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal swagger spec: %w", err)
	}

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("write swagger file: %w", err)
	}

	return nil
}

// buildSwaggerSpec 构建Swagger规范
func (sg *SwaggerGenerator) buildSwaggerSpec() map[string]interface{} {
	swagger := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "Go Template API",
			"description": "基于Go模板的RESTful API服务",
			"version":     "1.0.0",
			"contact": map[string]interface{}{
				"name":  "API Support",
				"email": "support@example.com",
			},
		},
		"servers": []map[string]interface{}{
			{
				"url":         "http://localhost:9322",
				"description": "开发环境",
			},
		},
		"paths":      sg.buildPaths(),
		"components": sg.buildComponents(),
	}

	return swagger
}

// buildPaths 构建API路径
func (sg *SwaggerGenerator) buildPaths() map[string]interface{} {
	paths := make(map[string]interface{})

	for _, route := range sg.server.Routes() {
		path := route.Path
		method := strings.ToLower(route.Method)

		if paths[path] == nil {
			paths[path] = make(map[string]interface{})
		}

		pathItem := paths[path].(map[string]interface{})

		operation := map[string]interface{}{
			"summary":     route.Name,
			"operationId": sg.generateOperationId(route),
			"tags":        []string{sg.extractTag(path)},
		}

		// 添加参数
		if strings.Contains(path, ":") {
			operation["parameters"] = sg.buildParameters(path)
		}

		// 添加响应
		operation["responses"] = sg.buildResponses(method)

		pathItem[method] = operation
	}

	return paths
}

// buildParameters 构建路径参数
func (sg *SwaggerGenerator) buildParameters(path string) []map[string]interface{} {
	var params []map[string]interface{}

	// 简单的参数提取，实际项目中可能需要更复杂的解析
	if strings.Contains(path, ":id") {
		params = append(params, map[string]interface{}{
			"name":        "id",
			"in":          "path",
			"required":    true,
			"description": "资源ID",
			"schema": map[string]interface{}{
				"type":   "integer",
				"format": "int64",
			},
		})
	}

	return params
}

// buildResponses 构建响应定义
func (sg *SwaggerGenerator) buildResponses(method string) map[string]interface{} {
	responses := map[string]interface{}{
		"200": map[string]interface{}{
			"description": "成功",
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"type": "object",
					},
				},
			},
		},
		"400": map[string]interface{}{
			"description": "请求错误",
		},
		"500": map[string]interface{}{
			"description": "服务器错误",
		},
	}

	// 为特定方法添加特定响应
	switch method {
	case "post":
		responses["201"] = map[string]interface{}{
			"description": "创建成功",
		}
	case "delete":
		responses["204"] = map[string]interface{}{
			"description": "删除成功",
		}
	}

	return responses
}

// buildComponents 构建组件定义
func (sg *SwaggerGenerator) buildComponents() map[string]interface{} {
	return map[string]interface{}{
		"schemas": map[string]interface{}{
			"User": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":   "integer",
						"format": "int64",
					},
					"username": map[string]interface{}{
						"type": "string",
					},
					"email": map[string]interface{}{
						"type":   "string",
						"format": "email",
					},
					"age": map[string]interface{}{
						"type": "integer",
					},
					"created_at": map[string]interface{}{
						"type":   "string",
						"format": "date-time",
					},
					"updated_at": map[string]interface{}{
						"type":   "string",
						"format": "date-time",
					},
				},
			},
		},
	}
}

// generateOperationId 生成操作ID
func (sg *SwaggerGenerator) generateOperationId(route *echo.Route) string {
	method := strings.ToLower(route.Method)
	path := strings.ReplaceAll(route.Path, ":", "")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.Trim(path, "_")

	return fmt.Sprintf("%s_%s", method, path)
}

// extractTag 提取标签
func (sg *SwaggerGenerator) extractTag(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		return strings.Title(parts[0])
	}
	return "Default"
}
