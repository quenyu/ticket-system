package delivery

import (
	"net/http"
	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/model"
	"ticket-system/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketHistoryHandler struct {
	Repo usecase.TicketHistoryRepository
}

func NewTicketHistoryHandler(repo usecase.TicketHistoryRepository) *TicketHistoryHandler {
	return &TicketHistoryHandler{Repo: repo}
}

func (h *TicketHistoryHandler) GetHistory(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "400",
			Message: "Invalid ticket id",
		})
		return
	}
	history, err := h.Repo.GetByTicketID(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "500",
			Message: err.Error(),
		})
		return
	}
	if history == nil {
		history = []*domain.TicketHistory{}
	}
	c.JSON(http.StatusOK, history)
}
