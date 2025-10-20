package handlers

import (
	"talky-space-be/dtos"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RoutingMessage(rg *gin.RouterGroup) {
	message := rg.Group("/messages")
	{
		{
			message.POST("/create", h.CreateMessage)
		}
	}
}

func (h *Handler) CreateMessage(c *gin.Context) {
	req := &dtos.CreateMessageRequest{}
	req.SenderID = c.GetString("user_id")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	res, err := h.service.CreateMessage(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create message with error: " + err.Error()})
		return
	}
	c.JSON(200, res)
}
