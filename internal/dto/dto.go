package dto

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}
