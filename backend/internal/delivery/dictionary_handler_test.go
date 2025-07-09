package delivery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ticket-system/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mocks for usecase interfaces (unique names for this file)
type testDepartmentDH struct {
	id   int16
	name string
}

func (d *testDepartmentDH) DepartmentID() int16    { return d.id }
func (d *testDepartmentDH) DepartmentName() string { return d.name }

type mockDepartmentRepoDH struct {
	listFunc func() ([]interface {
		DepartmentID() int16
		DepartmentName() string
	}, error)
}

func (m *mockDepartmentRepoDH) List() ([]interface {
	DepartmentID() int16
	DepartmentName() string
}, error) {
	if m.listFunc != nil {
		return m.listFunc()
	}
	return []interface {
		DepartmentID() int16
		DepartmentName() string
	}{
		&testDepartmentDH{1, "IT"},
	}, nil
}

type testStatusDH struct {
	id          int16
	code, label string
}

func (s *testStatusDH) ID() int16     { return s.id }
func (s *testStatusDH) Code() string  { return s.code }
func (s *testStatusDH) Label() string { return s.label }

type mockStatusRepoDH struct {
	listFunc func() ([]interface {
		ID() int16
		Code() string
		Label() string
	}, error)
}

func (m *mockStatusRepoDH) List() ([]interface {
	ID() int16
	Code() string
	Label() string
}, error) {
	if m.listFunc != nil {
		return m.listFunc()
	}
	return []interface {
		ID() int16
		Code() string
		Label() string
	}{
		&testStatusDH{1, "open", "Open"},
	}, nil
}

type testPriorityDH struct {
	id          int16
	code, label string
	level       int16
}

func (p *testPriorityDH) ID() int16     { return p.id }
func (p *testPriorityDH) Code() string  { return p.code }
func (p *testPriorityDH) Label() string { return p.label }
func (p *testPriorityDH) Level() int16  { return p.level }

type mockPriorityRepoDH struct {
	listFunc func() ([]interface {
		ID() int16
		Code() string
		Label() string
		Level() int16
	}, error)
}

func (m *mockPriorityRepoDH) List() ([]interface {
	ID() int16
	Code() string
	Label() string
	Level() int16
}, error) {
	if m.listFunc != nil {
		return m.listFunc()
	}
	return []interface {
		ID() int16
		Code() string
		Label() string
		Level() int16
	}{
		&testPriorityDH{1, "high", "High", 3},
	}, nil
}

func TestDictionaryHandler_DepartmentsList_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockDepartmentRepoDH{}
	svc := usecase.NewDepartmentService(repo)
	h := &DictionaryHandler{Departments: svc}
	r := gin.New()
	r.GET("/departments", h.DepartmentsList)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "IT")
}

func TestDictionaryHandler_DepartmentsList_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockDepartmentRepoDH{listFunc: func() ([]interface {
		DepartmentID() int16
		DepartmentName() string
	}, error) {
		return nil, errors.New("fail")
	}}
	svc := usecase.NewDepartmentService(repo)
	h := &DictionaryHandler{Departments: svc}
	r := gin.New()
	r.GET("/departments", h.DepartmentsList)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "fail")
}

func TestDictionaryHandler_TicketStatusesList_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockStatusRepoDH{}
	svc := usecase.NewTicketStatusService(repo)
	h := &DictionaryHandler{TicketStatuses: svc}
	r := gin.New()
	r.GET("/ticket_statuses", h.TicketStatusesList)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ticket_statuses", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "open")
}

func TestDictionaryHandler_TicketStatusesList_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockStatusRepoDH{listFunc: func() ([]interface {
		ID() int16
		Code() string
		Label() string
	}, error) {
		return nil, errors.New("fail")
	}}
	svc := usecase.NewTicketStatusService(repo)
	h := &DictionaryHandler{TicketStatuses: svc}
	r := gin.New()
	r.GET("/ticket_statuses", h.TicketStatusesList)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ticket_statuses", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "fail")
}

func TestDictionaryHandler_TicketPrioritiesList_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockPriorityRepoDH{}
	svc := usecase.NewTicketPriorityService(repo)
	h := &DictionaryHandler{TicketPriorities: svc}
	r := gin.New()
	r.GET("/ticket_priorities", h.TicketPrioritiesList)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ticket_priorities", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "high")
}

func TestDictionaryHandler_TicketPrioritiesList_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockPriorityRepoDH{listFunc: func() ([]interface {
		ID() int16
		Code() string
		Label() string
		Level() int16
	}, error) {
		return nil, errors.New("fail")
	}}
	svc := usecase.NewTicketPriorityService(repo)
	h := &DictionaryHandler{TicketPriorities: svc}
	r := gin.New()
	r.GET("/ticket_priorities", h.TicketPrioritiesList)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ticket_priorities", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "fail")
}
