package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ken0911208818/dcard-demo-kenny/lib/constant"
	"github.com/ken0911208818/dcard-demo-kenny/lib/middleware"
	"github.com/ken0911208818/dcard-demo-kenny/model"
	"gorm.io/gorm"
	"net/http"
)

func ParisCreate(c *gin.Context) {
	// 當需要取得一個物件時 使用MustGet 但如果是單一的值直接使用GetInt GetString 看當下的value
	// MustGet 若是key not exist 會發生 panic
	db := c.MustGet(constant.Session).(*gorm.DB)
	userId := c.GetString(constant.UserId)
	//	從資料庫選擇一個使用者 沒有被配對過在pairs
	q := `SELECT * FROM users WHERE id != ? and not exists (SELECT 1 FROM pairs WHERE user_id_one = ?) ORDER BY random() LIMIT 1;`
	user := model.User{}
	result := db.Find(&user).Exec(q, userId, userId)
	if err := result.Error; err != nil {
		middleware.SendErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	if found := result.RowsAffected; found == 0 {
		middleware.SendErrorResponse(c, http.StatusForbidden, errors.New("already have pair"))
		return
	}
	var pair model.Pair

	pair.UserIdOne = userId
	pair.UserIdTwo = user.Id

	result = db.Create(&pair)

	if err := result.Error; err != nil {
		middleware.SendErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if affect := result.RowsAffected; affect == 0 {
		middleware.SendErrorResponse(c, http.StatusForbidden, errors.New("already have pair"))
		return
	}
	c.Set(constant.StatusCode, http.StatusCreated)
	c.Set(constant.Output, map[string]interface{}{"user_id_two": pair.UserIdTwo})
}

func PairGetOne(c *gin.Context) {
	userId := c.GetString(constant.UserId)
	db := c.MustGet(constant.Session).(*gorm.DB)
	var pair model.Pair
	result := db.Where(`user_id_one = ?`, userId).First(&pair)
	if err := result.Error; err != nil {
		middleware.SendErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	if affect := result.RowsAffected; affect == 0 {
		middleware.SendErrorResponse(c, http.StatusForbidden, errNotFound)
		return
	}
	c.Set(constant.StatusCode, http.StatusCreated)
	c.Set(constant.Output, map[string]interface{}{"user_id_two": pair.UserIdTwo})
}
