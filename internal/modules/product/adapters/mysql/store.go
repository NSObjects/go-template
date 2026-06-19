package productmysql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/NSObjects/go-template/internal/modules/product/domain"
	"github.com/NSObjects/go-template/internal/modules/product/usecase"
	"github.com/NSObjects/go-template/internal/platform/apperr"
)

const productTableName = "products"

var _ usecase.Store = (*Store)(nil)

// Store persists products in MySQL through the boot-owned GORM connection.
type Store struct {
	db *gorm.DB
}

type productRecord struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"`
	SKU        string `gorm:"size:64;not null;uniqueIndex"`
	Name       string `gorm:"size:160;not null"`
	PriceCents int64  `gorm:"not null"`
	Active     bool   `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (productRecord) TableName() string {
	return productTableName
}

// NewStore creates a MySQL-backed product store and prepares its table.
func NewStore(db *gorm.DB) (*Store, error) {
	if db == nil {
		return nil, errors.New("product mysql store requires db")
	}
	if err := db.AutoMigrate(&productRecord{}); err != nil {
		return nil, fmt.Errorf("migrate product table: %w", err)
	}
	return &Store{db: db}, nil
}

// Create persists a product and returns the assigned id and timestamps.
func (s *Store) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	exists, err := s.skuExists(ctx, product.SKU())
	if err != nil {
		return domain.Product{}, err
	}
	if exists {
		return domain.Product{}, apperr.NewConflict("product sku already exists")
	}

	record := recordFromDomain(product)
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.Product{}, err
	}
	return domainFromRecord(record)
}

// FindByID returns one persisted product.
func (s *Store) FindByID(ctx context.Context, id int64) (domain.Product, error) {
	var record productRecord
	if err := s.db.WithContext(ctx).First(&record, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Product{}, apperr.NewNotFound("product")
		}
		return domain.Product{}, err
	}
	return domainFromRecord(record)
}

// List returns products matching the validated filter.
func (s *Store) List(ctx context.Context, filter usecase.ListFilter) ([]domain.Product, int, error) {
	query := applyListFilter(s.db.WithContext(ctx).Model(&productRecord{}), filter)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []productRecord
	if err := query.
		Order("id DESC").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&records).Error; err != nil {
		return nil, int(total), err
	}

	products := make([]domain.Product, 0, len(records))
	for _, record := range records {
		product, err := domainFromRecord(record)
		if err != nil {
			return nil, int(total), err
		}
		products = append(products, product)
	}
	return products, int(total), nil
}

// Update replaces mutable product fields while preserving SKU and creation time.
func (s *Store) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	var record productRecord
	if err := s.db.WithContext(ctx).First(&record, product.ID()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Product{}, apperr.NewNotFound("product")
		}
		return domain.Product{}, err
	}

	record.Name = product.Name()
	record.PriceCents = product.PriceCents()
	record.Active = product.Active()
	if err := s.db.WithContext(ctx).Save(&record).Error; err != nil {
		return domain.Product{}, err
	}
	return domainFromRecord(record)
}

func (s *Store) skuExists(ctx context.Context, sku string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&productRecord{}).
		Where("sku = ?", sku).
		Count(&count).Error
	return count > 0, err
}

func applyListFilter(query *gorm.DB, filter usecase.ListFilter) *gorm.DB {
	if filter.ActiveOnly {
		query = query.Where("active = ?", true)
	}
	if queryText := strings.TrimSpace(filter.Query); queryText != "" {
		like := "%" + strings.ToLower(queryText) + "%"
		query = query.Where("LOWER(sku) LIKE ? OR LOWER(name) LIKE ?", like, like)
	}
	return query
}

func recordFromDomain(product domain.Product) productRecord {
	return productRecord{
		ID:         product.ID(),
		SKU:        product.SKU(),
		Name:       product.Name(),
		PriceCents: product.PriceCents(),
		Active:     product.Active(),
		CreatedAt:  product.CreatedAt(),
		UpdatedAt:  product.UpdatedAt(),
	}
}

func domainFromRecord(record productRecord) (domain.Product, error) {
	return domain.Restore(
		record.ID,
		record.SKU,
		record.Name,
		record.PriceCents,
		record.Active,
		record.CreatedAt,
		record.UpdatedAt,
	)
}
