package main

import (
	"net/http"
	"prakerja-sesi-7/db"
	"prakerja-sesi-7/dto"
	"prakerja-sesi-7/entity"
	"prakerja-sesi-7/middleware"
	"prakerja-sesi-7/utils/internal_jwt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type RegisterResponse struct {
	Status string `json:"status"`
	Data   struct {
		Age      uint   `json:"age"`
		Email    string `json:"email"`
		ID       uint   `json:"id"`
		Username string `json:"username"`
	} `json:"data"`
}

type LoginResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data"`
}

type UpdateUserResponse struct {
	Status   string `json:"status"`
	UserData struct {
		ID        uint      `json:"id"`
		Email     string    `json:"email"`
		Username  string    `json:"username"`
		Age       uint      `json:"age"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
}

type CreatePhotoResponse struct {
	Status   string `json:"status"`
	UserData struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Caption   string    `json:"caption"`
		PhotoURL  string    `json:"photo_url"`
		UserID    uint      `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"data"`
}

type UpdatePhotoResponse struct {
	Status   string `json:"status"`
	UserData struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Caption   string    `json:"caption"`
		PhotoURL  string    `json:"photo_url"`
		UserID    uint      `json:"user_id"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
}

type Photo struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title" gorm:"not null; type:varchar(255)"`
	Caption   string    `json:"caption" gorm:"not null; type:varchar(255)"`
	PhotoURL  string    `json:"photo_url" gorm:"not null"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user" gorm:"foreignkey:UserID"`
}

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type UserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UserResponseComment struct {
	ID       uint   `'json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UserResponseSocial struct {
	ID       uint   `'json:"id"`
	Username string `json:"username"`
	//ProfileImageURL string `json:"profile_image_url"`
}

type PhotoCommentResponse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url"`
	UserID   uint   `json:"user_id"`
}

type PhotoResponse struct {
	ID        uint         `json:"id"`
	Title     string       `json:"title"`
	Caption   string       `json:"caption"`
	PhotoURL  string       `json:"photo_url"`
	UserID    uint         `json:"user_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	User      UserResponse `json:"user"`
}

type CreateCommentResponse struct {
	Status   string `json:"status"`
	UserData struct {
		ID        uint      `json:"id"`
		Message   string    `json:"message"`
		PhotoID   uint      `json:"photo_id"`
		UserID    uint      `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"data"`
}

type Comment struct {
	ID        uint      `json:"id"`
	Message   string    `json:"message" gorm:"not null"`
	PhotoID   uint      `json:"photo_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user" gorm:"foreignkey:UserID"`
	Photo     Photo     `json:"photo" gorm:"foreignkey:PhotoID"`
}

type CommentResponse struct {
	ID        uint                 `json:"id"`
	Message   string               `json:"message"`
	PhotoID   uint                 `json:"photo_id"`
	UserID    uint                 `json:"user_id"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	User      UserResponseComment  `json:"user"`
	Photo     PhotoCommentResponse `json:"photo"`
}

type UpdateCommentResponse struct {
	Status   string `json:"status"`
	UserData struct {
		ID        uint      `json:"id"`
		Message   string    `json:"message"`
		PhotoID   uint      `json:"photo_id"`
		UserID    uint      `json:"user_id"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
}

type SocialMedia struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name" gorm:"not null"`
	SocialMediaURL string    `json:"social_media_url" gorm:"not null"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	User           User      `json:"user" gorm:"foreignkey:UserID"`
}

type SocialMediaResponse struct {
	ID             uint               `json:"id"`
	Name           string             `json:"name"`
	SocialMediaURL string             `json:"social_media_url" `
	UserID         uint               `json:"user_id" `
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	User           UserResponseSocial `json:"user"`
}

type CreateSocialResponse struct {
	Status   string `json:"status"`
	UserData struct {
		ID             uint      `json:"id"`
		Name           string    `json:"name"`
		SocialMediaURL string    `json:"social_media_url"`
		UserID         uint      `json:"user_id"`
		CreatedAt      time.Time `json:"created_at"`
	} `json:"data"`
}

type UpdateSocialResponse struct {
	Status   string `json:"status"`
	UserData struct {
		ID             uint      `json:"id"`
		Name           string    `json:"name"`
		SocialMediaURL string    `json:"social_media_url"`
		UserID         uint      `json:"user_id"`
		UpdatedAt      time.Time `json:"Updated_at"`
	} `json:"data"`
}

func Register(ctx *gin.Context) {
	var payload dto.LoginRequestDto

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"messsage": err.Error(),
		})
		return
	}

	pg := db.GetDB()

	pass, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 8)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"messsage": err.Error(),
		})
		return
	}

	user := entity.User{
		Age:      payload.Age,
		Email:    payload.Email,
		Password: string(pass),
		Username: payload.Username,
	}

	if err := pg.Create(&user).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"messsage": err.Error(),
		})
		return
	}

	response := RegisterResponse{
		Status: "success",
		Data: struct {
			Age      uint   `json:"age"`
			Email    string `json:"email"`
			ID       uint   `json:"id"`
			Username string `json:"username"`
		}{
			Age:      user.Age,
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}

	ctx.JSON(http.StatusCreated, response)
}

func Login(ctx *gin.Context) {
	var payload dto.LoginRequestDto

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"messsage": err.Error(),
		})
		return
	}

	pg := db.GetDB()

	var user entity.User

	if err := pg.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"messsage": "invalid email/password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(payload.Password)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"messsage": "invalid email/password",
		})
		return
	}

	claim := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
	}

	token := internal_jwt.GenerateToken(claim)

	response := LoginResponse{
		Status: "success",
		Data: struct {
			Token string `json:"token"`
		}{
			Token: token,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

func UpdateUserData(ctx *gin.Context) {
	//authUserID := ctx.GetInt("userId")

	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User ID must be a number",
		})
		return
	}

	/*if uint(userID) != uint(authUserID) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Forbidden: You are not allowed to update this user",
		})
		return
	}*/

	var payload dto.UpdateUserRequestDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid JSON",
		})
		return
	}

	pg := db.GetDB()

	var user entity.User
	if err := pg.First(&user, userID).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	if payload.Email != "" {
		user.Email = payload.Email
	}
	if payload.Username != "" {
		user.Username = payload.Username
	}
	if payload.Age != 0 {
		user.Age = payload.Age
	}

	user.UpdatedAt = time.Now()

	if err := pg.Save(&user).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user",
		})
		return
	}

	response := UpdateUserResponse{
		Status: "success",
		UserData: struct {
			ID        uint      `json:"id"`
			Email     string    `json:"email"`
			Username  string    `json:"username"`
			Age       uint      `json:"age"`
			UpdatedAt time.Time `json:"updated_at"`
		}{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Age:       user.Age,
			UpdatedAt: user.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

func DeleteUserData(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"messsage": "user id has to be a number",
		},
		)
		return
	}

	pg := db.GetDB()

	user := entity.User{}

	if err := pg.Delete(&user, userId).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"messsage": err.Error(),
		},
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your account has been succesfully deleted",
	})

}
func GetPhotos(ctx *gin.Context) {
	var photos []Photo
	var photoResponses []PhotoResponse
	pg := db.GetDB()

	if err := pg.Preload("User").Find(&photos).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, photo := range photos {
		photoResponses = append(photoResponses, PhotoResponse{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoURL:  photo.PhotoURL,
			UserID:    photo.UserID,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			User: UserResponse{
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"photos": photoResponses,
	})
}

func CreatePhotos(ctx *gin.Context) {
	userId, ok := ctx.MustGet("userId").(float64)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"messsage": "something went wrong",
		})
		return
	}

	_ = userId

	var payload dto.CreatePhotoRequestDto

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"messsage": err.Error(),
		})
		return
	}

	pg := db.GetDB()

	photo := entity.Photo{
		Title:    payload.Title,
		Caption:  payload.Caption,
		PhotoURL: payload.PhotoURL,
		UserID:   uint(userId),
	}

	if err := pg.Create(&photo).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"messsage": err.Error(),
		})
		return
	}

	response := CreatePhotoResponse{
		Status: "success",
		UserData: struct {
			ID        uint      `json:"id"`
			Title     string    `json:"title"`
			Caption   string    `json:"caption"`
			PhotoURL  string    `json:"photo_url"`
			UserID    uint      `json:"user_id"`
			CreatedAt time.Time `json:"created_at"`
		}{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoURL:  photo.PhotoURL,
			UserID:    uint(userId),
			CreatedAt: photo.CreatedAt,
		},
	}

	ctx.JSON(http.StatusCreated, response)

}

func UpdatePhotos(ctx *gin.Context) {
	photoId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id has to be a number",
		})
		return
	}

	var payload dto.CreatePhotoRequestDto

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	pg := db.GetDB()

	var photo entity.Photo

	if err := pg.First(&photo, "id = ?", photoId).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "product not found",
		})
		return
	}

	userId, ok := ctx.MustGet("userId").(float64)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	if uint(userId) != photo.UserID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "forbidden access",
		})
		return
	}

	photo.Title = payload.Title
	photo.Caption = payload.Caption
	photo.PhotoURL = payload.PhotoURL

	if err := pg.Save(&photo).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := UpdatePhotoResponse{
		Status: "success",
		UserData: struct {
			ID        uint      `json:"id"`
			Title     string    `json:"title"`
			Caption   string    `json:"caption"`
			PhotoURL  string    `json:"photo_url"`
			UserID    uint      `json:"user_id"`
			UpdatedAt time.Time `json:"updated_at"`
		}{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoURL:  photo.PhotoURL,
			UserID:    uint(userId),
			UpdatedAt: photo.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

func DeletePhotos(ctx *gin.Context) {
	photoId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"messsage": "id has to be a number",
		},
		)
		return
	}

	pg := db.GetDB()

	var photo entity.Photo
	if err := pg.First(&photo, photoId).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Photo not found",
		})
		return
	}

	// Delete the photo from the database
	if err := pg.Delete(&photo).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}

func GetComments(ctx *gin.Context) {
	var comments []Comment
	var commentResponses []CommentResponse

	pg := db.GetDB()

	if err := pg.Preload("User").Preload("Photo").Find(&comments).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, comment := range comments {
		commentResponses = append(commentResponses, CommentResponse{
			ID:        comment.ID,
			Message:   comment.Message,
			PhotoID:   comment.PhotoID,
			UserID:    comment.UserID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			User: UserResponseComment{
				ID:       comment.User.ID,
				Email:    comment.User.Email,
				Username: comment.User.Username,
			},
			Photo: PhotoCommentResponse{
				ID:       comment.Photo.ID,
				Title:    comment.Photo.Title,
				Caption:  comment.Photo.Caption,
				PhotoURL: comment.Photo.PhotoURL,
				UserID:   comment.Photo.UserID,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"comments": commentResponses,
	})
}

func CreateComments(ctx *gin.Context) {
	userId, ok := ctx.MustGet("userId").(float64)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"messsage": "something went wrong",
		})
		return
	}

	_ = userId

	var payload dto.CreateCommentRequestDto

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"messsage": err.Error(),
		})
		return
	}

	if payload.PhotoID == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "photoId is required",
		})
	}

	pg := db.GetDB()

	comment := entity.Comment{
		Message: payload.Message,
		PhotoID: payload.PhotoID,
		UserID:  uint(userId),
	}

	if err := pg.Create(&comment).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"messsage": err.Error(),
		})
		return
	}

	response := CreateCommentResponse{
		Status: "success",
		UserData: struct {
			ID        uint      `json:"id"`
			Message   string    `json:"message"`
			PhotoID   uint      `json:"photo_id"`
			UserID    uint      `json:"user_id"`
			CreatedAt time.Time `json:"created_at"`
		}{
			ID:        comment.ID,
			Message:   comment.Message,
			PhotoID:   comment.PhotoID,
			UserID:    uint(userId),
			CreatedAt: comment.CreatedAt,
		},
	}

	ctx.JSON(http.StatusCreated, response)
}

func UpdateComments(ctx *gin.Context) {
	commentId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id has to be a number",
		})
		return
	}

	var payload dto.UpdateCommentRequestDto

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	pg := db.GetDB()

	var comment entity.Comment

	if err := pg.First(&comment, "id = ?", commentId).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "comment not found",
		})
		return
	}

	userId, ok := ctx.MustGet("userId").(float64)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	if uint(userId) != comment.UserID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "forbidden access",
		})
		return
	}

	comment.Message = payload.Message

	if err := pg.Save(&comment).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := UpdateCommentResponse{
		Status: "success",
		UserData: struct {
			ID        uint      `json:"id"`
			Message   string    `json:"message"`
			PhotoID   uint      `json:"photo_id"`
			UserID    uint      `json:"user_id"`
			UpdatedAt time.Time `json:"updated_at"`
		}{
			ID:        comment.ID,
			Message:   comment.Message,
			PhotoID:   comment.PhotoID,
			UserID:    uint(userId),
			UpdatedAt: comment.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

func DeleteComments(ctx *gin.Context) {
	commentId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"messsage": "id has to be a number",
		},
		)
		return
	}

	pg := db.GetDB()

	var comment entity.Comment
	if err := pg.First(&comment, commentId).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Comment not found",
		})
		return
	}

	// Delete the comment from the database
	if err := pg.Delete(&comment).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}

func GetSocials(ctx *gin.Context) {
	var socialmedias []SocialMedia
	var socialmediaResponses []SocialMediaResponse

	pg := db.GetDB()

	if err := pg.Preload("User").Find(&socialmedias).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, social := range socialmedias {
		socialmediaResponses = append(socialmediaResponses, SocialMediaResponse{
			ID:             social.ID,
			Name:           social.Name,
			SocialMediaURL: social.SocialMediaURL,
			UserID:         social.UserID,
			CreatedAt:      social.CreatedAt,
			UpdatedAt:      social.UpdatedAt,
			User: UserResponseSocial{
				ID:       social.User.ID,
				Username: social.User.Username,
				//ProfileImageURL: social.User,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"socialmedias": socialmediaResponses,
	})
}

func CreateSocials(ctx *gin.Context) {
	userId, ok := ctx.MustGet("userId").(float64)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"messsage": "something went wrong",
		})
		return
	}

	_ = userId

	var payload dto.CreateSocialRequestDto

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"messsage": err.Error(),
		})
		return
	}

	pg := db.GetDB()

	social := entity.SocialMedia{
		Name:           payload.Name,
		SocialMediaURL: payload.SocialMediaURL,
		UserID:         uint(userId),
	}

	if err := pg.Create(&social).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"messsage": err.Error(),
		})
		return
	}

	response := CreateSocialResponse{
		Status: "success",
		UserData: struct {
			ID             uint      `json:"id"`
			Name           string    `json:"name"`
			SocialMediaURL string    `json:"social_media_url"`
			UserID         uint      `json:"user_id"`
			CreatedAt      time.Time `json:"created_at"`
		}{
			ID:             social.ID,
			Name:           social.Name,
			SocialMediaURL: social.SocialMediaURL,
			UserID:         uint(userId),
			CreatedAt:      social.CreatedAt,
		},
	}

	ctx.JSON(http.StatusCreated, response)
}

func UpdateSocials(ctx *gin.Context) {
	socialId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id has to be a number",
		})
		return
	}

	var payload dto.CreateSocialRequestDto

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	pg := db.GetDB()

	var social entity.SocialMedia

	if err := pg.First(&social, "id = ?", socialId).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "social media not found",
		})
		return
	}

	userId, ok := ctx.MustGet("userId").(float64)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	if uint(userId) != social.UserID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "forbidden access",
		})
		return
	}

	social.Name = payload.Name
	social.SocialMediaURL = payload.SocialMediaURL

	if err := pg.Save(&social).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := UpdateSocialResponse{
		Status: "success",
		UserData: struct {
			ID             uint      `json:"id"`
			Name           string    `json:"name"`
			SocialMediaURL string    `json:"social_media_url"`
			UserID         uint      `json:"user_id"`
			UpdatedAt      time.Time `json:"Updated_at"`
		}{
			ID:             social.ID,
			Name:           social.Name,
			SocialMediaURL: social.SocialMediaURL,
			UserID:         uint(userId),
			UpdatedAt:      social.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

func DeleteSocials(ctx *gin.Context) {
	socialId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"messsage": "id has to be a number",
		},
		)
		return
	}

	pg := db.GetDB()

	var social entity.SocialMedia
	if err := pg.First(&social, socialId).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "social media not found",
		})
		return
	}

	if err := pg.Delete(&social).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}

func init() {
	db.InitializeDB()
	pg := db.GetDB()
	pg.AutoMigrate(entity.User{}, entity.Photo{}, entity.Comment{}, entity.SocialMedia{})
}

func main() {

	r := gin.Default()

	users := r.Group("/users")
	{
		users.POST("/register", Register)
		users.POST("/login", Login)

		// Routes requiring authentication for user actions
		users1 := users.Group("/")
		users1.Use(middleware.Authentication)
		{
			users1.PUT("/:id", middleware.Authorization, UpdateUserData)
			users1.DELETE("/:id", middleware.Authorization, DeleteUserData)
		}
	}

	photos := r.Group("/photos")
	{
		photos.Use(middleware.Authentication)
		photos.GET("", GetPhotos)
		photos.POST("", CreatePhotos)
		photos.PUT("/:id", middleware.Authorization, UpdatePhotos)
		photos.DELETE("/:id", middleware.Authorization, DeletePhotos)
	}

	comments := r.Group("/comments")
	{
		comments.Use(middleware.Authentication)
		comments.GET("", GetComments)
		comments.POST("", CreateComments)
		comments.PUT("/:id", middleware.Authorization1, UpdateComments)
		comments.DELETE("/:id", middleware.Authorization1, DeleteComments)

	}

	socialmedias := r.Group("/socialmedias")
	{
		socialmedias.Use(middleware.Authentication)
		socialmedias.GET("", GetSocials)
		socialmedias.POST("", CreateSocials)
		socialmedias.PUT("/:id", middleware.Authorization2, UpdateSocials)
		socialmedias.DELETE("/:id", middleware.Authorization2, DeleteSocials)

	}

	r.Run(":8080")

}
