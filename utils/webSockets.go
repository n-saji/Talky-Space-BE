package utils

import (
	"encoding/json"
	"fmt"
	"talky-space-be/config"
	"talky-space-be/daos"
	"talky-space-be/dtos"
	"talky-space-be/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	Clients    map[uuid.UUID]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

type MessagePayload struct {
	Type        string    `json:"type"`
	ChatroomID  uuid.UUID `json:"chatroom_id"`
	SenderID    uuid.UUID `json:"user_id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	Content     string    `json:"content"`
	CreatedAt   int64     `json:"created_at"`
	Source      string    `json:"source,omitempty"`
}

// Exported Hub instance
var HubInstance = &Hub{
	Clients:    make(map[uuid.UUID]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Broadcast:  make(chan []byte),
}

// Start the hub in background (call this in main.go)
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = client

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}

		case message := <-h.Broadcast:
			var payload MessagePayload
			if err := json.Unmarshal(message, &payload); err != nil {
				fmt.Println("Error unmarshaling message:", err)
				continue
			}

			fmt.Println("message from socket: ", payload)

			dbConn := config.DBInit()
			db := daos.New(dbConn)
			members, err := db.GetChatroomMembersByChatroomID(payload.ChatroomID.String())
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			for _, member := range members {

				if conn, ok := h.Clients[member.UserID]; ok {

					if payload.Source == "server" {
						// Skip storing message if it's from server to avoid duplication
						continue
					}

					messageModel := models.CreateMessageRequestToMessageModel(&dtos.CreateMessageRequest{
						ChatroomID:  payload.ChatroomID.String(),
						SenderID:    payload.SenderID.String(),
						Content:     payload.Content,
						RecipientID: payload.RecipientID.String(),
					})
					_, err = db.CreateMessage(*messageModel)
					if err != nil {
						fmt.Println("Error storing message:", err)
					}

					if member.UserID == payload.SenderID {
						continue
					}
					conn.Send <- message
				}
			}
		}
	}
}

// Add helper function to broadcast directly
func BroadcastMessage(payload MessagePayload) {
	data, _ := json.Marshal(payload)
	HubInstance.Broadcast <- data
}

func (c *Client) ReadPump() {
	defer func() {
		HubInstance.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		HubInstance.Broadcast <- msg
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
