package dtos

type CreateMessageRequest struct {
	ChatroomID  string `json:"chatroom_id"`
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	Content     string `json:"content"`
}

type MessageResponse struct {
	Id         string `json:"id"`
	ChatroomID string `json:"chatroom_id"`
	UserId     string `json:"user_id"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"created_at"`
}
