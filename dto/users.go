package dto

type LoginRequestDto struct {
	Age      uint   `json:"age" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required"`
}

type UpdateUserRequestDto struct {
	Age      uint   `json:"age"`
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username"`
}

type LoginResponseDto struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}
