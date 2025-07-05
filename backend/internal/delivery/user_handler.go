package delivery

import (
	"net/http"
	"ticket-system/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.Repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}
	var result []gin.H
	for _, u := range users {
		result = append(result, gin.H{
			"id":       u.ID.String(),
			"username": u.Username,
		})
	}
	c.JSON(http.StatusOK, result)
}
