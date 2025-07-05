package delivery

import (
	"fmt"
	"net/http"
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
	tickets, err := h.TicketRepo.Search(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.APIErrorDetail{Code: 500, Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req domain.Ticket
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.APIErrorDetail{Code: 400, Message: "Validation failed", Details: err.Error()},
		})
		return
	}
	if err := h.TicketRepo.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.APIErrorDetail{Code: 500, Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *TicketHandler) GetTicketByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.APIErrorDetail{Code: 400, Message: "Invalid ticket id"},
		})
		return
	}
	ticket, err := h.TicketRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIError{
			Error: model.APIErrorDetail{Code: 404, Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.APIErrorDetail{Code: 400, Message: "Invalid ticket id"},
		})
		return
	}
	var req domain.Ticket
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.APIErrorDetail{Code: 400, Message: "Validation failed", Details: err.Error()},
		})
		return
	}
	req.ID = id
	if err := h.TicketRepo.Update(&req); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.APIErrorDetail{Code: 500, Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Error: model.APIErrorDetail{Code: 400, Message: "Invalid ticket id"},
		})
		return
	}
	if err := h.TicketRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Error: model.APIErrorDetail{Code: 500, Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
