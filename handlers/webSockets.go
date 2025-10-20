package handlers

import (
	"net/http"
	"talky-space-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handler) RoutingWebSockets(rg *gin.RouterGroup) {
	websockets := rg.Group("/ws")
	{
		websockets.GET("/connect", h.HandleWebSocket)
	}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &utils.Client{
		UserID: uuid.MustParse(userID.(string)),
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	utils.HubInstance.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
