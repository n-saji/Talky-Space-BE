package handlers

import (
	"talky-space-be/middleware"
	"talky-space-be/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *service.Service
}

func New(dbConn *gorm.DB) *Handler {
	return &Handler{
		service: service.New(dbConn),
	}
}

func (h *Handler) GetRouter() *gin.Engine {

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5050", "http://localhost:3000", "http://localhost:5173","http://192.168.1.54:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Set-Cookie", "token", "account_id"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie", "token", "account_id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	h.RoutingUser(&router.RouterGroup)
	h.AuthenticationChannel(&router.RouterGroup)
	h.RoutingChannel(&router.RouterGroup)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Add protected routes here
		h.RoutingMessage(protected)
		h.RoutingWebSockets(protected)
		h.ChatroomChannel(protected)
	}

	return router
}
