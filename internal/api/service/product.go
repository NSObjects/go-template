package service

import (
	"strconv"
	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/NSObjects/go-template/internal/utils"
	"github.com/labstack/echo/v4"
)

type productController struct {
	product biz.ProductUseCase
}

func NewProductController(h *biz.ProductHandler) RegisterRouter {
	return &productController{
		product: h,
	}
}

func (c *productController) RegisterRouter(g *echo.Group, m ...echo.MiddlewareFunc) {
	g.GET("/products", c.list).Name = "列表示例"
	g.POST("/products", c.create).Name = "创建示例"
	g.GET("/products/:id", c.detail).Name = "详情示例"
	g.PUT("/products/:id", c.update).Name = "更新示例"
	g.DELETE("/products/:id", c.remove).Name = "删除示例"
}

func (c *productController) list(ctx echo.Context) error {
	var p param.ProductParam
	if err := BindAndValidate(ctx, &p); err != nil { return err }
	bizCtx := utils.BuildContext(ctx)
	items, total, err := c.product.List(bizCtx, p)
	if err != nil { return err }
	return resp.ListDataResponse(items, total, ctx)
}

func (c *productController) detail(ctx echo.Context) error {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	bizCtx := utils.BuildContext(ctx)
	item, err := c.product.Detail(bizCtx, id)
	if err != nil { return err }
	return resp.OneDataResponse(item, ctx)
}

func (c *productController) create(ctx echo.Context) error {
	var b param.ProductBody
	if err := BindAndValidate(ctx, &b); err != nil { return err }
	bizCtx := utils.BuildContext(ctx)
	if err := c.product.Create(bizCtx, b); err != nil { return err }
	return resp.OperateSuccess(ctx)
}

func (c *productController) update(ctx echo.Context) error {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	var b param.ProductBody
	if err := BindAndValidate(ctx, &b); err != nil { return err }
	bizCtx := utils.BuildContext(ctx)
	if err := c.product.Update(bizCtx, id, b); err != nil { return err }
	return resp.OperateSuccess(ctx)
}

func (c *productController) remove(ctx echo.Context) error {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	bizCtx := utils.BuildContext(ctx)
	if err := c.product.Delete(bizCtx, id); err != nil { return err }
	return resp.OperateSuccess(ctx)
}
