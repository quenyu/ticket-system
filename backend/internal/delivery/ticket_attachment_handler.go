package delivery

import (
	"net/http"
	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/model"
	"ticket-system/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketAttachmentHandler struct {
	Repo usecase.TicketAttachmentRepository
}

func NewTicketAttachmentHandler(repo usecase.TicketAttachmentRepository) *TicketAttachmentHandler {
	return &TicketAttachmentHandler{Repo: repo}
}

func (h *TicketAttachmentHandler) GetAttachments(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "400",
			Message: "Invalid ticket id",
		})
		return
	}
	atts, err := h.Repo.GetByTicketID(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "500",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, atts)
}

func (h *TicketAttachmentHandler) AddAttachment(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "400",
			Message: "Invalid ticket id",
		})
		return
	}
	var att domain.TicketAttachment
	if err := c.ShouldBindJSON(&att); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "400",
			Message: "Validation failed",
			Details: err.Error(),
		})
		return
	}
	att.TicketID = ticketID
	if err := h.Repo.Create(&att); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "500",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, att)
}

func (h *TicketAttachmentHandler) GetAttachmentByID(c *gin.Context) {
	attIDStr := c.Param("att_id")
	attID, err := uuid.Parse(attIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "400",
			Message: "Invalid attachment id",
		})
		return
	}
	att, err := h.Repo.GetByID(attID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:    "404",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, att)
}

func (h *TicketAttachmentHandler) DeleteAttachment(c *gin.Context) {
	attIDStr := c.Param("att_id")
	attID, err := uuid.Parse(attIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "400",
			Message: "Invalid attachment id",
		})
		return
	}
	if err := h.Repo.Delete(attID); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "500",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
