package customerhttp

import (
	"github.com/labstack/echo/v5"

	"github.com/NSObjects/go-template/internal/modules/customer/usecase"
	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/server/httpreq"
	"github.com/NSObjects/go-template/internal/platform/server/httpresp"
)

const defaultPageSize = 20

// Handler adapts customer HTTP requests to the customer usecase.
type Handler struct {
	usecase *usecase.Usecase
}

// New creates a customer HTTP handler.
func New(uc *usecase.Usecase) *Handler {
	return &Handler{usecase: uc}
}

// Register mounts customer routes on group.
func Register(group *echo.Group, handler *Handler) {
	group.POST("/customers", handler.Create)
	group.GET("/customers", handler.List)
	group.GET("/customers/:id", handler.Get)
	group.PATCH("/customers/:id", handler.Update)
}

type createCustomerRequest struct {
	Name  string `json:"name" validate:"required,max=120"`
	Email string `json:"email" validate:"required,email"`
}

type updateCustomerRequest struct {
	Name   *string `json:"name" validate:"omitempty,max=120"`
	Email  *string `json:"email" validate:"omitempty,email"`
	Status *string `json:"status"`
}

func (r updateCustomerRequest) hasChanges() bool {
	return r.Name != nil || r.Email != nil || r.Status != nil
}

// Create handles customer creation requests.
func (h *Handler) Create(c *echo.Context) error {
	if err := h.ready(); err != nil {
		return err
	}
	var req createCustomerRequest
	if err := httpreq.BindAndValidate(c, &req); err != nil {
		return err
	}
	customer, err := h.usecase.Create(c.Request().Context(), usecase.CreateInput{Name: req.Name, Email: req.Email})
	if err != nil {
		return err
	}
	return httpresp.Created(c, customer)
}

// Get handles customer lookup requests.
func (h *Handler) Get(c *echo.Context) error {
	if err := h.ready(); err != nil {
		return err
	}
	id, err := httpreq.PathID(c, "id", "customer")
	if err != nil {
		return err
	}
	customer, err := h.usecase.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return httpresp.OK(c, customer)
}

// List handles customer list requests.
func (h *Handler) List(c *echo.Context) error {
	if err := h.ready(); err != nil {
		return err
	}
	page, pageSize, err := httpreq.Pagination(c, defaultPageSize)
	if err != nil {
		return err
	}
	output, err := h.usecase.List(c.Request().Context(), usecase.ListInput{
		Page:     page,
		PageSize: pageSize,
		Query:    c.QueryParam("q"),
		Status:   c.QueryParam("status"),
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

// Update handles customer partial update requests.
func (h *Handler) Update(c *echo.Context) error {
	if err := h.ready(); err != nil {
		return err
	}
	id, err := httpreq.PathID(c, "id", "customer")
	if err != nil {
		return err
	}
	var req updateCustomerRequest
	if bindErr := httpreq.BindAndValidate(c, &req); bindErr != nil {
		return bindErr
	}
	if !req.hasChanges() {
		return apperr.NewBadRequest("empty customer update")
	}
	customer, err := h.usecase.Update(c.Request().Context(), usecase.UpdateInput{
		ID:     id,
		Name:   req.Name,
		Email:  req.Email,
		Status: req.Status,
	})
	if err != nil {
		return err
	}
	return httpresp.OK(c, customer)
}

func (h *Handler) ready() error {
	if h == nil || h.usecase == nil {
		return apperr.New(apperr.ErrInternalServer, "customer usecase is not configured")
	}
	return nil
}
