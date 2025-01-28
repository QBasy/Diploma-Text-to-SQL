package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	middleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"time"
)

func RateLimiter() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  100,
	}

	store := memory.NewStore()

	instance := limiter.New(store, rate)

	return middleware.NewMiddleware(instance)
}
