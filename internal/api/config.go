package api

import (
	"net/http"

	"home-vue-go/internal/config"
	"home-vue-go/internal/database"

	"github.com/gin-gonic/gin"
)

// GetSiteConfig 获取站点配置（含底部年份配置）
func GetSiteConfig(db *database.Database, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		footerCfg, _ := cfg.LoadFooterYear()
		visitTimerCfg, _ := cfg.LoadVisitTimer()

		siteCfg, err := db.Client.SiteConfig.Get(ctx, 1)
		if err != nil {
			// 如果不存在，返回默认值 + 年份配置
		c.JSON(http.StatusOK, gin.H{
			"siteName":        "个人主页",
			"siteURL":         "https://example.com",
			"siteIcon":        "/favicon.ico",
			"siteDescription": "一个基于Vue3的个人主页",
			"siteKeywords":    "个人主页,Vue3",
			"userName":        "用户",
			"profileImageURL": "",
			"icpNumber":       "暂未填写",
			"policeNumber":    "暂未填写",
			"pageTitle":       "个人主页",
			"favicon":         "/favicon.ico",
			"umamiScript":     "",
			"umamiWebsiteId":  "",
			"iconLibrary":     "//lib.baomitu.com/font-awesome/6.5.0/css/all.min.css",
			"fontLibrary":     "",
			"showVisitTimer":  visitTimerCfg.ShowVisitTimer,
			"footerYearStart": footerCfg.Start,
			"footerYearEnd":   footerCfg.End,
		})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"siteName":        siteCfg.SiteName,
			"siteURL":         siteCfg.SiteURL,
			"siteIcon":        siteCfg.SiteIcon,
			"siteDescription": siteCfg.SiteDescription,
			"siteKeywords":    siteCfg.SiteKeywords,
			"userName":        siteCfg.UserName,
			"profileImageURL": siteCfg.ProfileImageURL,
			"icpNumber":       siteCfg.IcpNumber,
			"policeNumber":    siteCfg.PoliceNumber,
			"pageTitle":       siteCfg.PageTitle,
			"favicon":         siteCfg.Favicon,
			"umamiScript":     siteCfg.UmamiScript,
			"umamiWebsiteId":  siteCfg.UmamiWebsiteID,
			"iconLibrary":     siteCfg.IconLibrary,
			"fontLibrary":     siteCfg.FontLibrary,
			"showVisitTimer":  visitTimerCfg.ShowVisitTimer,
			"footerYearStart": footerCfg.Start,
			"footerYearEnd":   footerCfg.End,
		})
	}
}

// UpdateSiteConfig 更新站点配置（含底部年份配置）
func UpdateSiteConfig(db *database.Database, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			SiteName        string `json:"siteName"`
			SiteURL         string `json:"siteURL"`
			SiteIcon        string `json:"siteIcon"`
			SiteDescription string `json:"siteDescription"`
			SiteKeywords    string `json:"siteKeywords"`
			UserName        string `json:"userName"`
			ProfileImageURL string `json:"profileImageURL"`
			ICPNumber       string `json:"icpNumber"`
			PoliceNumber    string `json:"policeNumber"`
			PageTitle       string `json:"pageTitle"`
			Favicon         string `json:"favicon"`
			UmamiScript     string `json:"umamiScript"`
			UmamiWebsiteId  string `json:"umamiWebsiteId"`
			IconLibrary     string `json:"iconLibrary"`
			FontLibrary     string `json:"fontLibrary"`
			ShowVisitTimer  *bool  `json:"showVisitTimer"`

			FooterYearStart string `json:"footerYearStart"`
			FooterYearEnd   string `json:"footerYearEnd"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		update := db.Client.SiteConfig.UpdateOneID(1)

		if req.SiteName != "" {
			update.SetSiteName(req.SiteName)
		}
		if req.SiteURL != "" {
			update.SetSiteURL(req.SiteURL)
		}
		if req.SiteIcon != "" {
			update.SetSiteIcon(req.SiteIcon)
		}
		if req.SiteDescription != "" {
			update.SetSiteDescription(req.SiteDescription)
		}
		if req.SiteKeywords != "" {
			update.SetSiteKeywords(req.SiteKeywords)
		}
		if req.UserName != "" {
			update.SetUserName(req.UserName)
		}
		if req.ProfileImageURL != "" {
			update.SetProfileImageURL(req.ProfileImageURL)
		}
		if req.ICPNumber != "" {
			update.SetIcpNumber(req.ICPNumber)
		}
		if req.PoliceNumber != "" {
			update.SetPoliceNumber(req.PoliceNumber)
		}
		if req.PageTitle != "" {
			update.SetPageTitle(req.PageTitle)
		}
		if req.Favicon != "" {
			update.SetFavicon(req.Favicon)
		}
		if req.UmamiScript != "" {
			update.SetUmamiScript(req.UmamiScript)
		}
		if req.UmamiWebsiteId != "" {
			update.SetUmamiWebsiteID(req.UmamiWebsiteId)
		}
		if req.IconLibrary != "" {
			update.SetIconLibrary(req.IconLibrary)
		}
		if req.FontLibrary != "" {
			update.SetFontLibrary(req.FontLibrary)
		}

		siteCfg, err := update.Save(ctx)
		if err != nil {
			// 如果不存在，创建新的
			create := db.Client.SiteConfig.Create().
				SetSiteName(req.SiteName).
				SetSiteURL(req.SiteURL).
				SetSiteIcon(req.SiteIcon).
				SetSiteDescription(req.SiteDescription).
				SetSiteKeywords(req.SiteKeywords).
				SetUserName(req.UserName)
			if req.ProfileImageURL != "" {
				create.SetProfileImageURL(req.ProfileImageURL)
			}
			if req.ICPNumber != "" {
				create.SetIcpNumber(req.ICPNumber)
			}
			if req.PoliceNumber != "" {
				create.SetPoliceNumber(req.PoliceNumber)
			}
			if req.PageTitle != "" {
				create.SetPageTitle(req.PageTitle)
			}
			if req.Favicon != "" {
				create.SetFavicon(req.Favicon)
			}
			if req.UmamiScript != "" {
				create.SetUmamiScript(req.UmamiScript)
			}
			if req.UmamiWebsiteId != "" {
				create.SetUmamiWebsiteID(req.UmamiWebsiteId)
			}
			if req.IconLibrary != "" {
				create.SetIconLibrary(req.IconLibrary)
			}
			if req.FontLibrary != "" {
				create.SetFontLibrary(req.FontLibrary)
			}
			siteCfg, err = create.Save(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// 保存底部年份配置到独立文件
		_ = cfg.SaveFooterYear(&config.FooterYearConfig{
			Start: req.FooterYearStart,
			End:   req.FooterYearEnd,
		})

		// 保存访问时间显示配置
		if req.ShowVisitTimer != nil {
			_ = cfg.SaveVisitTimer(&config.VisitTimerConfig{
				ShowVisitTimer: *req.ShowVisitTimer,
			})
		}

		// 读取最新的配置用于返回
		footerCfg, _ := cfg.LoadFooterYear()
		visitTimerCfg, _ := cfg.LoadVisitTimer()

		c.JSON(http.StatusOK, gin.H{
			"siteName":        siteCfg.SiteName,
			"siteURL":         siteCfg.SiteURL,
			"siteIcon":        siteCfg.SiteIcon,
			"siteDescription": siteCfg.SiteDescription,
			"siteKeywords":    siteCfg.SiteKeywords,
			"userName":        siteCfg.UserName,
			"profileImageURL": siteCfg.ProfileImageURL,
			"icpNumber":       siteCfg.IcpNumber,
			"policeNumber":    siteCfg.PoliceNumber,
			"pageTitle":       siteCfg.PageTitle,
			"favicon":         siteCfg.Favicon,
			"umamiScript":     siteCfg.UmamiScript,
			"umamiWebsiteId":  siteCfg.UmamiWebsiteID,
			"iconLibrary":     siteCfg.IconLibrary,
			"fontLibrary":     siteCfg.FontLibrary,
			"showVisitTimer":  visitTimerCfg.ShowVisitTimer,
			"footerYearStart": footerCfg.Start,
			"footerYearEnd":   footerCfg.End,
		})
	}
}

// GetRotatingTexts 获取轮换文本配置
func GetRotatingTexts(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		textsCfg, err := cfg.LoadRotatingTexts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "加载轮换文本配置失败: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"texts": textsCfg.Texts,
		})
	}
}

// UpdateRotatingTexts 更新轮换文本配置
func UpdateRotatingTexts(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Texts []string `json:"texts"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 限制最多8条文本
		if len(req.Texts) > 8 {
			req.Texts = req.Texts[:8]
		}

		textsCfg := &config.RotatingTextsConfig{
			Texts: req.Texts,
		}

		if err := cfg.SaveRotatingTexts(textsCfg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存轮换文本配置失败: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "轮换文本配置已保存",
			"texts":   textsCfg.Texts,
		})
	}
}

