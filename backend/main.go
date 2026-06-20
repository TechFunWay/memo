package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"memo/config"
	"memo/database"
	"memo/handlers"
	"memo/logger"
	"memo/middleware"
	"memo/services"
)

// 命令行参数
var (
	port      = flag.String("port", "8905", "服务端口")
	dataDir   = flag.String("data-dir", "./data", "数据目录路径")
	webDir = flag.String("web-dir", "./static", "静态资源目录")
	resetAdmin = flag.Bool("reset-admin", false, "重置管理员密码（需配合 -admin-pass 或 ADMIN_PASS 环境变量）")
	adminPass  = flag.String("admin-pass", "", "新密码（与 -reset-admin 配合使用）")
)

func main() {
	flag.Parse()

	// 初始化日志
	if err := logger.InitLogger(*dataDir); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	sugarLogger := logger.GetSugarLogger()
	sugarLogger.Infof("🚀 备忘录 %s", config.AppVersion)
	sugarLogger.Infof("📁 Data directory: %s", *dataDir)
	sugarLogger.Infof("🌐 Web directory: %s", *webDir)

	// 初始化配置
	config.InitConfig(*dataDir)

	// 初始化数据库
	dbPath := filepath.Join(*dataDir, "db", "database.db")
	if err := database.InitDB(dbPath); err != nil {
		sugarLogger.Fatalf("❌ Failed to initialize database: %v", err)
	}
	sugarLogger.Infof("📄 Database path: %s", dbPath)

	// 初始化数据库表
	if err := database.InitTables(); err != nil {
		sugarLogger.Fatalf("❌ Failed to initialize tables: %v", err)
	}

	// 执行数据库升级
	sugarLogger.Info("🔄 Checking database upgrades...")
	upgradeSvc := services.NewUpgradeService(database.DB, config.AppVersion, *dataDir)
	if err := upgradeSvc.RunUpgrades(); err != nil {
		sugarLogger.Errorf("❌ Database upgrade failed: %v", err)
	}

	// 命令行重置管理员密码
	if *resetAdmin {
		pass := *adminPass
		if pass == "" {
			pass = os.Getenv("ADMIN_PASS")
		}
		if pass == "" {
			log.Fatal("❌ 请通过 -admin-pass 参数或 ADMIN_PASS 环境变量指定新密码")
		}
		authService := &services.AuthService{}
		// 前端密码已做 MD5 传输，但命令行输入的是明文，需要做一次 MD5
		hashedPass := services.MD5Hash(pass)
		if err := authService.AdminResetPassword(hashedPass); err != nil {
			log.Fatalf("❌ 重置管理员密码失败: %v", err)
		}
		sugarLogger.Info("✅ 管理员密码重置成功")
		os.Exit(0)
	}

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 上传文件静态服务
	r.Static("/uploads", config.UploadDir)

	// 静态文件服务 (在 NoRoute 中处理，避免与 /api/* 路由冲突)
	if _, err := os.Stat(*webDir); err == nil {
		// 使用外部 static 目录
		sugarLogger.Infof("✅ 使用外部 static 目录：%s", *webDir)
	} else {
		// 使用嵌入的 static 目录
		sugarLogger.Infof("✅ 使用嵌入的 static 目录")
	}
	sugarLogger.Infof("✅ 上传文件目录：%s", config.UploadDir)

	// API 路由
	api := r.Group("/api")
	{
		// 版本信息
		api.GET("/version", handlers.GetVersion)

		// 认证相关 (不需要登录)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.GET("/check", handlers.CheckAuth)
			auth.GET("/has-users", handlers.CheckHasUsers)
			auth.GET("/security-question", handlers.GetSecurityQuestion)
			auth.POST("/verify-security", handlers.VerifySecurityAnswer)
		}

		// 公开配置
		api.GET("/config/public", handlers.GetPublicConfig)
		api.GET("/config/system", handlers.GetSystemConfig)

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			authorized.GET("/auth/me", handlers.GetCurrentUser)
			authorized.PUT("/auth/password", handlers.ChangePassword)
			authorized.PUT("/auth/security-question", handlers.SetSecurityQuestion)
			authorized.GET("/auth/security-question/check", handlers.CheckSecurityQuestion)

			// 备忘录
			memos := authorized.Group("/memos")
			{
				memos.GET("", handlers.GetMemos)
				memos.POST("", handlers.CreateMemo)
				memos.GET("/:id", handlers.GetMemo)
				memos.PUT("/:id", handlers.UpdateMemo)
				memos.DELETE("/:id", handlers.DeleteMemo)
			}

			// 文件夹
			folders := authorized.Group("/folders")
			{
				folders.GET("", handlers.GetFolders)
				folders.POST("", handlers.CreateFolder)
				folders.PUT("/:id", handlers.UpdateFolder)
				folders.DELETE("/:id", handlers.DeleteFolder)
			}

			// 上传
			upload := authorized.Group("/upload")
			{
				upload.POST("/image", handlers.UploadImage)
			}

			// 用户管理 (仅管理员)
			users := authorized.Group("/users")
			users.Use(middleware.AdminMiddleware())
			{
				users.GET("", handlers.GetAllUsers)
				users.POST("", handlers.CreateUser)
				users.PUT("/:id", handlers.UpdateUser)
				users.POST("/:id/toggle", handlers.ToggleUserStatus)
				users.POST("/:id/reset-password", handlers.AdminResetUserPassword)
				users.DELETE("/:id", handlers.DeleteUser)
			}

			// 系统配置 (仅管理员)
			configs := authorized.Group("/configs")
			configs.Use(middleware.AdminMiddleware())
			{
				configs.GET("", handlers.GetConfigs)
				configs.POST("", handlers.UpdateConfig)
			}
		}
	}

	// SPA 路由：先尝试提供静态文件，找不到则返回 index.html
	r.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}
		// 尝试提供静态文件
		staticPath := filepath.Join(*webDir, c.Request.URL.Path)
		if _, err := os.Stat(staticPath); err == nil {
			c.File(staticPath)
			return
		}
		// SPA 路由回退
		c.File(filepath.Join(*webDir, "index.html"))
	})

	// 启动服务器
	addr := ":" + *port
	sugarLogger.Infof("🌐 Server starting on port %s", addr)
	sugarLogger.Infof("📌 API endpoint: http://localhost:%s/api", addr)

	if err := r.Run(addr); err != nil {
		sugarLogger.Fatalf("❌ Server failed to start: %v", err)
	}
}
