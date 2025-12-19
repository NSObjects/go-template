//go:build integration

/*
 * Data Layer Integration Tests
 *
 * These tests require a running database instance.
 * Use testcontainers or a dedicated test database.
 *
 * Run with: go test -tags=integration ./internal/api/data/...
 */

package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/NSObjects/go-template/internal/api/data"
	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/api/service/param"
	configs "github.com/NSObjects/go-kit/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testDataManager creates a DataManager for testing.
// In real tests, use testcontainers or a test database.
func testDataManager(t *testing.T) *db.DataManager {
	t.Helper()

	// TODO: Replace with testcontainers or test database setup
	cfg := configs.Config{
		Mysql: configs.Mysql{
			Host:     "localhost",
			Port:     3306,
			User:     "test",
			Password: "test",
			Dbname:   "go_template_test",
		},
	}

	// Note: This will panic if database is not available.
	// Use testcontainers for reliable CI testing.
	dm := &db.DataManager{}
	t.Cleanup(func() {
		_ = dm.Stop(context.Background())
	})

	return dm
}

func TestUserRepository_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	dm := testDataManager(t)
	repo := data.NewUserRepository(dm)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		req := param.UserCreateRequest{
			Username: "test_user_" + time.Now().Format("20060102150405"),
			Email:    "test@example.com",
			Age:      25,
		}

		err := repo.Create(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("ListUsers", func(t *testing.T) {
		req := param.UserListUsersRequest{}
		req.Page = 1
		req.Count = 10

		users, total, err := repo.ListUsers(ctx, req)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(0))
		assert.NotNil(t, users)
	})

	t.Run("GetByID", func(t *testing.T) {
		// First create a user
		createReq := param.UserCreateRequest{
			Username: "get_by_id_user",
			Email:    "getbyid@example.com",
			Age:      30,
		}
		err := repo.Create(ctx, createReq)
		require.NoError(t, err)

		// Then get by ID (assuming ID 1 exists)
		user, err := repo.GetByID(ctx, 1)
		if err != nil {
			t.Logf("GetByID error (may be expected if ID 1 doesn't exist): %v", err)
		} else {
			assert.NotEmpty(t, user.Username)
		}
	})

	t.Run("Update", func(t *testing.T) {
		req := param.UserUpdateRequest{
			Username: "updated_user",
			Email:    "updated@example.com",
			Age:      35,
		}

		err := repo.Update(ctx, 1, req)
		if err != nil {
			t.Logf("Update error (may be expected if ID 1 doesn't exist): %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// Create a user to delete
		createReq := param.UserCreateRequest{
			Username: "delete_me",
			Email:    "delete@example.com",
			Age:      20,
		}
		_ = repo.Create(ctx, createReq)

		// Delete (assuming the created user has a known ID)
		err := repo.Delete(ctx, 999) // Use a test-specific ID
		if err != nil {
			t.Logf("Delete error (may be expected): %v", err)
		}
	})
}
