package delivery

import (
	"net/http"

	"ticket-system/backend/internal/model"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService AuthService
}

type AuthService interface {
	Login(username, password string) (model.LoginResponse, error)
	Register(req model.UserRegisterRequest) (model.UserRegisterResponse, error)
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.APIErrorDetail{
				Code:    400,
				Message: "Validation failed",
				Details: err.Error(),
			},
		})
		return
	}
	resp, err := h.AuthService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.APIError{
			Error: model.APIErrorDetail{
				Code:    401,
				Message: "Invalid credentials",
			},
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.APIErrorDetail{
				Code:    400,
				Message: "Validation failed",
				Details: err.Error(),
			},
		})
		return
	}
	resp, err := h.AuthService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.APIErrorDetail{
				Code:    400,
				Message: err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusCreated, resp)
}
