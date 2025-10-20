package dtos


type CreateUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Name        string `json:"name" binding:"required"`
	Avatar      string `json:"avatar"`
	PhoneNumber string `json:"phone_number"`
}

type UserResponse struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	PhoneNumber string `json:"phone_number"`
}

type UserAuthenticationResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type UpdateUserRequest struct {
	Email       string `json:"email" binding:"omitempty,email"`
	Password    string `json:"password" binding:"omitempty,min=6"`
	Name        string `json:"name" binding:"omitempty"`
	Avatar      string `json:"avatar" binding:"omitempty"`
	PhoneNumber string `json:"phone_number" binding:"omitempty"`
}

