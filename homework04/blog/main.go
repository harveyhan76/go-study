package main

import (
	"blog/config"
	"blog/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 加载环境变量
	config.LoadEnv()

	// 初始化日志
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	// 初始化数据库
	db, err := config.InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// 创建Gin路由
	r := gin.Default()

	// 注册路由
	routes.SetupRoutes(r, db)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting on port " + port)
	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
