package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/ken0911208818/dcard-demo-kenny/lib/auth"
	"github.com/ken0911208818/dcard-demo-kenny/lib/constant"
	"github.com/ken0911208818/dcard-demo-kenny/lib/lua"
	"github.com/ken0911208818/dcard-demo-kenny/lib/vaildate"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

var (
	db          *gorm.DB
	redisClient *redis.Client
)

const (
	IPLimitPeriod     = 3600
	IPLimitTimeFormat = "2006-01-02 15:04:05"
	IPLimitMaximum    = 1000
)

func Init(database *gorm.DB, client *redis.Client) {
	db = database
	redisClient = client
}

func Plain() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constant.Db, db)
		c.Set(constant.StatusCode, nil)
		c.Set(constant.Output, nil)
		c.Set(constant.Error, nil)
		c.Next()

		statusCode := c.GetInt(constant.StatusCode)
		err := c.MustGet(constant.Error)
		output := c.MustGet(constant.Output)
		if err != nil {
			if validationErr, ok := err.(validator.ValidationErrors); ok {
				//Translate : 用途 將錯誤訊息依造定義的語言翻譯
				sendResponse(c, statusCode, map[string]interface{}{"error": validationErr.Translate(vaildate.BindingTrans)})
			} else {
				// 直接顯示錯誤訊息 不進行翻譯動作
				sendResponse(c, statusCode, map[string]interface{}{"error": err.(error).Error()})
			}
		} else {
			sendResponse(c, statusCode, output)
		}
	}
}

func sendResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
	c.Abort()
}

func SendErrorResponse(c *gin.Context, StatusCode int, err error) {
	c.Set(constant.StatusCode, StatusCode)
	c.Set(constant.Error, err)
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		userId, err := auth.Verify(token)
		if err != nil {
			SendErrorResponse(c, http.StatusUnauthorized, err)
		} else {
			if newToken, err := auth.Sign(userId); err != nil {
				SendErrorResponse(c, http.StatusInternalServerError, err)
			} else {
				c.Header("Authorization", newToken) // update JWT Token
				c.Set(constant.UserId, userId)
			}
		}
	}
}

func IPLimitIntercept() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString(constant.UserId)
		now := time.Now().Unix()
		key := c.Request.URL.Path + "-" + c.Request.Method + "-" + c.ClientIP() + "-" + userId
		script := redis.NewScript(lua.SCRIPT)
		args := []interface{}{now, IPLimitMaximum, IPLimitPeriod}

		value, err := script.Run(c, redisClient, []string{key}, args...).Result()
		// only when redis is disconnected or lua runtime error, error will show up. and it will be rollback.
		// if script's any redis operations are wrong, it will not get error because it is recognized as logical error
		// for example: wrong key
		if err != nil {
			sendResponse(c, http.StatusInternalServerError, err)
			return
		}

		result := value.([]interface{})
		remaining := result[0].(int64)
		// in normal situation: 0~9
		// otherwise, "-1" means too much requests in period
		if remaining == -1 {
			sendResponse(c, http.StatusTooManyRequests, err)
			return
		}
		t := result[1].(int64)
		reset := time.Unix(t, 0).Format(IPLimitTimeFormat)

		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Header("X-RateLimit-Reset", reset)
	}
}

func Tx() gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.Session(&gorm.Session{PrepareStmt: true})
		c.Set(constant.Session, tx)
		c.Set(constant.StatusCode, nil)
		c.Set(constant.Output, nil)
		c.Set(constant.Error, nil)
		c.Next()
		statusCode := c.GetInt(constant.StatusCode)
		err := c.MustGet(constant.Error)
		output := c.MustGet(constant.Output)
		if err != nil {
			if validationErr, ok := err.(validator.ValidationErrors); ok {
				//Translate : 用途 將錯誤訊息依造定義的語言翻譯
				sendResponse(c, statusCode, map[string]interface{}{"error": validationErr.Translate(vaildate.BindingTrans)})
			} else {
				// 直接顯示錯誤訊息 不進行翻譯動作
				sendResponse(c, statusCode, map[string]interface{}{"error": err.(error).Error()})
			}
		} else {
			sendResponse(c, statusCode, output)
		}
	}
}
