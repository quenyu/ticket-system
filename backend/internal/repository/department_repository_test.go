package repository

import (
	"os"
	"strings"
	"testing"

	"ticket-system/backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDepartmentTestDB_SQLite(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.Department{})
	require.NoError(t, err)
	dept := &domain.Department{ID: 1, Name: "IT"}
	require.NoError(t, db.Create(dept).Error)
	return db
}

func TestDepartmentRepository_SQLite(t *testing.T) {
	db := setupDepartmentTestDB_SQLite(t)
	repo := NewDepartmentRepository(db)

	depts, err := repo.List()
	assert.NoError(t, err)
	assert.NotEmpty(t, depts)
}

func setupDepartmentTestDB_Postgres(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5434 user=postgres password=password dbname=ticket_system_test sslmode=disable"
	if !strings.Contains("ticket_system_test", "_test") {
		t.Fatal("Refusing to run tests on non-test database!")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	db.Exec("TRUNCATE departments RESTART IDENTITY CASCADE")
	dept := &domain.Department{ID: 1, Name: "IT"}
	require.NoError(t, db.Create(dept).Error)
	return db
}

func TestDepartmentRepository_Postgres(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("set INTEGRATION_TEST=1 to run")
	}
	db := setupDepartmentTestDB_Postgres(t)
	repo := NewDepartmentRepository(db)

	depts, err := repo.List()
	assert.NoError(t, err)
	assert.NotEmpty(t, depts)
}
