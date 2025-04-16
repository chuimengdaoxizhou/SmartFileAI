package routes

import (
	"RAG/controllers"
	"RAG/middleware"
	"RAG/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Initialize Redis
	services.InitRedis()

	router := gin.Default()
	// 正确设置 CORS 头
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, multipart/form-data")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public routes
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST("/upload", controllers.ResumableUpload)
	router.GET("/check-chunk", controllers.CheckChunk)
	router.POST("upload2", controllers.UploadFileHandler)
	// Protected routes
	protected := router.Group("/").Use(middleware.AuthMiddleware())
	{
		protected.POST("/logout", controllers.Logout)
		protected.DELETE("/account", controllers.DeleteAccount)
		protected.GET("/online-users", controllers.GetOnlineUsers)
		// 上传相关路由
	}

	return router
}
