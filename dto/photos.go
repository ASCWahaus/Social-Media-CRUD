package dto

type CreatePhotoRequestDto struct {
	Title    string `json:"title" binding:"required"`
	Caption  string `json:"caption" binding:"required"`
	PhotoURL string `json:"photo_url" binding:"required"`
}
