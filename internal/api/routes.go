package api

import (
	"home-vue-go/internal/config"
	"home-vue-go/internal/database"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *database.Database, cfg *config.Config) {
	// 初始化日志缓冲区
	InitLogBuffer()
	
	// 添加日志中间件
	r.Use(LoggingMiddleware())
	
	api := r.Group("/api")
	{
		// 公开API - 获取数据
		api.GET("/sites", GetSites(db))
		api.GET("/contacts", GetContacts(db))
		api.GET("/config", GetSiteConfig(db, cfg))
		api.GET("/frontend-config", GetFrontendConfig(db))
		api.GET("/rotating-texts", GetRotatingTexts(cfg))
		
		// 访问统计API（公开，用于记录访问）
		api.POST("/track-visit", TrackVisit(db))

		// 认证API
		auth := api.Group("/auth")
		{
			auth.POST("/login", Login(db, cfg))
		}

		// 受保护的管理API
		admin := api.Group("/admin")
		admin.Use(JWTAuthMiddleware(cfg.JWTSecret))
		{
			// 站点管理
			admin.GET("/sites", GetSites(db))
			admin.POST("/sites", CreateSite(db))
			admin.PUT("/sites/:id", UpdateSite(db))
			admin.DELETE("/sites/:id", DeleteSite(db))

			// 联系方式管理
			admin.GET("/contacts", GetContacts(db))
			admin.POST("/contacts", CreateContact(db))
			admin.PUT("/contacts/:id", UpdateContact(db))
			admin.DELETE("/contacts/:id", DeleteContact(db))

			// 站点配置管理
			admin.GET("/config", GetSiteConfig(db, cfg))
			admin.PUT("/config", UpdateSiteConfig(db, cfg))
			
			// 轮换文本配置
			admin.GET("/rotating-texts", GetRotatingTexts(cfg))
			admin.PUT("/rotating-texts", UpdateRotatingTexts(cfg))

			// 文件上传
			admin.POST("/upload", UploadFile(cfg))

			// 统计API
			admin.GET("/stats", GetStats(db))
			admin.GET("/charts", GetChartData(db))
			admin.GET("/recent-visits", GetRecentVisits(db))
			admin.POST("/notify-update", NotifyConfigUpdate())

			// 用户管理
			admin.PUT("/change-password", ChangePassword(db))
			
			// 日志API
			admin.GET("/logs", GetBackendLogs(cfg))
			
			// 登录历史API
			admin.GET("/login-history", GetLoginHistory(db))
		}
	}

	// 静态文件服务 - 提供上传的图片
	r.Static("/uploads", cfg.UploadDir)
}
