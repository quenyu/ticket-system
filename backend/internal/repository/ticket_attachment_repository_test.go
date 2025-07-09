package repository

import (
	"os"
	"testing"
	"time"

	"ticket-system/backend/internal/domain"

	"strings"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAttachmentTestDB_SQLite(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.Role{}, &domain.Department{}, &domain.User{}, &domain.Ticket{}, &domain.TicketAttachment{})
	require.NoError(t, err)
	role := &domain.Role{ID: uuid.New(), Name: "user"}
	require.NoError(t, db.Create(role).Error)
	dept := &domain.Department{ID: 1, Name: "IT"}
	require.NoError(t, db.Create(dept).Error)
	user := &domain.User{ID: uuid.New(), Username: "user", PasswordHash: "hash", Email: "u@e.com", RoleID: role.ID, DepartmentID: dept.ID, CreatedAt: time.Now()}
	require.NoError(t, db.Create(user).Error)
	assignee := &domain.User{ID: uuid.New(), Username: "assignee", PasswordHash: "hash2", Email: "a@e.com", RoleID: role.ID, DepartmentID: dept.ID, CreatedAt: time.Now()}
	require.NoError(t, db.Create(assignee).Error)
	ticket := &domain.Ticket{ID: uuid.New(), Title: "T", StatusID: 1, PriorityID: 1, CreatorID: user.ID, AssigneeID: assignee.ID, DepartmentID: dept.ID, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	require.NoError(t, db.Create(ticket).Error)
	return db
}

func TestTicketAttachmentRepository_SQLite(t *testing.T) {
	db := setupAttachmentTestDB_SQLite(t)
	repo := NewTicketAttachmentRepository(db)
	var user domain.User
	require.NoError(t, db.First(&user).Error)
	var ticket domain.Ticket
	require.NoError(t, db.First(&ticket).Error)
	path := "/tmp/test.txt"
	attachment := &domain.TicketAttachment{
		ID:         uuid.New(),
		TicketID:   ticket.ID,
		Filename:   "test.txt",
		FilePath:   &path,
		UploadedBy: user.ID,
		UploadedAt: time.Now(),
	}
	err := repo.Create(attachment)
	assert.NoError(t, err)
}

func setupAttachmentTestDB_Postgres(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5434 user=postgres password=password dbname=ticket_system_test sslmode=disable"
	if !strings.Contains("ticket_system_test", "_test") {
		t.Fatal("Refusing to run tests on non-test database!")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	db.Exec("TRUNCATE ticket_attachments, tickets, users, roles, departments RESTART IDENTITY CASCADE")
	role := &domain.Role{ID: uuid.New(), Name: "user"}
	require.NoError(t, db.Create(role).Error)
	dept := &domain.Department{ID: 1, Name: "IT"}
	require.NoError(t, db.Create(dept).Error)
	user := &domain.User{ID: uuid.New(), Username: "user", PasswordHash: "hash", Email: "u@e.com", RoleID: role.ID, DepartmentID: dept.ID, CreatedAt: time.Now()}
	require.NoError(t, db.Create(user).Error)
	assignee := &domain.User{ID: uuid.New(), Username: "assignee", PasswordHash: "hash2", Email: "a@e.com", RoleID: role.ID, DepartmentID: dept.ID, CreatedAt: time.Now()}
	require.NoError(t, db.Create(assignee).Error)
	ticket := &domain.Ticket{ID: uuid.New(), Title: "T", StatusID: 1, PriorityID: 1, CreatorID: user.ID, AssigneeID: assignee.ID, DepartmentID: dept.ID, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	require.NoError(t, db.Create(ticket).Error)
	return db
}

func TestTicketAttachmentRepository_Postgres(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("set INTEGRATION_TEST=1 to run")
	}
	db := setupAttachmentTestDB_Postgres(t)
	repo := NewTicketAttachmentRepository(db)
	var user domain.User
	require.NoError(t, db.First(&user).Error)
	var ticket domain.Ticket
	require.NoError(t, db.First(&ticket).Error)
	path := "/tmp/test.txt"
	attachment := &domain.TicketAttachment{
		ID:         uuid.New(),
		TicketID:   ticket.ID,
		Filename:   "test.txt",
		FilePath:   &path,
		UploadedBy: user.ID,
		UploadedAt: time.Now(),
	}
	err := repo.Create(attachment)
	assert.NoError(t, err)
}
