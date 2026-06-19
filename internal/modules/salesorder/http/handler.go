package salesorderhttp

import (
	"github.com/labstack/echo/v5"

	"github.com/NSObjects/go-template/internal/modules/salesorder/usecase"
	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/server/httpreq"
	"github.com/NSObjects/go-template/internal/platform/server/httpresp"
)

const defaultPageSize = 20

// Handler adapts sales order HTTP requests to the sales order usecase.
type Handler struct {
	usecase *usecase.Usecase
}

// New creates a sales order HTTP handler.
func New(uc *usecase.Usecase) *Handler {
	return &Handler{usecase: uc}
}

// Register mounts sales order routes on group.
func Register(group *echo.Group, handler *Handler) {
	group.POST("/sales-orders", handler.Create)
	group.GET("/sales-orders", handler.List)
	group.GET("/sales-orders/:id", handler.Get)
}

type createOrderRequest struct {
	CustomerID int64             `json:"customer_id" validate:"required,gt=0"`
	Items      []createItemInput `json:"items" validate:"required,min=1,dive"`
}

type createItemInput struct {
	ProductID      int64 `json:"product_id" validate:"required,gt=0"`
	Quantity       int   `json:"quantity" validate:"required,gt=0"`
	UnitPriceCents int64 `json:"unit_price_cents" validate:"gte=0"`
}

// Create handles sales order creation requests.
func (h *Handler) Create(c *echo.Context) error {
	if err := h.ready(); err != nil {
		return err
	}
	var req createOrderRequest
	if err := httpreq.BindAndValidate(c, &req); err != nil {
		return err
	}
	items := make([]usecase.CreateItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, usecase.CreateItemInput{
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
		})
	}
	order, err := h.usecase.Create(c.Request().Context(), usecase.CreateInput{
		CustomerID: req.CustomerID,
		Items:      items,
	})
	if err != nil {
		return err
	}
	return httpresp.Created(c, order)
}

// Get handles sales order lookup requests.
func (h *Handler) Get(c *echo.Context) error {
	if err := h.ready(); err != nil {
		return err
	}
	id, err := httpreq.PathID(c, "id", "order")
	if err != nil {
		return err
	}
	order, err := h.usecase.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return httpresp.OK(c, order)
}

// List handles sales order list requests.
func (h *Handler) List(c *echo.Context) error {
	if err := h.ready(); err != nil {
		return err
	}
	page, pageSize, err := httpreq.Pagination(c, defaultPageSize)
	if err != nil {
		return err
	}
	customerID, err := httpreq.QueryInt64(c, "customer_id", 0)
	if err != nil {
		return err
	}
	output, err := h.usecase.List(c.Request().Context(), usecase.ListInput{
		Page:       page,
		PageSize:   pageSize,
		CustomerID: customerID,
	})
	if err != nil {
		return err
	}
	pageMeta, err := httpresp.NewPageMeta(output.Page, output.PageSize, output.Total)
	if err != nil {
		return err
	}
	return httpresp.List(c, output.Items, pageMeta)
}

func (h *Handler) ready() error {
	if h == nil || h.usecase == nil {
		return apperr.New(apperr.ErrInternalServer, "sales order usecase is not configured")
	}
	return nil
}
