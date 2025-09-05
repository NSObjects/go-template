package biz

import (
	"context"
	"github.com/NSObjects/go-template/internal/api/data"
	"github.com/NSObjects/go-template/internal/api/data/model"
	"github.com/NSObjects/go-template/internal/api/service/param"
)

// ProductUseCase Product业务用例接口
type ProductUseCase interface {
	List(ctx context.Context, p param.ProductParam) ([]param.ProductResponse, int64, error)
	Create(ctx context.Context, b param.ProductBody) error
	Update(ctx context.Context, id int64, b param.ProductBody) error
	Delete(ctx context.Context, id int64) error
	Detail(ctx context.Context, id int64) (*param.ProductResponse, error)
}

// ProductHandler Product业务处理器
type ProductHandler struct {
	dm *data.DataManager
}

// NewProductHandler 创建Product业务处理器
func NewProductHandler(dm *data.DataManager) *ProductHandler {
	return &ProductHandler{dm: dm}
}

// List 获取Product列表
func (h *ProductHandler) List(ctx context.Context, p param.ProductParam) ([]param.ProductResponse, int64, error) {
	// TODO: 实现列表查询逻辑
	// 示例：
	// var models []model.Product
	// if err := h.dm.MySQLWithContext(ctx).Offset(p.Offset()).Limit(p.Limit()).Find(&models).Error; err != nil {
	//     return nil, 0, code.WrapDatabaseError(err, "query Product list")
	// }
	// var total int64
	// h.dm.MySQLWithContext(ctx).Model(&model.Product{}).Count(&total)
	// return convertProductToResponses(models), total, nil
	return nil, 0, nil
}

// Create 创建Product
func (h *ProductHandler) Create(ctx context.Context, b param.ProductBody) error {
	// TODO: 实现创建逻辑
	// 示例：
	// model := &model.Product{
	//     // 设置字段
	//     CreatedAt: time.Now(),
	// }
	// if err := h.dm.MySQLWithContext(ctx).Create(model).Error; err != nil {
	//     return nil, code.WrapDatabaseError(err, "create Product")
	// }
	// 创建成功
	return nil
}

// Update 更新Product
func (h *ProductHandler) Update(ctx context.Context, id int64, b param.ProductBody) error {
	// TODO: 实现更新逻辑
	// 示例：
	// var model model.Product
	// if err := h.dm.MySQLWithContext(ctx).First(&model, id).Error; err != nil {
	//     if errors.Is(err, gorm.ErrRecordNotFound) {
	//         return nil, code.WrapNotFoundError(nil, "Product not found")
	//     }
	//     return nil, code.WrapDatabaseError(err, "query Product")
	// }
	// // 更新字段
	// model.UpdatedAt = time.Now()
	// if err := h.dm.MySQLWithContext(ctx).Save(&model).Error; err != nil {
	//     return code.WrapDatabaseError(err, "update Product")
	// }
	// 更新成功
	return nil
}

// Delete 删除Product
func (h *ProductHandler) Delete(ctx context.Context, id int64) error {
	// TODO: 实现删除逻辑
	// 示例：
	// if err := h.dm.MySQLWithContext(ctx).Delete(&model.Product{}, id).Error; err != nil {
	//     return code.WrapDatabaseError(err, "delete Product")
	// }
	// return nil
	return nil
}

// Detail 获取Product详情
func (h *ProductHandler) Detail(ctx context.Context, id int64) (*param.ProductResponse, error) {
	// TODO: 实现详情查询逻辑
	// 示例：
	// var model model.Product
	// if err := h.dm.MySQLWithContext(ctx).First(&model, id).Error; err != nil {
	//     if errors.Is(err, gorm.ErrRecordNotFound) {
	//         return nil, code.WrapNotFoundError(nil, "Product not found")
	//     }
	//     return nil, code.WrapDatabaseError(err, "query Product")
	// }
	// return convertProductToResponse(&model), nil
	return nil, nil
}

// convertProductToResponse 转换为响应结构
func convertProductToResponse(model *model.Product) *param.ProductResponse {
	// TODO: 实现转换逻辑
	return &param.ProductResponse{
		// ID: model.ID,
		// CreatedAt: model.CreatedAt,
		// UpdatedAt: model.UpdatedAt,
	}
}

// convertProductToResponses 转换为响应结构列表
func convertProductToResponses(models []model.Product) []param.ProductResponse {
	responses := make([]param.ProductResponse, len(models))
	for i, model := range models {
		responses[i] = *convertProductToResponse(&model)
	}
	return responses
}
