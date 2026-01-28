package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"home-vue-go/internal/database"

	"github.com/gin-gonic/gin"
)

// GetStats 获取统计数据
func GetStats(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		
		// 获取站点数量
		sites, _ := db.Client.Site.Query().All(ctx)
		totalSites := len(sites)
		
		// 获取总访问量（从内存记录）
		records := getVisitRecords()
		totalViews := len(records)
		
		// 获取独立访客数（去重IP）
		uniqueIPs := make(map[string]bool)
		for _, r := range records {
			uniqueIPs[r.IP] = true
		}
		uniqueVisitors := len(uniqueIPs)
		
		// 获取今日访问量
		today := time.Now()
		todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
		todayViews := 0
		for _, r := range records {
			if r.VisitTime.After(todayStart) {
				todayViews++
			}
		}
		
		stats := gin.H{
			"totalViews":      totalViews,
			"uniqueVisitors": uniqueVisitors,
			"todayViews":      todayViews,
			"totalSites":      totalSites,
		}

		c.JSON(http.StatusOK, stats)
	}
}

// GetChartData 获取图表数据
func GetChartData(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		period := c.DefaultQuery("period", "7")

		// 根据period确定天数
		periodDays := 7
		if period == "30" {
			periodDays = 30
		} else if period == "90" {
			periodDays = 90
		}
		
		// 获取真实趋势数据（从内存记录）
		records := getVisitRecords()
		trend := []gin.H{}
		now := time.Now()
		
		for i := 0; i < periodDays; i++ {
			date := now.AddDate(0, 0, -(periodDays-1-i))
			dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
			dayEnd := dayStart.Add(24 * time.Hour)
			
			// 统计当天的访问量
			count := 0
			for _, r := range records {
				if r.VisitTime.After(dayStart) && r.VisitTime.Before(dayEnd) {
					count++
				}
			}
			
			trend = append(trend, gin.H{
				"label": date.Format("1/2"),
				"value": count,
			})
		}
		
		// 获取访问来源数据（基于referer）
		sourceMap := make(map[string]int)
		total := len(records)
		
		for _, r := range records {
			ref := strings.ToLower(r.Referer)
			if ref == "" {
				sourceMap["直接访问"]++
			} else if strings.Contains(ref, "google") || strings.Contains(ref, "baidu") || strings.Contains(ref, "bing") || strings.Contains(ref, "yahoo") || strings.Contains(ref, "sogou") {
				sourceMap["搜索引擎"]++
			} else if strings.Contains(ref, "twitter") || strings.Contains(ref, "facebook") || strings.Contains(ref, "weibo") || strings.Contains(ref, "wechat") || strings.Contains(ref, "qq") {
				sourceMap["社交媒体"]++
			} else {
				sourceMap["其他"]++
			}
		}
		
		sources := []gin.H{}
		if total > 0 {
			for label, count := range sourceMap {
				sources = append(sources, gin.H{
					"label": label,
					"value": (count * 100) / total, // 转换为百分比
				})
			}
		} else {
			// 如果没有数据，返回默认值
			sources = []gin.H{
				{"label": "直接访问", "value": 100},
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"trend":   trend,
			"sources": sources,
		})
	}
}


// GetRecentVisits 获取最近访问记录
func GetRecentVisits(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 支持通过 query 参数自定义条数，默认 5，最大 50
		limit := 5
		if q := c.DefaultQuery("limit", "5"); q != "" {
			if n, err := parseInt(q); err == nil && n > 0 {
				limit = n
				if limit > 50 {
					limit = 50
				}
			}
		}

		// 从内存记录获取最近访问
		records := getVisitRecords()
		result := make([]gin.H, 0, limit)
		now := time.Now()

		// 取最近 limit 条
		start := len(records) - limit
		if start < 0 {
			start = 0
		}

		for i := len(records) - 1; i >= start && i >= 0; i-- {
			r := records[i]

			// 计算相对时间
			diff := now.Sub(r.VisitTime)
			var timeStr string
			if diff < time.Minute {
				timeStr = "刚刚"
			} else if diff < time.Hour {
				timeStr = fmt.Sprintf("%.0f分钟前", diff.Minutes())
			} else if diff < 24*time.Hour {
				timeStr = fmt.Sprintf("%.0f小时前", diff.Hours())
			} else {
				timeStr = fmt.Sprintf("%.0f天前", diff.Hours()/24)
			}

			result = append(result, gin.H{
				"path": r.Path,
				"ip":   r.IP,
				"time": timeStr,
			})
		}

		c.JSON(http.StatusOK, result)
	}
}

// NotifyConfigUpdate 通知配置更新（用于热重载）
func NotifyConfigUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以添加通知逻辑，比如通过WebSocket或SSE通知前端
		// 目前简单返回成功
		c.JSON(http.StatusOK, gin.H{
			"message": "配置更新通知已发送",
		})
	}
}
