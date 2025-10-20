package dtos

type CreateChatroomMemberRequest struct {
	ChatroomID string `json:"chatroom_id" binding:"required,uuid"`
	UserID     string `json:"user_id" binding:"required,uuid"`
}

type ChatroomMemberResponse struct {
	Id         string `json:"id"`
	ChatroomID string `json:"chatroom_id"`
	UserID     string `json:"user_id"`
	JoinedAt   int64  `json:"joined_at"`
}

type UpdateChatroomMemberRequest struct {
	Id         string `json:"id" binding:"required,uuid"`
	ChatroomID string `json:"chatroom_id" binding:"omitempty,uuid"`
	UserID     string `json:"user_id" binding:"omitempty,uuid"`
}
type DeleteChatroomMemberRequest struct {
	Id string `json:"id" binding:"required,uuid"`
}
type GetChatroomMembersRequest struct {
	ChatroomID string `json:"chatroom_id" binding:"required,uuid"`
}
type GetChatroomMembersResponse struct {
	Members []ChatroomMemberResponse `json:"members"`
}
type IsUserInChatroomRequest struct {
	ChatroomID string `json:"chatroom_id" binding:"required,uuid"`
	UserID     string `json:"user_id" binding:"required,uuid"`
}
type IsUserInChatroomResponse struct {
	IsMember bool `json:"is_member"`
}
