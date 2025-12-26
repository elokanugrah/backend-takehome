package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=4"`
	Email    string `json:"email" binding:"required,min=6"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
