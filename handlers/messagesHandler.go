package handlers

import (
	"talky-space-be/dtos"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RoutingMessage(rg *gin.RouterGroup) {
	message := rg.Group("/messages")
	{
		{
			message.POST("/new-chat", h.InitiateNewChat)
			message.GET("/chatroom/:cid",h.FetchChatroomMessages)
		}
	}
}

func (h *Handler) InitiateNewChat(c *gin.Context) {
	req := &dtos.CreateMessageRequest{}
	req.SenderID = c.GetString("user_id")
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	res, err := h.service.InitiateNewChat(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create message with error: " + err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) FetchChatroomMessages(c *gin.Context) {
	chatroomID := c.Param("cid")
	res, err := h.service.FetchMessagesByChatroomID(chatroomID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch messages with error: " + err.Error()})
		return
	}
	c.JSON(200, res)
}