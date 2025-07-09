package repository

import (
	"os"
	"testing"

	"ticket-system/backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupStatusTestDB_SQLite(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.TicketStatus{})
	require.NoError(t, err)
	status := &domain.TicketStatus{ID: 1, Code: "open", Label: "Open"}
	require.NoError(t, db.Create(status).Error)
	return db
}

func TestTicketStatusRepository_SQLite(t *testing.T) {
	db := setupStatusTestDB_SQLite(t)
	repo := NewTicketStatusRepository(db)

	statuses, err := repo.List()
	assert.NoError(t, err)
	assert.NotEmpty(t, statuses)
}

func setupStatusTestDB_Postgres(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5434 user=postgres password=password dbname=ticket_system sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	db.Exec("TRUNCATE ticket_statuses RESTART IDENTITY CASCADE")
	status := &domain.TicketStatus{ID: 1, Code: "open", Label: "Open"}
	require.NoError(t, db.Create(status).Error)
	return db
}

func TestTicketStatusRepository_Postgres(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("set INTEGRATION_TEST=1 to run")
	}
	db := setupStatusTestDB_Postgres(t)
	repo := NewTicketStatusRepository(db)

	statuses, err := repo.List()
	assert.NoError(t, err)
	assert.NotEmpty(t, statuses)
}
