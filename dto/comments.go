package dto

type CreateCommentRequestDto struct {
	Message string `json:"message" binding:"required"`
	PhotoID uint   `json:"photo_id" binding:"required"`
}

type UpdateCommentRequestDto struct {
	Message string `json:"message" binding:"required"`
}
