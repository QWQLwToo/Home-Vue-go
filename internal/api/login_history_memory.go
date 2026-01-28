package api

import (
	"sync"
	"time"
)

// 内存存储登录历史（临时方案，等Ent代码生成后改用数据库）
type LoginHistoryRecord struct {
	Username  string
	IP        string
	Location  string
	UserAgent string
	LoginTime time.Time
	Success   bool
}

var (
	loginHistoryRecords []LoginHistoryRecord
	loginHistoryMutex   sync.RWMutex
	maxLoginHistory     = 1000 // 最多保存1000条登录记录
)

// 添加登录历史记录
func addLoginHistory(username, ip, location, userAgent string, success bool) {
	loginHistoryMutex.Lock()
	defer loginHistoryMutex.Unlock()
	
	loginHistoryRecords = append(loginHistoryRecords, LoginHistoryRecord{
		Username:  username,
		IP:        ip,
		Location:  location,
		UserAgent: userAgent,
		LoginTime: time.Now(),
		Success:   success,
	})
	
	// 限制记录数量
	if len(loginHistoryRecords) > maxLoginHistory {
		loginHistoryRecords = loginHistoryRecords[len(loginHistoryRecords)-maxLoginHistory:]
	}
}

// 获取登录历史记录
func getLoginHistoryRecords(limit int) []LoginHistoryRecord {
	loginHistoryMutex.RLock()
	defer loginHistoryMutex.RUnlock()
	
	if limit <= 0 || limit > len(loginHistoryRecords) {
		limit = len(loginHistoryRecords)
	}
	
	// 返回最后limit条记录（最新的）
	start := len(loginHistoryRecords) - limit
	if start < 0 {
		start = 0
	}
	
	result := make([]LoginHistoryRecord, 0, limit)
	for i := len(loginHistoryRecords) - 1; i >= start; i-- {
		result = append(result, loginHistoryRecords[i])
	}
	
	// 反转顺序，使最新的在前
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	
	return result
}
