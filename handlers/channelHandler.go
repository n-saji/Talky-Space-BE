package handlers

import (
	"net/http"
	"talky-space-be/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RoutingChannel(rc *gin.RouterGroup) {
	rc.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	channel := rc.Group("/channel")
	{
		protected := channel.Group("/")
		protected.Use(middleware.AuthMiddleware())

		protected.GET("/:id", func(c *gin.Context) {
			// Handle getting channel by ID
		})
		protected.POST("/", func(c *gin.Context) {
			// Handle creating a new channel
		})
	}
}
