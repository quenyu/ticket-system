package delivery

import (
	"io"
	"net/http"
	"strings"
	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/model"
	"ticket-system/backend/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketAttachmentHandler struct {
	Repo       usecase.TicketAttachmentRepository
	TicketRepo usecase.TicketRepository
}

func NewTicketAttachmentHandler(repo usecase.TicketAttachmentRepository, ticketRepo usecase.TicketRepository) *TicketAttachmentHandler {
	return &TicketAttachmentHandler{Repo: repo, TicketRepo: ticketRepo}
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

	ticket, err := h.TicketRepo.GetByID(ticketID)
	if err != nil || ticket == nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "400",
			Message: "Ticket not found",
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{
			Code:    "400",
			Message: "No file uploaded",
		})
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIError{
			Code:    "500",
			Message: "Failed to read file",
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

	att := domain.TicketAttachment{
		ID:              uuid.New(),
		TicketID:        ticketID,
		TicketCreatedAt: ticket.CreatedAt,
		Filename:        header.Filename,
		FileData:        fileData,
		UploadedBy:      userID,
		UploadedAt:      time.Now().UTC(),
	}

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

func (h *TicketAttachmentHandler) DownloadAttachment(c *gin.Context) {
	attIDStr := c.Param("att_id")
	attID, err := uuid.Parse(attIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Code: "400", Message: "Invalid attachment id"})
		return
	}
	att, err := h.Repo.GetByID(attID)
	if err != nil || att.FileData == nil {
		c.JSON(http.StatusNotFound, model.APIError{Code: "404", Message: "Attachment not found"})
		return
	}
	filename := att.Filename
	contentType := "application/octet-stream"
	if strings.HasSuffix(strings.ToLower(filename), ".jpg") || strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(strings.ToLower(filename), ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(strings.ToLower(filename), ".gif") {
		contentType = "image/gif"
	}
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Data(http.StatusOK, contentType, att.FileData)
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

	roleIDRaw, exists := c.Get("role_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.APIError{
			Code:    "401",
			Message: "No role_id in token",
		})
		return
	}
	roleID, ok := roleIDRaw.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.APIError{
			Code:    "401",
			Message: "role_id in token is not a string",
		})
		return
	}
	isAdmin := roleID == "00000000-0000-0000-0000-000000000002"

	att, err := h.Repo.GetByID(attID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIError{
			Code:    "404",
			Message: "Attachment not found",
		})
		return
	}
	if att.UploadedBy != userID && !isAdmin {
		c.JSON(http.StatusForbidden, model.APIError{
			Code:    "403",
			Message: "You can only delete your own attachments",
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
