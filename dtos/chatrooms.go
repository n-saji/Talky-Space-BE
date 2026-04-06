package dtos

type CreateChatroomRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsGroup     bool   `json:"is_group"`
	CreatedBy   string `json:"created_by" binding:"required"`
}

type ChatroomResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsGroup     bool   `json:"is_group"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type UpdateChatroomRequest struct {
	Id          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
}

type ChatroomInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsGroup     bool   `json:"is_group"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}