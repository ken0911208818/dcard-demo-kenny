package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ken0911208818/dcard-demo-kenny/lib/auth"
	"github.com/ken0911208818/dcard-demo-kenny/lib/constant"
	"github.com/ken0911208818/dcard-demo-kenny/lib/middleware"
	"github.com/ken0911208818/dcard-demo-kenny/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		middleware.SendErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	user := model.User{}
	db := c.MustGet(constant.Db).(*gorm.DB)
	result := db.Where(`Email = ?`, input.Email).First(&user)
	if err := result.Error; err != nil {
		middleware.SendErrorResponse(c, http.StatusNotFound, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(input.Password)); err != nil {
		middleware.SendErrorResponse(c, http.StatusUnauthorized, errors.New("Incorrect Email  Password"))
		return
	}

	if newToken, err := auth.Sign(user.Id); err != nil {
		middleware.SendErrorResponse(c, http.StatusInternalServerError, err)
		return
	} else {
		// update JWT Token
		c.Header("Authorization", newToken)
		// allow CORS
		c.Header("Access-Control-Expose-Headers", "Authorization")
		c.Set(constant.StatusCode, http.StatusCreated)
		c.Set(constant.Output, map[string]interface{}{"userId": user.Id})
	}
}
