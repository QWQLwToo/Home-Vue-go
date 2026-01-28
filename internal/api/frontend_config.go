package api

import (
	"net/http"

	"home-vue-go/internal/database"

	"github.com/gin-gonic/gin"
)

// GetFrontendConfig 获取前端配置（用于index.html等）
func GetFrontendConfig(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		config, err := db.Client.SiteConfig.Get(ctx, 1)
		if err != nil {
			// 如果不存在，返回默认值
			c.JSON(http.StatusOK, gin.H{
				"title":          "个人主页",
				"keywords":       "个人主页,Vue3",
				"description":    "一个基于Vue3的个人主页",
				"favicon":        "/favicon.ico",
				"umamiScript":   "",
				"umamiWebsiteId": "",
				"iconLibrary":   "//lib.baomitu.com/font-awesome/6.5.0/css/all.min.css",
				"fontLibrary":    "",
			})
			return
		}

		pageTitle := config.PageTitle
		if pageTitle == "" {
			pageTitle = "个人主页"
		}
		favicon := config.Favicon
		if favicon == "" {
			favicon = "/favicon.ico"
		}
		iconLibrary := config.IconLibrary
		if iconLibrary == "" {
			iconLibrary = "//lib.baomitu.com/font-awesome/6.5.0/css/all.min.css"
		}

		// 确保返回完整的配置信息，包括站点名称和URL
		c.JSON(http.StatusOK, gin.H{
			"title":          pageTitle,
			"siteName":       config.SiteName,
			"siteURL":        config.SiteURL,
			"keywords":       config.SiteKeywords,
			"description":    config.SiteDescription,
			"favicon":        favicon,
			"umamiScript":    config.UmamiScript,
			"umamiWebsiteId":  config.UmamiWebsiteID,
			"iconLibrary":    iconLibrary,
			"fontLibrary":    config.FontLibrary,
		})
	}
}
