package repository

import (
	"os"
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

func setupTicketTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	// Миграция всех нужных структур
	err = db.AutoMigrate(
		&domain.Role{},
		&domain.Department{},
		&domain.User{},
		&domain.TicketStatus{},
		&domain.TicketPriority{},
		&domain.Ticket{},
	)
	require.NoError(t, err)
	role := &domain.Role{ID: uuid.New(), Name: "user"}
	require.NoError(t, db.Create(role).Error)
	dept := &domain.Department{ID: 1, Name: "IT"}
	require.NoError(t, db.Create(dept).Error)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "user",
		PasswordHash: "hash",
		Email:        "u@e.com",
		RoleID:       role.ID,
		DepartmentID: dept.ID,
		CreatedAt:    time.Now(),
	}
	require.NoError(t, db.Create(user).Error)
	status := &domain.TicketStatus{ID: 1, Code: "open", Label: "Open"}
	priority := &domain.TicketPriority{ID: 1, Code: "high", Label: "High", Level: 1}
	require.NoError(t, db.Create(status).Error)
	require.NoError(t, db.Create(priority).Error)
	return db
}

func TestTicketRepository_CRUD(t *testing.T) {
	db := setupTicketTestDB(t)
	repo := NewTicketRepository(db)

	var user domain.User
	require.NoError(t, db.First(&user).Error)

	ticket := &domain.Ticket{
		ID:           uuid.New(),
		Title:        "Test Ticket",
		Description:  "Test Description",
		StatusID:     1,
		PriorityID:   1,
		CreatorID:    user.ID,
		AssigneeID:   user.ID, // исправлено
		DepartmentID: 1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create
	err := repo.Create(ticket)
	assert.NoError(t, err)

	// GetByID
	got, err := repo.GetByID(ticket.ID)
	assert.NoError(t, err)
	assert.Equal(t, ticket.Title, got.Title)

	// Update
	got.Title = "Updated Title"
	err = repo.Update(got)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ticket.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)

	// Delete
	err = repo.Delete(ticket.ID)
	assert.NoError(t, err)
	_, err = repo.GetByID(ticket.ID)
	assert.Error(t, err)
}

func TestTicketRepository_CRUD_SQLite(t *testing.T) {
	db := setupTicketTestDB(t)
	repo := NewTicketRepository(db)

	var user domain.User
	require.NoError(t, db.First(&user).Error)

	ticket := &domain.Ticket{
		ID:           uuid.New(),
		Title:        "Test Ticket",
		Description:  "Test Description",
		StatusID:     1,
		PriorityID:   1,
		CreatorID:    user.ID,
		AssigneeID:   user.ID, // исправлено
		DepartmentID: 1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create
	err := repo.Create(ticket)
	assert.NoError(t, err)

	// GetByID
	got, err := repo.GetByID(ticket.ID)
	assert.NoError(t, err)
	assert.Equal(t, ticket.Title, got.Title)

	// Update
	got.Title = "Updated Title"
	err = repo.Update(got)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ticket.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)

	// Delete
	err = repo.Delete(ticket.ID)
	assert.NoError(t, err)
	_, err = repo.GetByID(ticket.ID)
	assert.Error(t, err)
}

func setupPostgresTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5434 user=postgres password=password dbname=ticket_system sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	db.Exec("TRUNCATE tickets, users, ticket_statuses, ticket_priorities, departments, roles RESTART IDENTITY CASCADE")
	role := &domain.Role{ID: uuid.New(), Name: "user"}
	require.NoError(t, db.Create(role).Error)
	dept := &domain.Department{ID: 1, Name: "IT"}
	require.NoError(t, db.Create(dept).Error)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "user",
		PasswordHash: "hash",
		Email:        "u@e.com",
		RoleID:       role.ID,
		DepartmentID: dept.ID,
		CreatedAt:    time.Now(),
	}
	require.NoError(t, db.Create(user).Error)
	status := &domain.TicketStatus{ID: 1, Code: "open", Label: "Open"}
	priority := &domain.TicketPriority{ID: 1, Code: "high", Label: "High", Level: 1}
	require.NoError(t, db.Create(status).Error)
	require.NoError(t, db.Create(priority).Error)
	return db
}

func TestTicketRepository_CRUD_Postgres(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("set INTEGRATION_TEST=1 to run")
	}
	db := setupPostgresTestDB(t)
	repo := NewTicketRepository(db)

	var user domain.User
	require.NoError(t, db.First(&user).Error)

	ticket := &domain.Ticket{
		ID:           uuid.New(),
		Title:        "Integration Ticket",
		Description:  "Integration Description",
		StatusID:     1,
		PriorityID:   1,
		CreatorID:    user.ID,
		AssigneeID:   user.ID, // исправлено
		DepartmentID: 1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create
	err := repo.Create(ticket)
	assert.NoError(t, err)

	// GetByID
	got, err := repo.GetByID(ticket.ID)
	assert.NoError(t, err)
	assert.Equal(t, ticket.Title, got.Title)

	// Update
	got.Title = "Updated Integration Title"
	err = repo.Update(got)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ticket.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Integration Title", updated.Title)

	// Delete
	err = repo.Delete(ticket.ID)
	assert.NoError(t, err)
	_, err = repo.GetByID(ticket.ID)
	assert.Error(t, err)
}
