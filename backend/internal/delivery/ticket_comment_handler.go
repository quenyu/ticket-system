package delivery

import (
	"net/http"
	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/model"
	"ticket-system/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketCommentHandler struct {
	Repo usecase.TicketCommentRepository
}

func NewTicketCommentHandler(repo usecase.TicketCommentRepository) *TicketCommentHandler {
	return &TicketCommentHandler{Repo: repo}
}

func (h *TicketCommentHandler) GetComments(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.APIErrorDetail{Code: 400, Message: "Invalid ticket id"},
		})
		return
	}
	comments, err := h.Repo.GetByTicketID(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.APIErrorDetail{Code: 500, Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (h *TicketCommentHandler) AddComment(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.APIErrorDetail{Code: 400, Message: "Invalid ticket id"},
		})
		return
	}
	var comment domain.TicketComment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.APIErrorDetail{Code: 400, Message: "Validation failed", Details: err.Error()},
		})
		return
	}
	comment.TicketID = ticketID
	if err := h.Repo.Create(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.APIErrorDetail{Code: 500, Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusCreated, comment)
}
