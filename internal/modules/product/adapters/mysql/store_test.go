package productmysql

import (
	"testing"
	"time"

	"github.com/NSObjects/go-template/internal/modules/product/domain"
)

func TestRecordDomainRoundTripPreservesProductFields(t *testing.T) {
	createdAt := time.Date(2026, 6, 19, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	product, err := domain.Restore(42, "SKU-1", "Starter", 1999, true, createdAt, updatedAt)
	if err != nil {
		t.Fatalf("Restore() error = %v", err)
	}

	roundTripped, err := domainFromRecord(recordFromDomain(product))
	if err != nil {
		t.Fatalf("domainFromRecord() error = %v", err)
	}

	if roundTripped.ID() != product.ID() {
		t.Fatalf("ID = %d, want %d", roundTripped.ID(), product.ID())
	}
	if roundTripped.SKU() != product.SKU() {
		t.Fatalf("SKU = %q, want %q", roundTripped.SKU(), product.SKU())
	}
	if roundTripped.Name() != product.Name() {
		t.Fatalf("Name = %q, want %q", roundTripped.Name(), product.Name())
	}
	if roundTripped.PriceCents() != product.PriceCents() {
		t.Fatalf("PriceCents = %d, want %d", roundTripped.PriceCents(), product.PriceCents())
	}
	if roundTripped.Active() != product.Active() {
		t.Fatalf("Active = %t, want %t", roundTripped.Active(), product.Active())
	}
	if !roundTripped.CreatedAt().Equal(product.CreatedAt()) {
		t.Fatalf("CreatedAt = %v, want %v", roundTripped.CreatedAt(), product.CreatedAt())
	}
	if !roundTripped.UpdatedAt().Equal(product.UpdatedAt()) {
		t.Fatalf("UpdatedAt = %v, want %v", roundTripped.UpdatedAt(), product.UpdatedAt())
	}
}

func TestDomainFromRecordRejectsInvalidPersistedProduct(t *testing.T) {
	_, err := domainFromRecord(productRecord{
		ID:         1,
		SKU:        "",
		Name:       "Starter",
		PriceCents: 1999,
		Active:     true,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	})
	if err == nil {
		t.Fatal("domainFromRecord() error = nil, want invalid SKU error")
	}
}
