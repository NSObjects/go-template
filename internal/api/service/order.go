/*
 * Generated from OpenAPI3 document
 * Module: Order
 */

package service

import (
	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/NSObjects/go-template/internal/utils"
	"github.com/labstack/echo/v4"
)

type OrderController struct {
	order biz.OrderUseCase
}

func NewOrderController(h biz.OrderUseCase) RegisterRouter {
	return &OrderController{ order: h }
}

func (c *OrderController) RegisterRouter(g *echo.Group, m ...echo.MiddlewareFunc) {

	g.GET("/orders", c.ListOrders).Name = "获取订单列表"

	g.POST("/orders", c.CreateOrder).Name = "创建订单"

}

// TODO: 实现控制器方法



func (c *OrderController) ListOrders(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.OrderListOrdersRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	list, total, err := c.order.ListOrders(bizCtx, req)
	if err != nil {
		return err
	}
	
	// 返回列表数据
	return resp.ListDataResponse(list, total, ctx)
}



func (c *OrderController) CreateOrder(ctx echo.Context) error {
	// TODO: 绑定和验证请求参数
	var req param.OrderCreateOrderRequest
	if err := BindAndValidate(ctx, &req); err != nil {
		return err
	}
	
	// 调用业务逻辑 - 构造包含链路追踪信息的context
	bizCtx := utils.BuildContext(ctx)
	if err := c.order.CreateOrder(bizCtx, req); err != nil {
		return err
	}
	
	// 返回操作成功
	return resp.OperateSuccess(ctx)
}


