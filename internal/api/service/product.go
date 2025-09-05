/*
 * Generated from OpenAPI3 document
 * Module: Product
 */

package service

import (
	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/NSObjects/go-template/internal/utils"
	"github.com/labstack/echo/v4"
)

type ProductController struct {
	product biz.ProductUseCase
}

func NewProductController(h biz.ProductUseCase) RegisterRouter {
	return &ProductController{ product: h }
}

func (c *ProductController) RegisterRouter(g *echo.Group, m ...echo.MiddlewareFunc) {

	g.GET("/products", c.ListProducts).Name = "获取产品列表"

	g.POST("/products", c.CreateProduct).Name = "创建产品"

}

// TODO: 实现控制器方法



func (c *ProductController) ListProducts(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.ProductListProductsRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	list, total, err := c.product.ListProducts(bizCtx, req)
	if err != nil {
		return err
	}
	
	// 返回列表数据
	return resp.ListDataResponse(list, total, ctx)
}



func (c *ProductController) CreateProduct(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.ProductCreateProductRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	if err := c.product.CreateProduct(bizCtx, req); err != nil {
		return err
	}
	
	// 返回操作成功
	return resp.OperateSuccess(ctx)
}


