package delivery

import (
	"net/http"
	"ticket-system/backend/internal/model"
	"ticket-system/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type DictionaryHandler struct {
	Departments      *usecase.DepartmentService
	TicketStatuses   *usecase.TicketStatusService
	TicketPriorities *usecase.TicketPriorityService
}

func NewDictionaryHandler(dep *usecase.DepartmentService, stat *usecase.TicketStatusService, prio *usecase.TicketPriorityService) *DictionaryHandler {
	return &DictionaryHandler{
		Departments:      dep,
		TicketStatuses:   stat,
		TicketPriorities: prio,
	}
}

func (h *DictionaryHandler) DepartmentsList(c *gin.Context) {
	items, err := h.Departments.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "500",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *DictionaryHandler) TicketStatusesList(c *gin.Context) {
	items, err := h.TicketStatuses.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "500",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *DictionaryHandler) TicketPrioritiesList(c *gin.Context) {
	items, err := h.TicketPriorities.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "500",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, items)
}
