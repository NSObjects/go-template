package code

// ErrorType 错误类型枚举
type ErrorType int

const (
	// InternalError 内部错误：数据库、Redis、Kafka等组件错误
	// 这些错误会记录详细日志，对外返回HTTP 500
	InternalError ErrorType = iota

	// BusinessError 业务错误：参数验证、权限、业务逻辑错误
	// 这些错误有明确的HTTP状态码和业务错误码，给前端明确提示
	BusinessError
)

// ErrorCategory 错误分类
type ErrorCategory string

const (
	// 内部错误分类
	CategoryDatabase ErrorCategory = "database" // 数据库错误
	CategoryRedis    ErrorCategory = "redis"    // Redis错误
	CategoryKafka    ErrorCategory = "kafka"    // Kafka错误
	CategoryExternal ErrorCategory = "external" // 第三方服务错误
	CategorySystem   ErrorCategory = "system"   // 系统错误

	// 业务错误分类
	CategoryAuth       ErrorCategory = "auth"       // 认证错误
	CategoryPermission ErrorCategory = "permission" // 权限错误
	CategoryValidation ErrorCategory = "validation" // 参数验证错误
	CategoryBusiness   ErrorCategory = "business"   // 业务逻辑错误
)

// ErrorInfo 错误信息结构
type ErrorInfo struct {
	Type     ErrorType     `json:"type"`     // 错误类型
	Category ErrorCategory `json:"category"` // 错误分类
	Code     int           `json:"code"`     // 业务错误码
	Message  string        `json:"message"`  // 错误消息
	Details  string        `json:"details"`  // 详细错误信息（内部错误时记录）
}

// IsInternal 判断是否为内部错误
func (e *ErrorInfo) IsInternal() bool {
	return e.Type == InternalError
}

// IsBusiness 判断是否为业务错误
func (e *ErrorInfo) IsBusiness() bool {
	return e.Type == BusinessError
}
