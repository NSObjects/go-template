package mysql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/data"
	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/code"
	"github.com/NSObjects/go-template/internal/configs"
	user "github.com/NSObjects/go-template/internal/modules/user"
	"github.com/NSObjects/go-template/internal/platform/module"
)

const (
	ModuleName   = "user-storage-mysql"
	ProviderName = "mysql"
)

var ErrInvalidConfig = errors.New("invalid user mysql storage config")

// Module exposes MySQL as an optional user.storage provider.
type Module struct {
	cfg configs.MysqlConfig
}

// New creates a MySQL user storage provider module from configuration.
func New(cfg configs.MysqlConfig) Module {
	return Module{cfg: cfg}
}

// Descriptor reports whether MySQL can satisfy user.storage.
func (m Module) Descriptor() module.Descriptor {
	status := module.CapabilityDisabled
	if m.cfg.Enabled {
		status = module.CapabilityEnabled
		if !hasRequiredConfig(m.cfg) {
			status = module.CapabilityUnavailable
		}
	}
	return module.Descriptor{
		Name: ModuleName,
		Kind: module.CapabilityModule,
		Provides: []module.Capability{
			{
				Name:     user.StorageCapability,
				Provider: ProviderName,
				Status:   status,
				Value:    m.Repository(),
			},
		},
	}
}

// Repository returns a lazy MySQL-backed user repository.
func (m Module) Repository() biz.UserRepository {
	return &repository{cfg: m.cfg}
}

// Validate returns a startup-blocking error when selected MySQL storage config is unusable.
func (m Module) Validate() error {
	if !m.cfg.Enabled {
		return nil
	}
	if !hasRequiredConfig(m.cfg) {
		return fmt.Errorf("%w: host, port, user, and database are required", ErrInvalidConfig)
	}
	return nil
}

// Start verifies the selected MySQL storage provider can initialize and reach MySQL.
func (m Module) Start(ctx context.Context) error {
	if !m.cfg.Enabled {
		return nil
	}
	if err := m.Validate(); err != nil {
		return err
	}
	dataManager, err := db.NewDataManager(configs.Config{Mysql: m.cfg})
	if err != nil {
		return code.WrapDatabaseError(err, "initialize user mysql storage")
	}
	defer dataManager.Shutdown(context.Background())
	if err := dataManager.Start(ctx); err != nil {
		return code.WrapDatabaseError(err, "start user mysql storage")
	}
	return nil
}

type repository struct {
	cfg  configs.MysqlConfig
	once sync.Once
	repo biz.UserRepository
	err  error
}

func (r *repository) ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {
	repo, err := r.load()
	if err != nil {
		return nil, 0, err
	}
	return repo.ListUsers(ctx, req)
}

func (r *repository) Create(ctx context.Context, req param.UserCreateRequest) error {
	repo, err := r.load()
	if err != nil {
		return err
	}
	return repo.Create(ctx, req)
}

func (r *repository) GetByID(ctx context.Context, id int64) (param.UserData, error) {
	repo, err := r.load()
	if err != nil {
		return param.UserData{}, err
	}
	return repo.GetByID(ctx, id)
}

func (r *repository) Update(ctx context.Context, id int64, req param.UserUpdateRequest) error {
	repo, err := r.load()
	if err != nil {
		return err
	}
	return repo.Update(ctx, id, req)
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	repo, err := r.load()
	if err != nil {
		return err
	}
	return repo.Delete(ctx, id)
}

func (r *repository) load() (biz.UserRepository, error) {
	r.once.Do(func() {
		dataManager, err := db.NewDataManager(configs.Config{Mysql: r.cfg})
		if err != nil {
			r.err = code.WrapDatabaseError(err, "initialize user mysql repository")
			return
		}
		r.repo = data.NewUserRepository(dataManager)
	})
	return r.repo, r.err
}

func hasRequiredConfig(cfg configs.MysqlConfig) bool {
	return !blank(cfg.Host) && !blank(cfg.Port) && !blank(cfg.User) && !blank(cfg.Database)
}

func blank(value string) bool {
	return strings.TrimSpace(value) == ""
}
