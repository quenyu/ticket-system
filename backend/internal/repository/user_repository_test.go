package repository

import (
	"os"
	"strings"
	"testing"
	"time"

	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserTestDB_SQLite(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.Role{}, &domain.Department{}, &domain.User{})
	require.NoError(t, err)
	role := &domain.Role{ID: uuid.New(), Name: "user"}
	require.NoError(t, db.Create(role).Error)
	dept := &domain.Department{ID: 1, Name: "IT"}
	require.NoError(t, db.Create(dept).Error)
	return db
}

func TestUserRepository_SQLite(t *testing.T) {
	db := setupUserTestDB_SQLite(t)
	repo := NewUserRepository(db)

	var role domain.Role
	require.NoError(t, db.First(&role).Error)
	var dept domain.Department
	require.NoError(t, db.First(&dept).Error)

	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: "hash123",
		Email:        "test@example.com",
		RoleID:       role.ID,
		DepartmentID: dept.ID,
		CreatedAt:    time.Now(),
	}

	// Create
	err := repo.Create(user)
	assert.NoError(t, err)

	// FindByUsername
	got, err := repo.FindByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, got.Email)
}

func setupUserTestDB_Postgres(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5434 user=postgres password=password dbname=ticket_system_test sslmode=disable"
	if !strings.Contains("ticket_system_test", "_test") {
		t.Fatal("Refusing to run tests on non-test database!")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	db.Exec("TRUNCATE users, roles, departments RESTART IDENTITY CASCADE")
	role := &domain.Role{ID: uuid.New(), Name: "user"}
	require.NoError(t, db.Create(role).Error)
	dept := &domain.Department{ID: 1, Name: "IT"}
	require.NoError(t, db.Create(dept).Error)
	return db
}

func TestUserRepository_Postgres(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("set INTEGRATION_TEST=1 to run")
	}
	db := setupUserTestDB_Postgres(t)
	repo := NewUserRepository(db)

	var role domain.Role
	require.NoError(t, db.First(&role).Error)
	var dept domain.Department
	require.NoError(t, db.First(&dept).Error)

	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: "hash123",
		Email:        "test@example.com",
		RoleID:       role.ID,
		DepartmentID: dept.ID,
		CreatedAt:    time.Now(),
	}

	// Create
	err := repo.Create(user)
	assert.NoError(t, err)

	// FindByUsername
	got, err := repo.FindByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, got.Email)
}
