package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"ticket-system/backend/internal/config"
	"ticket-system/backend/internal/delivery"
	"ticket-system/backend/internal/middleware"
	"ticket-system/backend/internal/repository"
	"ticket-system/backend/internal/usecase"
)

func main() {
	cfg := config.LoadConfig()

	db, err := repository.NewDB(cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := usecase.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := delivery.NewAuthHandler(authService)

	userHandler := delivery.NewUserHandler(userRepo)

	depRepo := repository.NewDepartmentRepository(db)
	statRepo := repository.NewTicketStatusRepository(db)
	prioRepo := repository.NewTicketPriorityRepository(db)
	depService := usecase.NewDepartmentService(depRepo)
	statService := usecase.NewTicketStatusService(statRepo)
	prioService := usecase.NewTicketPriorityService(prioRepo)

	dictHandler := delivery.NewDictionaryHandler(depService, statService, prioService)

	ticketRepo := repository.NewTicketRepository(db)
	historyRepo := repository.NewTicketHistoryRepository(db)
	ticketHandler := delivery.NewTicketHandler(
		ticketRepo,
		statRepo,
		prioRepo,
		depRepo,
		userRepo,
		historyRepo,
	)

	historyHandler := delivery.NewTicketHistoryHandler(historyRepo)

	commentRepo := repository.NewTicketCommentRepository(db)
	commentHandler := delivery.NewTicketCommentHandler(commentRepo)

	attachmentRepo := repository.NewTicketAttachmentRepository(db)
	attachmentHandler := delivery.NewTicketAttachmentHandler(attachmentRepo)

	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(middleware.ErrorLogger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/register", authHandler.Register)
	api.GET("/departments", dictHandler.DepartmentsList)
	api.GET("/ticket_statuses", dictHandler.TicketStatusesList)
	api.GET("/ticket_priorities", dictHandler.TicketPrioritiesList)
	api.GET("/users", userHandler.ListUsers)

	protected := api.Group("")
	protected.Use(middleware.JWT(cfg.JWTSecret))
	// Ticket
	protected.GET("/tickets", ticketHandler.GetTickets)
	protected.POST("/tickets", ticketHandler.CreateTicket)
	protected.GET("/tickets/:id", ticketHandler.GetTicketByID)
	protected.PATCH("/tickets/:id", ticketHandler.UpdateTicket)
	protected.DELETE("/tickets/:id", ticketHandler.DeleteTicket)
	// Ticket history
	protected.GET("/tickets/:id/history", historyHandler.GetHistory)
	// Ticket comments
	protected.GET("/tickets/:id/comments", commentHandler.GetComments)
	protected.POST("/tickets/:id/comments", commentHandler.AddComment)
	// Ticket attachments
	protected.GET("/tickets/:id/attachments", attachmentHandler.GetAttachments)
	protected.POST("/tickets/:id/attachments", attachmentHandler.AddAttachment)
	protected.GET("/tickets/:id/attachments/:att_id", attachmentHandler.GetAttachmentByID)
	protected.DELETE("/tickets/:id/attachments/:att_id", attachmentHandler.DeleteAttachment)

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
