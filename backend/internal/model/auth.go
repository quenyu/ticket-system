package model

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  LoginUserDTO `json:"user"`
}

type LoginUserDTO struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	RoleID       string `json:"role_id"`
	DepartmentID int16  `json:"department_id"`
}

type UserRegisterRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Email        string `json:"email" binding:"required"`
	RoleID       string `json:"role_id" binding:"required"`
	DepartmentID int16  `json:"department_id" binding:"required"`
}

type UserRegisterResponse struct {
	UserID string `json:"user_id"`
}
