package middleware

import (
	"fmt"
	"net/http"
	"prakerja-sesi-7/db"
	"prakerja-sesi-7/entity"
	"prakerja-sesi-7/utils/internal_jwt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Authentication(ctx *gin.Context) {
	jwtToken := ctx.Request.Header.Get("Authorization")

	fmt.Printf("\ntoken: %s\n", jwtToken)

	claim, err := internal_jwt.ValidateToken(jwtToken)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
			"errorMessage": "unauthorized",
		})
		return
	}

	fmt.Println("claim:", claim)

	ctx.Set("userId", claim["id"])

	ctx.Next()
}

func Authorization(ctx *gin.Context) {
	pg := db.GetDB()

	photoId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"messsage": "id has to be a number",
		},
		)
		return
	}

	var photo entity.Photo

	err = pg.Debug().First(&photo, "id = ?", photoId).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]string{
			"errorMessage": "photo not found",
		})
		return
	}

	userId, ok := ctx.MustGet("userId").(float64)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]string{
			"errorMessage": "forbidden access",
		})
		return
	}

	if uint(userId) != photo.UserID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]string{
			"errorMessage": "forbidden access",
		})
		return
	}

	ctx.Next()
}

func Authorization1(ctx *gin.Context) {
	pg := db.GetDB()

	commentId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"messsage": "id has to be a number",
		},
		)
		return
	}

	var comment entity.Comment

	err = pg.Debug().First(&comment, "id = ?", commentId).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]string{
			"errorMessage": "comment not found",
		})
		return
	}

	userId, ok := ctx.MustGet("userId").(float64)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]string{
			"errorMessage": "forbidden access",
		})
		return
	}

	if uint(userId) != comment.UserID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]string{
			"errorMessage": "forbidden access",
		})
		return
	}

	ctx.Next()
}

func Authorization2(ctx *gin.Context) {
	pg := db.GetDB()

	socialmediaId, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"messsage": "id has to be a number",
		},
		)
		return
	}

	var socialmedias entity.SocialMedia

	err = pg.Debug().First(&socialmedias, "id = ?", socialmediaId).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]string{
			"errorMessage": "social media not found",
		})
		return
	}

	userId, ok := ctx.MustGet("userId").(float64)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]string{
			"errorMessage": "forbidden access",
		})
		return
	}

	if uint(userId) != socialmedias.UserID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, map[string]string{
			"errorMessage": "forbidden access",
		})
		return
	}

	ctx.Next()
}
