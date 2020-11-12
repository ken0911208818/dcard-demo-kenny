package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ken0911208818/dcard-demo-kenny/lib/auth"
	"github.com/ken0911208818/dcard-demo-kenny/lib/constant"
	"github.com/ken0911208818/dcard-demo-kenny/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

func UserCreate(c *gin.Context)  {
	var user struct{
		model.User
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.Set(constant.StatusCode, http.StatusBadRequest)
		c.Set(constant.Error, err)
		return
	}
	uuid, err := uuid.NewRandom()
	if err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
		return
	}
	user.Id = uuid.String()


	if digest, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
		return
	} else {
		user.PasswordDigest = string(digest)
	}
	db := c.MustGet(constant.Db).(*gorm.DB)
	q := `insert into users(id, email, password_digest, name)
			select ?, ?, ?, ?
			where not exists (select 1 from users where email = ?)`
	result := db.Exec(q, user.Id, user.Email, user.PasswordDigest, user.Name, user.Email)
	if  err = result.Error; err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
		return
	}
	if result.RowsAffected == 0 {
		c.Set(constant.StatusCode, http.StatusForbidden)
		c.Set(constant.Error, errors.New("the email is already used"))
		return
	}

	if newToken, err := auth.Sign(user.Id); err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
	} else {
		// update JWT Token
		c.Header("Authorization", newToken)
		// allow CORS
		c.Header("Access-Control-Expose-Headers", "Authorization")
		c.Set(constant.StatusCode, http.StatusCreated)
		c.Set(constant.Output, map[string]interface{}{"userId": user.Id})
	}
}
