package delivery

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mocks for usecase interfaces (unique names for this file)
type testStatus struct {
	id          int16
	code, label string
}

func (s *testStatus) ID() int16     { return s.id }
func (s *testStatus) Code() string  { return s.code }
func (s *testStatus) Label() string { return s.label }

type mockStatusRepoTH struct{}

func (m *mockStatusRepoTH) List() ([]interface {
	ID() int16
	Code() string
	Label() string
}, error) {
	return []interface {
		ID() int16
		Code() string
		Label() string
	}{
		&testStatus{1, "open", "Open"},
	}, nil
}

type testPriority struct {
	id          int16
	code, label string
	level       int16
}

func (p *testPriority) ID() int16     { return p.id }
func (p *testPriority) Code() string  { return p.code }
func (p *testPriority) Label() string { return p.label }
func (p *testPriority) Level() int16  { return p.level }

type mockPriorityRepoTH struct{}

func (m *mockPriorityRepoTH) List() ([]interface {
	ID() int16
	Code() string
	Label() string
	Level() int16
}, error) {
	return []interface {
		ID() int16
		Code() string
		Label() string
		Level() int16
	}{
		&testPriority{1, "high", "High", 3},
	}, nil
}

type testDepartment struct {
	id   int16
	name string
}

func (d *testDepartment) DepartmentID() int16    { return d.id }
func (d *testDepartment) DepartmentName() string { return d.name }

type mockDepartmentRepoTH struct{}

func (m *mockDepartmentRepoTH) List() ([]interface {
	DepartmentID() int16
	DepartmentName() string
}, error) {
	return []interface {
		DepartmentID() int16
		DepartmentName() string
	}{
		&testDepartment{1, "IT"},
	}, nil
}

func setupTestTicketHandler(t *testing.T) *TicketHandler {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&domain.Ticket{}, &domain.User{}, &domain.TicketStatus{}, &domain.TicketPriority{}, &domain.Department{}, &domain.TicketHistory{}))

	user := &domain.User{ID: uuid.New(), Username: "alice", DepartmentID: 1}
	require.NoError(t, db.Create(user).Error)

	ticket := &domain.Ticket{
		ID:           uuid.New(),
		Title:        "Test Ticket",
		Description:  "Test Desc",
		StatusID:     1,
		PriorityID:   1,
		DepartmentID: 1,
		AssigneeID:   user.ID,
		CreatedAt:    time.Now(),
	}
	require.NoError(t, db.Create(ticket).Error)

	ticketRepo := repository.NewTicketRepository(db)
	statusRepo := &mockStatusRepoTH{}
	prioRepo := &mockPriorityRepoTH{}
	deptRepo := &mockDepartmentRepoTH{}
	userRepo := repository.NewUserRepository(db)
	historyRepo := repository.NewTicketHistoryRepository(db)

	ticketHandler := NewTicketHandler(
		ticketRepo,
		statusRepo,
		prioRepo,
		deptRepo,
		userRepo,
		historyRepo,
	)

	return ticketHandler
}

func TestTicketHandler_GetTickets_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := setupTestTicketHandler(t)
	r := gin.New()
	r.GET("/tickets", h.GetTickets)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tickets", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Test Ticket")
	assert.Contains(t, w.Body.String(), "alice")
}

func TestTicketHandler_GetTickets_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	ticketRepo := repository.NewTicketRepository(db)
	statusRepo := &mockStatusRepoTH{}
	prioRepo := &mockPriorityRepoTH{}
	deptRepo := &mockDepartmentRepoTH{}
	userRepo := repository.NewUserRepository(db)
	h := NewTicketHandler(ticketRepo, statusRepo, prioRepo, deptRepo, userRepo, nil)
	r := gin.New()
	r.GET("/tickets", h.GetTickets)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tickets", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get tickets")
}

func TestTicketHandler_CreateTicket_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := setupTestTicketHandler(t)
	r := gin.New()
	r.POST("/tickets", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		h.CreateTicket(c)
	})
	body := `{"title":"New Ticket","description":"Desc","status_id":1,"priority_id":1,"department_id":1,"assignee_id":"` + uuid.New().String() + `"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tickets", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	assert.Contains(t, w.Body.String(), "New Ticket")
}

func TestTicketHandler_CreateTicket_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := setupTestTicketHandler(t)
	r := gin.New()
	r.POST("/tickets", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		h.CreateTicket(c)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tickets", strings.NewReader("not a json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestTicketHandler_CreateTicket_NoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := setupTestTicketHandler(t)
	r := gin.New()
	r.POST("/tickets", h.CreateTicket)
	body := `{"title":"New Ticket","description":"Desc","status_id":1,"priority_id":1,"department_id":1,"assignee_id":"` + uuid.New().String() + `"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tickets", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "User not authenticated")
}

func TestTicketHandler_CreateTicket_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := setupTestTicketHandler(t)
	r := gin.New()
	r.POST("/tickets", func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		h.CreateTicket(c)
	})
	body := `{"title":"","description":"","status_id":0,"priority_id":0,"department_id":0}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tickets", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "Title is required")
}
