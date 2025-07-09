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

func setupCommentTestDB_SQLite(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.Role{}, &domain.Department{}, &domain.User{}, &domain.Ticket{}, &domain.TicketComment{})
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

func TestTicketCommentRepository_CRUD_SQLite(t *testing.T) {
	db := setupCommentTestDB_SQLite(t)
	repo := NewTicketCommentRepository(db)
	var user domain.User
	require.NoError(t, db.First(&user).Error)
	var ticket domain.Ticket
	require.NoError(t, db.First(&ticket).Error)
	comment := &domain.TicketComment{
		ID:        uuid.New(),
		TicketID:  ticket.ID,
		AuthorID:  user.ID,
		Content:   "Test comment",
		CreatedAt: time.Now(),
	}
	err := repo.Create(comment)
	assert.NoError(t, err)
	got, err := repo.GetByID(comment.ID)
	assert.NoError(t, err)
	assert.Equal(t, comment.Content, got.Content)
	err = repo.Delete(comment.ID)
	assert.NoError(t, err)
	_, err = repo.GetByID(comment.ID)
	assert.Error(t, err)
}

func setupCommentTestDB_Postgres(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5434 user=postgres password=password dbname=ticket_system sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	db.Exec("TRUNCATE ticket_comments, tickets, users, roles, departments RESTART IDENTITY CASCADE")
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

func TestTicketCommentRepository_CRUD_Postgres(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("set INTEGRATION_TEST=1 to run")
	}
	db := setupCommentTestDB_Postgres(t)
	repo := NewTicketCommentRepository(db)
	var user domain.User
	require.NoError(t, db.First(&user).Error)
	var ticket domain.Ticket
	require.NoError(t, db.First(&ticket).Error)
	comment := &domain.TicketComment{
		ID:        uuid.New(),
		TicketID:  ticket.ID,
		AuthorID:  user.ID,
		Content:   "Integration comment",
		CreatedAt: time.Now(),
	}
	err := repo.Create(comment)
	assert.NoError(t, err)
	got, err := repo.GetByID(comment.ID)
	assert.NoError(t, err)
	assert.Equal(t, comment.Content, got.Content)
	err = repo.Delete(comment.ID)
	assert.NoError(t, err)
	_, err = repo.GetByID(comment.ID)
	assert.Error(t, err)
}
