/*
 * OpenAPI3 类型定义
 */

package openapi

// OpenAPI3 结构定义
type OpenAPI3 struct {
	OpenAPI    string              `yaml:"openapi" json:"openapi"`
	Info       Info                `yaml:"info" json:"info"`
	Servers    []Server            `yaml:"servers,omitempty" json:"servers,omitempty"`
	Paths      map[string]PathItem `yaml:"paths" json:"paths"`
	Components Components          `yaml:"components,omitempty" json:"components,omitempty"`
	Tags       []Tag               `yaml:"tags,omitempty" json:"tags,omitempty"`
}

type Info struct {
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Version     string `yaml:"version" json:"version"`
}

type Server struct {
	URL         string `yaml:"url" json:"url"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

type PathItem struct {
	Get    *Operation `yaml:"get,omitempty" json:"get,omitempty"`
	Post   *Operation `yaml:"post,omitempty" json:"post,omitempty"`
	Put    *Operation `yaml:"put,omitempty" json:"put,omitempty"`
	Delete *Operation `yaml:"delete,omitempty" json:"delete,omitempty"`
	Patch  *Operation `yaml:"patch,omitempty" json:"patch,omitempty"`
}

type Operation struct {
	Tags        []string            `yaml:"tags,omitempty" json:"tags,omitempty"`
	Summary     string              `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description string              `yaml:"description,omitempty" json:"description,omitempty"`
	OperationID string              `yaml:"operationId,omitempty" json:"operationId,omitempty"`
	Parameters  []Parameter         `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	RequestBody *RequestBody        `yaml:"requestBody,omitempty" json:"requestBody,omitempty"`
	Responses   map[string]Response `yaml:"responses" json:"responses"`
}

type Parameter struct {
	Name        string  `yaml:"name" json:"name"`
	In          string  `yaml:"in" json:"in"`
	Description string  `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool    `yaml:"required,omitempty" json:"required,omitempty"`
	Schema      *Schema `yaml:"schema,omitempty" json:"schema,omitempty"`
	Ref         string  `yaml:"$ref,omitempty" json:"$ref,omitempty"`
}

type RequestBody struct {
	Description string               `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool                 `yaml:"required,omitempty" json:"required,omitempty"`
	Content     map[string]MediaType `yaml:"content" json:"content"`
}

type Response struct {
	Description string               `yaml:"description" json:"description"`
	Content     map[string]MediaType `yaml:"content,omitempty" json:"content,omitempty"`
	Ref         string               `yaml:"$ref,omitempty" json:"$ref,omitempty"`
}

type MediaType struct {
	Schema *Schema `yaml:"schema,omitempty" json:"schema,omitempty"`
}

type Schema struct {
	Type        string             `yaml:"type,omitempty" json:"type,omitempty"`
	Format      string             `yaml:"format,omitempty" json:"format,omitempty"`
	Description string             `yaml:"description,omitempty" json:"description,omitempty"`
	Properties  map[string]*Schema `yaml:"properties,omitempty" json:"properties,omitempty"`
	Items       *Schema            `yaml:"items,omitempty" json:"items,omitempty"`
	Required    []string           `yaml:"required,omitempty" json:"required,omitempty"`
	Ref         string             `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Default     interface{}        `yaml:"default,omitempty" json:"default,omitempty"`
	// 验证相关字段
	MinLength *int     `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	MaxLength *int     `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	Pattern   string   `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	Minimum   *float64 `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	Maximum   *float64 `yaml:"maximum,omitempty" json:"maximum,omitempty"`
	MinItems  *int     `yaml:"minItems,omitempty" json:"minItems,omitempty"`
	MaxItems  *int     `yaml:"maxItems,omitempty" json:"maxItems,omitempty"`
}

type Components struct {
	Schemas       map[string]*Schema      `yaml:"schemas,omitempty" json:"schemas,omitempty"`
	Parameters    map[string]*Parameter   `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Responses     map[string]*Response    `yaml:"responses,omitempty" json:"responses,omitempty"`
	RequestBodies map[string]*RequestBody `yaml:"requestBodies,omitempty" json:"requestBodies,omitempty"`
}

type Tag struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

// API模块信息
type APIModule struct {
	Name          string
	Tag           string
	Operations    []APIOperation
	Schemas       map[string]*Schema
	Parameters    map[string]*Parameter
	RequestBodies map[string]*RequestBody
}

type APIOperation struct {
	Method      string
	Path        string
	Summary     string
	Description string
	OperationID string
	Parameters  []Parameter
	RequestBody *RequestBody
	Responses   map[string]Response
	Tag         string
}
