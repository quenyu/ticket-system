package delivery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestUserRepo(t *testing.T) *repository.UserRepository {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&domain.User{}))
	user := &domain.User{ID: uuid.New(), Username: "alice", DepartmentID: 1}
	repo := repository.NewUserRepository(db)
	require.NoError(t, repo.Create(user))
	return repo
}

func TestUserHandler_ListUsers_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := setupTestUserRepo(t)
	h := &UserHandler{Repo: repo}
	r := gin.New()
	r.GET("/users", h.ListUsers)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "alice")
}

func TestUserHandler_ListUsers_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	repo := repository.NewUserRepository(db)
	h := &UserHandler{Repo: repo}
	r := gin.New()
	r.GET("/users", h.ListUsers)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "failed to fetch users")
}
