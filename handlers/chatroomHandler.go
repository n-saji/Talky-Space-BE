package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) ChatroomChannel(rg *gin.RouterGroup) {
	chatroomGrp := rg.Group("/chatrooms")
	{
		chatroomGrp.POST("/", h.CreateChatroom)
		chatroomGrp.GET("/", h.GetUserChatrooms)
		chatroomGrp.GET("/:chatroom_id", h.GetChatroomDetails)
		chatroomGrp.PUT("/:chatroom_id", h.UpdateChatroom)
		chatroomGrp.DELETE("/:chatroom_id", h.DeleteChatroom)
		chatroomGrp.GET("/find-by-users/user1/:uid1/user2/:uid2", h.FindChatroomByUsers)
	}
}

func (h *Handler) CreateChatroom(c *gin.Context) {
	// Implementation for creating a chatroom
}

func (h *Handler) GetUserChatrooms(c *gin.Context) {
	// Implementation for retrieving user's chatrooms

}

func (h *Handler) GetChatroomDetails(c *gin.Context) {
	// Implementation for retrieving chatroom details
	chatroom_id := c.Param("chatroom_id")

	res, err := h.service.GetChatroomByID(chatroom_id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) UpdateChatroom(c *gin.Context) {
	// Implementation for updating a chatroom
}

func (h *Handler) DeleteChatroom(c *gin.Context) {
	// Implementation for deleting a chatroom
}

func (h *Handler) FindChatroomByUsers(c *gin.Context) {
	uid1 := c.Param("uid1")
	uid2 := c.Param("uid2")

	res, err := h.service.FindChatroomByUsers(uid1, uid2)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}
