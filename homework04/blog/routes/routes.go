package routes

import (
	"blog/handlers"
	"blog/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 初始化处理器
	authHandler := handlers.NewAuthHandler(db)
	postHandler := handlers.NewPostHandler(db)
	commentHandler := handlers.NewCommentHandler(db)

	// 认证路由
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
	}

	// 文章路由
	posts := r.Group("/api/posts")
	{
		posts.GET("", postHandler.GetPosts)
		posts.GET("/:id", postHandler.GetPost)
		posts.POST("", middleware.AuthMiddleware(), postHandler.CreatePost)
		posts.PUT("/:id", middleware.AuthMiddleware(), postHandler.UpdatePost)
		posts.DELETE("/:id", middleware.AuthMiddleware(), postHandler.DeletePost)
	}

	// 评论路由
	comments := r.Group("/api/posts/:id/comments")
	{
		comments.GET("", commentHandler.GetComments)
		comments.POST("", middleware.AuthMiddleware(), commentHandler.CreateComment)
		comments.DELETE("/:commentId", middleware.AuthMiddleware(), commentHandler.DeleteComment)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Blog API is running",
		})
	})
}
