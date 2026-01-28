package api

import (
	"net/http"
	"sync"
	"time"

	"home-vue-go/internal/database"

	"github.com/gin-gonic/gin"
)

// VisitRecord 表示一次访问记录（暂存在内存中）
type VisitRecord struct {
	Path      string
	IP        string
	UserAgent string
	Referer   string
	VisitTime time.Time
}

// 简单的内存存储与读写锁
var (
	visitMutex   sync.RWMutex
	visitRecords []VisitRecord
	maxVisits    = 1000 // 最多保留的访问记录条数
)

// TrackVisit 记录访问
func TrackVisit(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Path    string `json:"path"`
			Referer string `json:"referer"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}

		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		// 记录访问到内存
		visitMutex.Lock()
		visitRecords = append(visitRecords, VisitRecord{
			Path:      req.Path,
			IP:        ip,
			UserAgent: userAgent,
			Referer:   req.Referer,
			VisitTime: time.Now(),
		})
		
		// 限制记录数量
		if len(visitRecords) > maxVisits {
			visitRecords = visitRecords[len(visitRecords)-maxVisits:]
		}
		visitMutex.Unlock()
		
		// TODO: 等Ent代码生成后，改用数据库存储
		// ctx := c.Request.Context()
		// if db.Client.Visit != nil {
		// 	db.Client.Visit.Create()...
		// }

		c.JSON(http.StatusOK, gin.H{"message": "访问已记录"})
	}
}

// 获取访问记录（从内存）
func getVisitRecords() []VisitRecord {
	visitMutex.RLock()
	defer visitMutex.RUnlock()
	return visitRecords
}
