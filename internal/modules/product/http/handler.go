package producthttp

import (
	"github.com/labstack/echo/v5"

	"github.com/NSObjects/go-template/internal/modules/product/usecase"
	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/server/httpreq"
	"github.com/NSObjects/go-template/internal/platform/server/httpresp"
)

const defaultPageSize = 20

// Handler adapts product HTTP requests to the product usecase.
type Handler struct {
	usecase *usecase.Usecase
}

// New creates a product HTTP handler.
func New(uc *usecase.Usecase) *Handler {
	return &Handler{usecase: uc}
}

// Register mounts product routes on group.
func Register(group *echo.Group, handler *Handler) {
	group.POST("/products", handler.Create)
	group.GET("/products", handler.List)
	group.GET("/products/:id", handler.Get)
	group.PATCH("/products/:id", handler.Update)
}

type createProductRequest struct {
	SKU        string `json:"sku" validate:"required,max=64"`
	Name       string `json:"name" validate:"required,max=160"`
	PriceCents int64  `json:"price_cents" validate:"gte=0"`
}

type updateProductRequest struct {
	Name       *string `json:"name" validate:"omitempty,max=160"`
	PriceCents *int64  `json:"price_cents" validate:"omitempty,gte=0"`
	Active     *bool   `json:"active"`
}

func (r updateProductRequest) hasChanges() bool {
	return r.Name != nil || r.PriceCents != nil || r.Active != nil
}

// Create handles product creation requests.
func (h *Handler) Create(c *echo.Context) error {
	if err := h.ready(); err != nil {
		return err
	}
	var req createProductRequest
	if err := httpreq.BindAndValidate(c, &req); err != nil {
		return err
	}
	product, err := h.usecase.Create(c.Request().Context(), usecase.CreateInput{
		SKU:        req.SKU,
		Name:       req.Name,
		PriceCents: req.PriceCents,
	})
	if err != nil {
		return err
	}
	return httpresp.Created(c, product)
}

// Get handles product lookup requests.
func (h *Handler) Get(c *echo.Context) error {
	if err := h.ready(); err != nil {
		return err
	}
	id, err := httpreq.PathID(c, "id", "product")
	if err != nil {
		return err
	}
	product, err := h.usecase.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return httpresp.OK(c, product)
}

// List handles product list requests.
func (h *Handler) List(c *echo.Context) error {
	if err := h.ready(); err != nil {
		return err
	}
	page, pageSize, err := httpreq.Pagination(c, defaultPageSize)
	if err != nil {
		return err
	}
	activeOnly, err := httpreq.QueryBool(c, "active_only", false)
	if err != nil {
		return err
	}
	output, err := h.usecase.List(c.Request().Context(), usecase.ListInput{
		Page:       page,
		PageSize:   pageSize,
		Query:      c.QueryParam("q"),
		ActiveOnly: activeOnly,
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

// Update handles product partial update requests.
func (h *Handler) Update(c *echo.Context) error {
	if err := h.ready(); err != nil {
		return err
	}
	id, err := httpreq.PathID(c, "id", "product")
	if err != nil {
		return err
	}
	var req updateProductRequest
	if bindErr := httpreq.BindAndValidate(c, &req); bindErr != nil {
		return bindErr
	}
	if !req.hasChanges() {
		return apperr.NewBadRequest("empty product update")
	}
	product, err := h.usecase.Update(c.Request().Context(), usecase.UpdateInput{
		ID:         id,
		Name:       req.Name,
		PriceCents: req.PriceCents,
		Active:     req.Active,
	})
	if err != nil {
		return err
	}
	return httpresp.OK(c, product)
}

func (h *Handler) ready() error {
	if h == nil || h.usecase == nil {
		return apperr.New(apperr.ErrInternalServer, "product usecase is not configured")
	}
	return nil
}
