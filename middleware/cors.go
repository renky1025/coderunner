package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors CORS（跨域资源共享）中间件
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "lang", "Authorization", "accept", "origin", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "X-Requested-With", "Cache-Control", "x-tenant", "x-client", "X-CSRF-Token", "Content-Type", "access-control-allow-origin", "access-control-allow-headers", "Content-Length", "Accept-Encoding"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           86400,
	})
}
