/*
 * 默认模板渲染函数
 */

package templates

import (
	"fmt"
)

// RenderBiz 生成默认业务逻辑模板
func RenderBiz(pascal, packagePath string) string {
	return "package biz\n\n" +
		"import (\n" +
		"\t\"context\"\n" +
		fmt.Sprintf("\t\"%s/internal/api/data\"\n", packagePath) +
		fmt.Sprintf("\t\"%s/internal/api/data/model\"\n", packagePath) +
		fmt.Sprintf("\t\"%s/internal/api/service/param\"\n", packagePath) +
		")\n\n" +
		"// " + pascal + "UseCase " + pascal + "业务用例接口\n" +
		fmt.Sprintf("type %sUseCase interface {\n", pascal) +
		fmt.Sprintf("\tList(ctx context.Context, p param.%[1]sParam) ([]param.%[1]sResponse, int64, error)\n", pascal) +
		fmt.Sprintf("\tCreate(ctx context.Context, b param.%[1]sBody) error\n", pascal) +
		fmt.Sprintf("\tUpdate(ctx context.Context, id int64, b param.%[1]sBody) error\n", pascal) +
		"\tDelete(ctx context.Context, id int64) error\n" +
		fmt.Sprintf("\tDetail(ctx context.Context, id int64) (*param.%[1]sResponse, error)\n", pascal) +
		"}\n\n" +
		"// " + pascal + "Handler " + pascal + "业务处理器\n" +
		fmt.Sprintf("type %sHandler struct {\n", pascal) +
		"\tdm *data.DataManager\n" +
		"}\n\n" +
		fmt.Sprintf("// New%[1]sHandler 创建%[1]s业务处理器\n", pascal) +
		fmt.Sprintf("func New%[1]sHandler(dm *data.DataManager) *%[1]sHandler {\n", pascal) +
		fmt.Sprintf("\treturn &%[1]sHandler{dm: dm}\n", pascal) +
		"}\n\n" +
		"// List 获取" + pascal + "列表\n" +
		fmt.Sprintf("func (h *%[1]sHandler) List(ctx context.Context, p param.%[1]sParam) ([]param.%[1]sResponse, int64, error) {\n", pascal) +
		"\t// TODO: 实现列表查询逻辑\n" +
		"\t// 示例：\n" +
		"\t// var models []model." + pascal + "\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).Offset(p.Offset()).Limit(p.Limit()).Find(&models).Error; err != nil {\n" +
		"\t//     return nil, 0, code.WrapDatabaseError(err, \"query " + pascal + " list\")\n" +
		"\t// }\n" +
		"\t// var total int64\n" +
		"\t// h.dm.MySQLWithContext(ctx).Model(&model." + pascal + "{}).Count(&total)\n" +
		"\t// return convert%sToResponses(models), total, nil\n" +
		"\treturn nil, 0, nil\n" +
		"}\n\n" +
		"// Create 创建" + pascal + "\n" +
		fmt.Sprintf("func (h *%[1]sHandler) Create(ctx context.Context, b param.%[1]sBody) error {\n", pascal) +
		"\t// TODO: 实现创建逻辑\n" +
		"\t// 示例：\n" +
		"\t// model := &model." + pascal + "{\n" +
		"\t//     // 设置字段\n" +
		"\t//     CreatedAt: time.Now(),\n" +
		"\t// }\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).Create(model).Error; err != nil {\n" +
		"\t//     return nil, code.WrapDatabaseError(err, \"create " + pascal + "\")\n" +
		"\t// }\n" +
		"\t// 创建成功\n" +
		"\treturn nil\n" +
		"}\n\n" +
		"// Update 更新" + pascal + "\n" +
		fmt.Sprintf("func (h *%[1]sHandler) Update(ctx context.Context, id int64, b param.%[1]sBody) error {\n", pascal) +
		"\t// TODO: 实现更新逻辑\n" +
		"\t// 示例：\n" +
		"\t// var model model." + pascal + "\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).First(&model, id).Error; err != nil {\n" +
		"\t//     if errors.Is(err, gorm.ErrRecordNotFound) {\n" +
		"\t//         return nil, code.WrapNotFoundError(nil, \"" + pascal + " not found\")\n" +
		"\t//     }\n" +
		"\t//     return nil, code.WrapDatabaseError(err, \"query " + pascal + "\")\n" +
		"\t// }\n" +
		"\t// // 更新字段\n" +
		"\t// model.UpdatedAt = time.Now()\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).Save(&model).Error; err != nil {\n" +
		"\t//     return code.WrapDatabaseError(err, \"update " + pascal + "\")\n" +
		"\t// }\n" +
		"\t// 更新成功\n" +
		"\treturn nil\n" +
		"}\n\n" +
		"// Delete 删除" + pascal + "\n" +
		fmt.Sprintf("func (h *%[1]sHandler) Delete(ctx context.Context, id int64) error {\n", pascal) +
		"\t// TODO: 实现删除逻辑\n" +
		"\t// 示例：\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).Delete(&model." + pascal + "{}, id).Error; err != nil {\n" +
		"\t//     return code.WrapDatabaseError(err, \"delete " + pascal + "\")\n" +
		"\t// }\n" +
		"\t// return nil\n" +
		"\treturn nil\n" +
		"}\n\n" +
		"// Detail 获取" + pascal + "详情\n" +
		fmt.Sprintf("func (h *%[1]sHandler) Detail(ctx context.Context, id int64) (*param.%[1]sResponse, error) {\n", pascal) +
		"\t// TODO: 实现详情查询逻辑\n" +
		"\t// 示例：\n" +
		"\t// var model model." + pascal + "\n" +
		"\t// if err := h.dm.MySQLWithContext(ctx).First(&model, id).Error; err != nil {\n" +
		"\t//     if errors.Is(err, gorm.ErrRecordNotFound) {\n" +
		"\t//         return nil, code.WrapNotFoundError(nil, \"" + pascal + " not found\")\n" +
		"\t//     }\n" +
		"\t//     return nil, code.WrapDatabaseError(err, \"query " + pascal + "\")\n" +
		"\t// }\n" +
		"\t// return convert%sToResponse(&model), nil\n" +
		"\treturn nil, nil\n" +
		"}\n\n" +
		"// convert%sToResponse 转换为响应结构\n" +
		fmt.Sprintf("func convert%sToResponse(model *model.%[1]s) *param.%[1]sResponse {\n", pascal, pascal) +
		"\t// TODO: 实现转换逻辑\n" +
		"\treturn &param." + pascal + "Response{\n" +
		"\t\t// ID: model.ID,\n" +
		"\t\t// CreatedAt: model.CreatedAt,\n" +
		"\t\t// UpdatedAt: model.UpdatedAt,\n" +
		"\t}\n" +
		"}\n\n" +
		"// convert%sToResponses 转换为响应结构列表\n" +
		fmt.Sprintf("func convert%sToResponses(models []model.%[1]s) []param.%[1]sResponse {\n", pascal, pascal) +
		"\tresponses := make([]param." + pascal + "Response, len(models))\n" +
		"\tfor i, model := range models {\n" +
		"\t\tresponses[i] = *convert" + pascal + "ToResponse(&model)\n" +
		"\t}\n" +
		"\treturn responses\n" +
		"}\n"
}

// RenderService 生成默认服务层模板
func RenderService(pascal, camel, baseRoute, packagePath string) string {
	return "package service\n\n" +
		"import (\n" +
		"\t\"strconv\"\n" +
		fmt.Sprintf("\t\"%s/internal/api/biz\"\n", packagePath) +
		fmt.Sprintf("\t\"%s/internal/api/service/param\"\n", packagePath) +
		fmt.Sprintf("\t\"%s/internal/resp\"\n", packagePath) +
		fmt.Sprintf("\t\"%s/internal/utils\"\n", packagePath) +
		"\t\"github.com/labstack/echo/v4\"\n" +
		")\n\n" +
		fmt.Sprintf("type %sController struct {\n\t%s biz.%sUseCase\n}\n\n", camel, camel, pascal) +
		fmt.Sprintf("func New%[1]sController(h *biz.%[1]sHandler) RegisterRouter {\n\treturn &%[2]sController{%[2]s: h}\n}\n\n", pascal, camel) +
		fmt.Sprintf("func (c *%[1]sController) RegisterRouter(g *echo.Group, m ...echo.MiddlewareFunc) {\n\tg.GET(\"%[2]s\", c.list).Name = \"列表示例\"\n\tg.POST(\"%[2]s\", c.create).Name = \"创建示例\"\n\tg.GET(\"%[2]s/:id\", c.detail).Name = \"详情示例\"\n\tg.PUT(\"%[2]s/:id\", c.update).Name = \"更新示例\"\n\tg.DELETE(\"%[2]s/:id\", c.remove).Name = \"删除示例\"\n}\n\n", camel, baseRoute) +
		fmt.Sprintf("func (c *%[1]sController) list(ctx echo.Context) error {\n\tvar p param.%[2]sParam\n\tif err := BindAndValidate(ctx, &p); err != nil { return err }\n\tbizCtx := utils.BuildContext(ctx)\n\titems, total, err := c.%[3]s.List(bizCtx, p)\n\tif err != nil { return err }\n\treturn resp.ListDataResponse(items, total, ctx)\n}\n\n", camel, pascal, camel) +
		fmt.Sprintf("func (c *%[1]sController) detail(ctx echo.Context) error {\n\tid, _ := strconv.ParseInt(ctx.Param(\"id\"), 10, 64)\n\tbizCtx := utils.BuildContext(ctx)\n\titem, err := c.%[2]s.Detail(bizCtx, id)\n\tif err != nil { return err }\n\treturn resp.OneDataResponse(item, ctx)\n}\n\n", camel, camel) +
		fmt.Sprintf("func (c *%[1]sController) create(ctx echo.Context) error {\n\tvar b param.%[2]sBody\n\tif err := BindAndValidate(ctx, &b); err != nil { return err }\n\tbizCtx := utils.BuildContext(ctx)\n\tif err := c.%[3]s.Create(bizCtx, b); err != nil { return err }\n\treturn resp.OperateSuccess(ctx)\n}\n\n", camel, pascal, camel) +
		fmt.Sprintf("func (c *%[1]sController) update(ctx echo.Context) error {\n\tid, _ := strconv.ParseInt(ctx.Param(\"id\"), 10, 64)\n\tvar b param.%[2]sBody\n\tif err := BindAndValidate(ctx, &b); err != nil { return err }\n\tbizCtx := utils.BuildContext(ctx)\n\tif err := c.%[3]s.Update(bizCtx, id, b); err != nil { return err }\n\treturn resp.OperateSuccess(ctx)\n}\n\n", camel, pascal, camel) +
		fmt.Sprintf("func (c *%[1]sController) remove(ctx echo.Context) error {\n\tid, _ := strconv.ParseInt(ctx.Param(\"id\"), 10, 64)\n\tbizCtx := utils.BuildContext(ctx)\n\tif err := c.%[2]s.Delete(bizCtx, id); err != nil { return err }\n\treturn resp.OperateSuccess(ctx)\n}\n", camel, camel)
}

// RenderParam 生成默认参数模板
func RenderParam(pascal, packagePath string) string {
	return "package param\n\n" +
		"import \"time\"\n\n" +
		"// " + pascal + "Param 查询参数\n" +
		fmt.Sprintf("type %sParam struct {\n", pascal) +
		"\tPage  int    `json:\"page\" form:\"page\" query:\"page\"`\n" +
		"\tCount int    `json:\"count\" form:\"count\" query:\"count\"`\n" +
		"\tName  string `json:\"name\" form:\"name\" query:\"name\"`\n" +
		"\t// TODO: 添加更多查询字段\n" +
		"}\n\n" +
		"// Limit 获取限制数量\n" +
		fmt.Sprintf("func (p %sParam) Limit() int {\n", pascal) +
		"\tif p.Count <= 0 {\n" +
		"\t\treturn 10\n" +
		"\t}\n" +
		"\treturn p.Count\n" +
		"}\n\n" +
		"// Offset 获取偏移量\n" +
		fmt.Sprintf("func (p %sParam) Offset() int {\n", pascal) +
		"\tif p.Page <= 1 {\n" +
		"\t\treturn 0\n" +
		"\t}\n" +
		"\treturn (p.Page - 1) * p.Limit()\n" +
		"}\n\n" +
		"// " + pascal + "Body 创建/更新请求体\n" +
		fmt.Sprintf("type %sBody struct {\n", pascal) +
		"\tName        string `json:\"name\" validate:\"required\"`\n" +
		"\tDescription string `json:\"description\"`\n" +
		"\t// TODO: 添加更多字段\n" +
		"}\n\n" +
		"// " + pascal + "Response 响应结构\n" +
		fmt.Sprintf("type %sResponse struct {\n", pascal) +
		"\tID          uint      `json:\"id\"`\n" +
		"\tName        string    `json:\"name\"`\n" +
		"\tDescription string    `json:\"description\"`\n" +
		"\tCreatedAt   time.Time `json:\"created_at\"`\n" +
		"\tUpdatedAt   time.Time `json:\"updated_at\"`\n" +
		"\t// TODO: 添加更多返回字段\n" +
		"}\n"
}

// RenderModel 生成默认数据模型模板
func RenderModel(pascal, table, packagePath string) string {
	return "package model\n\n" +
		"import (\n" +
		"\t\"time\"\n" +
		"\t\"gorm.io/gorm\"\n" +
		")\n\n" +
		"// " + pascal + " 数据模型\n" +
		fmt.Sprintf("type %s struct {\n", pascal) +
		"\tID          uint           `gorm:\"primaryKey;autoIncrement\" json:\"id\"`\n" +
		"\tName        string         `gorm:\"column:name;type:varchar(100);not null\" json:\"name\"`\n" +
		"\tDescription string         `gorm:\"column:description;type:text\" json:\"description\"`\n" +
		"\tStatus      int            `gorm:\"column:status;type:int;default:1\" json:\"status\"`\n" +
		"\tCreatedAt   time.Time      `gorm:\"column:created_at\" json:\"created_at\"`\n" +
		"\tUpdatedAt   time.Time      `gorm:\"column:updated_at\" json:\"updated_at\"`\n" +
		"\tDeletedAt   gorm.DeletedAt `gorm:\"column:deleted_at;index\" json:\"-\"`\n" +
		"\t// TODO: 添加更多字段\n" +
		"}\n\n" +
		"// TableName 指定表名\n" +
		fmt.Sprintf("func (%s) TableName() string {\n", pascal) +
		fmt.Sprintf("\treturn \"%s\"\n", table) +
		"}\n"
}

// RenderCode 生成业务错误码文件
func RenderCode(pascal, table, packagePath string) string {
	// 计算错误码起始值（基于表名）
	baseCode := 100000 + int(table[0])*1000 + int(table[len(table)-1])*10

	return "package code\n\n" +
		"//go:generate codegen -type=int\n" +
		"//go:generate codegen -type=int -doc -output ./error_code_generated.md\n\n" +
		fmt.Sprintf("// %s相关错误码\n", pascal) +
		"const (\n" +
		fmt.Sprintf("\t// Err%sNotFound - 404: %s not found.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sNotFound int = iota + %d\n", pascal, baseCode) +
		fmt.Sprintf("\t// Err%sAlreadyExists - 400: %s already exists.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sAlreadyExists\n", pascal) +
		fmt.Sprintf("\t// Err%sInvalidData - 400: %s invalid data.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sInvalidData\n", pascal) +
		fmt.Sprintf("\t// Err%sPermissionDenied - 403: %s permission denied.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sPermissionDenied\n", pascal) +
		fmt.Sprintf("\t// Err%sInUse - 400: %s is in use.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sInUse\n", pascal) +
		fmt.Sprintf("\t// Err%sCreateFailed - 500: %s create failed.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sCreateFailed\n", pascal) +
		fmt.Sprintf("\t// Err%sUpdateFailed - 500: %s update failed.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sUpdateFailed\n", pascal) +
		fmt.Sprintf("\t// Err%sDeleteFailed - 500: %s delete failed.\n", pascal, pascal) +
		fmt.Sprintf("\tErr%sDeleteFailed\n", pascal) +
		")\n"
}
