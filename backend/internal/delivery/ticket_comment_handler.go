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
			Code:    "400",
			Message: "Invalid ticket id",
		})
		return
	}
	comments, err := h.Repo.GetByTicketID(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "500",
			Message: err.Error(),
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
			Code:    "400",
			Message: "Invalid ticket id",
		})
		return
	}
	var comment domain.TicketComment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "400",
			Message: "Validation failed",
			Details: err.Error(),
		})
		return
	}
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.APIError{
			Code:    "401",
			Message: "No user_id in token",
		})
		return
	}
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.APIError{
			Code:    "401",
			Message: "user_id in token is not a string",
		})
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.APIError{
			Code:    "401",
			Message: "Invalid user_id in token",
		})
		return
	}
	comment.AuthorID = userID

	comment.TicketID = ticketID
	if comment.ID == uuid.Nil {
		comment.ID = uuid.New()
	}
	if err := h.Repo.Create(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "500",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, comment)
}
