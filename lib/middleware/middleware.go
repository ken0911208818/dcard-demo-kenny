package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/ken0911208818/dcard-demo-kenny/lib/constant"
	"github.com/ken0911208818/dcard-demo-kenny/lib/vaildate"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	redisClient *redis.Client
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
