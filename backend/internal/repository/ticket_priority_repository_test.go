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

func setupPriorityTestDB_SQLite(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.TicketPriority{})
	require.NoError(t, err)
	priority := &domain.TicketPriority{ID: 1, Code: "high", Label: "High", Level: 1}
	require.NoError(t, db.Create(priority).Error)
	return db
}

func TestTicketPriorityRepository_SQLite(t *testing.T) {
	db := setupPriorityTestDB_SQLite(t)
	repo := NewTicketPriorityRepository(db)

	priorities, err := repo.List()
	assert.NoError(t, err)
	assert.NotEmpty(t, priorities)
}

func setupTicketPriorityTestDB_Postgres(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5434 user=postgres password=password dbname=ticket_system_test sslmode=disable"
	if !strings.Contains("ticket_system_test", "_test") {
		t.Fatal("Refusing to run tests on non-test database!")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	db.Exec("TRUNCATE ticket_priorities RESTART IDENTITY CASCADE")
	priority := &domain.TicketPriority{ID: 1, Code: "high", Label: "High", Level: 1}
	require.NoError(t, db.Create(priority).Error)
	return db
}

func TestTicketPriorityRepository_Postgres(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("set INTEGRATION_TEST=1 to run")
	}
	db := setupTicketPriorityTestDB_Postgres(t)
	repo := NewTicketPriorityRepository(db)

	priorities, err := repo.List()
	assert.NoError(t, err)
	assert.NotEmpty(t, priorities)
}
