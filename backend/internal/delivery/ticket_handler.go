package delivery

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/model"
	"ticket-system/backend/internal/usecase"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	TicketRepo usecase.TicketRepository
}

type TicketFilter struct {
	StatusID     *int16
	AssigneeID   *string
	DepartmentID *int16
	Q            string
	Limit        int
	Offset       int
}

func NewTicketHandler(repo usecase.TicketRepository) *TicketHandler {
	return &TicketHandler{TicketRepo: repo}
}

func (h *TicketHandler) validateTicket(ticket *domain.Ticket) *model.APIError {
	if strings.TrimSpace(ticket.Title) == "" {
		return &model.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Title is required",
		}
	}

	if len(ticket.Title) > 255 {
		return &model.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Title must be less than 255 characters",
		}
	}

	if strings.TrimSpace(ticket.Description) == "" {
		return &model.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Description is required",
		}
	}

	if len(ticket.Description) > 10000 {
		return &model.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Description must be less than 10000 characters",
		}
	}

	if ticket.StatusID <= 0 {
		return &model.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Valid status is required",
		}
	}

	if ticket.PriorityID <= 0 {
		return &model.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Valid priority is required",
		}
	}

	if ticket.DepartmentID <= 0 {
		return &model.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Valid department is required",
		}
	}

	return nil
}

func (h *TicketHandler) GetTickets(c *gin.Context) {
	var filter usecase.TicketFilter
	if v := c.Query("status_id"); v != "" {
		var id int16
		if _, err := fmt.Sscan(v, &id); err == nil {
			filter.StatusID = &id
		}
	}
	if v := c.Query("assignee_id"); v != "" {
		filter.AssigneeID = &v
	}
	if v := c.Query("department_id"); v != "" {
		var id int16
		if _, err := fmt.Sscan(v, &id); err == nil {
			filter.DepartmentID = &id
		}
	}
	filter.Q = c.Query("q")
	if v := c.Query("limit"); v != "" {
		fmt.Sscan(v, &filter.Limit)
	}
	if v := c.Query("offset"); v != "" {
		fmt.Sscan(v, &filter.Offset)
	}

	if filter.AssigneeID != nil {
		if _, err := uuid.Parse(*filter.AssigneeID); err != nil {
			c.JSON(http.StatusBadRequest, model.APIError{
				Code:    "INVALID_UUID",
				Message: "Invalid assignee_id format",
			})
			return
		}
	}

	if filter.StatusID != nil {
		if _, err := strconv.ParseInt(fmt.Sprintf("%d", *filter.StatusID), 10, 16); err != nil {
			c.JSON(http.StatusBadRequest, model.APIError{
				Code:    "INVALID_STATUS_ID",
				Message: "Invalid status_id format",
			})
			return
		}
	}

	if filter.DepartmentID != nil {
		if _, err := strconv.ParseInt(fmt.Sprintf("%d", *filter.DepartmentID), 10, 16); err != nil {
			c.JSON(http.StatusBadRequest, model.APIError{
				Code:    "INVALID_DEPARTMENT_ID",
				Message: "Invalid department_id format",
			})
			return
		}
	}

	tickets, err := h.TicketRepo.Search(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "DATABASE_ERROR",
			Message: "Failed to get tickets",
		})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req domain.Ticket
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "INVALID_JSON",
			Message: "Invalid request body",
		})
		return
	}

	if validationError := h.validateTicket(&req); validationError != nil {
		c.JSON(http.StatusBadRequest, *validationError)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}
	req.CreatorID = userID.(uuid.UUID)

	if err := h.TicketRepo.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "DATABASE_ERROR",
			Message: "Failed to create ticket",
		})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *TicketHandler) GetTicketByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "MISSING_ID",
			Message: "Ticket ID is required",
		})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "INVALID_UUID",
			Message: "Invalid ticket ID format",
		})
		return
	}

	ticket, err := h.TicketRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:    "TICKET_NOT_FOUND",
			Message: "Ticket not found",
		})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "MISSING_ID",
			Message: "Ticket ID is required",
		})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "INVALID_UUID",
			Message: "Invalid ticket ID format",
		})
		return
	}

	var req domain.Ticket
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "INVALID_JSON",
			Message: "Invalid request body",
		})
		return
	}

	if validationError := h.validateTicket(&req); validationError != nil {
		c.JSON(http.StatusBadRequest, *validationError)
		return
	}

	req.ID = id

	if err := h.TicketRepo.Update(&req); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "DATABASE_ERROR",
			Message: "Failed to update ticket",
		})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "MISSING_ID",
			Message: "Ticket ID is required",
		})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "INVALID_UUID",
			Message: "Invalid ticket ID format",
		})
		return
	}

	if err := h.TicketRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "DATABASE_ERROR",
			Message: "Failed to delete ticket",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}
