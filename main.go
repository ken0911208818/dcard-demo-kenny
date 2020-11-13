package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ken0911208818/dcard-demo-kenny/handler"
	"github.com/ken0911208818/dcard-demo-kenny/lib/auth"
	"github.com/ken0911208818/dcard-demo-kenny/lib/config"
	"github.com/ken0911208818/dcard-demo-kenny/lib/middleware"
	"github.com/ken0911208818/dcard-demo-kenny/lib/vaildate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func init() {
	//postgres
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.GetStr("DB_HOST"),
		config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_NAME"),
		config.GetStr("DB_PORT"), config.GetStr("DB_SSL_MODE"))

	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		log.Panic("DB connection initialization failed:", err)
	}
	//redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.GetStr("REDIS_ENDPOINT"),
		Password: config.GetStr("REDIS_PASSWORD"),
		Network:  "tcp",
		PoolSize: config.GetInt("REDIS_POOL_SIZE"),
	})
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Panic("Redis connection initialization failed:", err)
	}
	// jwt setting
	secretKey := config.GetBytes("SECRET_KEY")
	tokenLifeTime := time.Duration(config.GetInt("TOKEN_LIFETIME")) * time.Minute

	auth.Init(secretKey, tokenLifeTime)
	middleware.Init(db, redisClient)
	vaildate.Init(config.GetStr("LOCALE"))
}
func main() {
	router := setupRouter()
	router.Run(":1234")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		authRouter := v1.Group("/auth/")
		{
			authRouter.POST("/", middleware.Plain(), handler.Login)
		}
		userRouter := v1.Group("/users/")
		{
			userRouter.POST("/", middleware.Plain(), handler.UserCreate)
		}
	}
	return r
}
