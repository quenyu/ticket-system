package delivery

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/model"
	"ticket-system/backend/internal/repository"
	"ticket-system/backend/internal/usecase"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	TicketRepo     usecase.TicketRepository
	StatusRepo     usecase.TicketStatusRepository
	PriorityRepo   usecase.TicketPriorityRepository
	DepartmentRepo usecase.DepartmentRepository
	UserRepo       *repository.UserRepository
	HistoryRepo    *repository.TicketHistoryRepository
}

type TicketFilter struct {
	StatusID     *int16
	AssigneeID   *string
	DepartmentID *int16
	Q            string
	Limit        int
	Offset       int
}

func NewTicketHandler(
	repo usecase.TicketRepository,
	statusRepo usecase.TicketStatusRepository,
	priorityRepo usecase.TicketPriorityRepository,
	departmentRepo usecase.DepartmentRepository,
	userRepo *repository.UserRepository,
	historyRepo *repository.TicketHistoryRepository,
) *TicketHandler {
	return &TicketHandler{
		TicketRepo:     repo,
		StatusRepo:     statusRepo,
		PriorityRepo:   priorityRepo,
		DepartmentRepo: departmentRepo,
		UserRepo:       userRepo,
		HistoryRepo:    historyRepo,
	}
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

	// Enrich tickets with labels/names
	statuses, _ := h.StatusRepo.List()
	priorities, _ := h.PriorityRepo.List()
	departments, _ := h.DepartmentRepo.List()
	users, _ := h.UserRepo.List()

	statusMap := make(map[int16]string)
	for _, s := range statuses {
		statusMap[s.ID] = s.Label
	}
	priorityMap := make(map[int16]string)
	for _, p := range priorities {
		priorityMap[p.ID] = p.Label
	}
	departmentMap := make(map[int16]string)
	for _, d := range departments {
		departmentMap[d.ID] = d.Name
	}
	userMap := make(map[string]string)
	for _, u := range users {
		userMap[u.ID.String()] = u.Username
	}

	var result []gin.H
	for _, t := range tickets {
		result = append(result, gin.H{
			"ticket_id":       t.ID,
			"title":           t.Title,
			"description":     t.Description,
			"status_id":       t.StatusID,
			"status_label":    statusMap[t.StatusID],
			"priority_id":     t.PriorityID,
			"priority_label":  priorityMap[t.PriorityID],
			"department_id":   t.DepartmentID,
			"department_name": departmentMap[t.DepartmentID],
			"assignee_id":     t.AssigneeID,
			"assignee_name":   userMap[t.AssigneeID.String()],
			"creator_id":      t.CreatorID,
			"created_at":      t.CreatedAt,
			"updated_at":      t.UpdatedAt,
			"deleted_at":      t.DeletedAt,
		})
	}
	c.JSON(http.StatusOK, result)
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

	if req.AssigneeID == uuid.Nil {
		var raw map[string]interface{}
		if err := c.ShouldBindJSON(&raw); err == nil {
			if val, ok := raw["assignee_id"].(string); ok && val != "" {
				parsed, err := uuid.Parse(val)
				if err == nil {
					req.AssigneeID = parsed
				}
			}
		}
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

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "INVALID_UUID",
			Message: "user_id in context is not a string",
		})
		return
	}

	parsedUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "INVALID_UUID",
			Message: "Invalid user_id in token",
		})
		return
	}
	req.CreatorID = parsedUUID

	if err := h.TicketRepo.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "DATABASE_ERROR",
			Message: "Failed to create ticket",
		})
		return
	}

	history := &domain.TicketHistory{
		TicketID:        req.ID,
		TicketCreatedAt: req.CreatedAt,
		ChangedAt:       time.Now(),
		ChangedBy:       req.CreatorID,
		FieldName:       "ticket",
		OldValue:        "",
		NewValue:        "created",
	}
	_ = h.HistoryRepo.Create(history)
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

	oldTicket, err := h.TicketRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:    "TICKET_NOT_FOUND",
			Message: "Ticket not found",
		})
		return
	}

	if err := h.TicketRepo.Update(&req); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "DATABASE_ERROR",
			Message: "Failed to update ticket",
		})
		return
	}

	userID, _ := c.Get("user_id")
	userUUID, _ := uuid.Parse(userID.(string))
	now := time.Now()
	if oldTicket.StatusID != req.StatusID {
		h.HistoryRepo.Create(&domain.TicketHistory{
			TicketID:        req.ID,
			TicketCreatedAt: oldTicket.CreatedAt,
			ChangedAt:       now,
			ChangedBy:       userUUID,
			FieldName:       "status",
			OldValue:        strconv.Itoa(int(oldTicket.StatusID)),
			NewValue:        strconv.Itoa(int(req.StatusID)),
		})
	}
	if oldTicket.PriorityID != req.PriorityID {
		h.HistoryRepo.Create(&domain.TicketHistory{
			TicketID:        req.ID,
			TicketCreatedAt: oldTicket.CreatedAt,
			ChangedAt:       now,
			ChangedBy:       userUUID,
			FieldName:       "priority",
			OldValue:        strconv.Itoa(int(oldTicket.PriorityID)),
			NewValue:        strconv.Itoa(int(req.PriorityID)),
		})
	}
	if oldTicket.Description != req.Description {
		h.HistoryRepo.Create(&domain.TicketHistory{
			TicketID:        req.ID,
			TicketCreatedAt: oldTicket.CreatedAt,
			ChangedAt:       now,
			ChangedBy:       userUUID,
			FieldName:       "description",
			OldValue:        oldTicket.Description,
			NewValue:        req.Description,
		})
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
