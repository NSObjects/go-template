package param

import "time"

// ProductParam 查询参数
type ProductParam struct {
	Page  int    `json:"page" form:"page" query:"page"`
	Count int    `json:"count" form:"count" query:"count"`
	Name  string `json:"name" form:"name" query:"name"`
	// TODO: 添加更多查询字段
}

// Limit 获取限制数量
func (p ProductParam) Limit() int {
	if p.Count <= 0 {
		return 10
	}
	return p.Count
}

// Offset 获取偏移量
func (p ProductParam) Offset() int {
	if p.Page <= 1 {
		return 0
	}
	return (p.Page - 1) * p.Limit()
}

// ProductBody 创建/更新请求体
type ProductBody struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	// TODO: 添加更多字段
}

// ProductResponse 响应结构
type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// TODO: 添加更多返回字段
}
