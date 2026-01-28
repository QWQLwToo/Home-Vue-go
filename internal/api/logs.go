package api

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"home-vue-go/internal/config"

	"github.com/gin-gonic/gin"
)

var logBuffer []string
var maxLogLines = 1000

// 初始化日志缓冲区
func InitLogBuffer() {
	logBuffer = make([]string, 0, maxLogLines)
}

// 添加日志到缓冲区
func AddLogToBuffer(logLine string) {
	logBuffer = append(logBuffer, logLine)
	if len(logBuffer) > maxLogLines {
		logBuffer = logBuffer[len(logBuffer)-maxLogLines:]
	}
}

// 获取后端日志
func GetBackendLogs(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		lines := c.DefaultQuery("lines", "100")
		linesInt := 100
		if lines == "all" {
			linesInt = maxLogLines
		} else {
			// 尝试解析lines参数
			if n, err := parseIntLogs(lines); err == nil && n > 0 {
				linesInt = n
				if linesInt > maxLogLines {
					linesInt = maxLogLines
				}
			}
		}

		// 从缓冲区获取日志
		logs := make([]string, 0)
		if len(logBuffer) > linesInt {
			logs = logBuffer[len(logBuffer)-linesInt:]
		} else {
			logs = logBuffer
		}

		// 尝试从日志文件读取（如果存在）
		logFile := filepath.Join(cfg.DataDir, "app.log")
		if fileLogs, err := readLogFile(logFile, linesInt); err == nil && len(fileLogs) > 0 {
			// 合并文件日志和缓冲区日志
			allLogs := append(fileLogs, logs...)
			// 去重并排序
			logs = deduplicateLogs(allLogs)
			if len(logs) > linesInt {
				logs = logs[len(logs)-linesInt:]
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"logs":  logs,
			"count": len(logs),
		})
	}
}

// 读取日志文件
func readLogFile(filePath string, maxLines int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) >= maxLines*2 {
			// 只保留最后maxLines行
			lines = lines[len(lines)-maxLines:]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 返回最后maxLines行
	if len(lines) > maxLines {
		return lines[len(lines)-maxLines:], nil
	}
	return lines, nil
}

// 去重日志
func deduplicateLogs(logs []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for i := len(logs) - 1; i >= 0; i-- {
		if !seen[logs[i]] {
			seen[logs[i]] = true
			result = append([]string{logs[i]}, result...)
		}
	}
	return result
}

// 简单的整数解析（logs.go专用）
func parseIntLogs(s string) (int, error) {
	var n int
	for _, char := range s {
		if char >= '0' && char <= '9' {
			n = n*10 + int(char-'0')
		} else {
			return 0, io.EOF
		}
	}
	return n, nil
}

// 日志中间件 - 记录请求日志
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		ip := c.ClientIP()

		c.Next()

		latency := time.Since(start)
		method := c.Request.Method
		statusCode := c.Writer.Status()

		logLine := formatLogLine(start, method, statusCode, latency, path, query, ip)
		AddLogToBuffer(logLine)
	}
}

func formatLogLine(timestamp time.Time, method string, statusCode int, latency time.Duration, path, query, ip string) string {
	if query != "" {
		path = path + "?" + query
	}
	// 使用Asia/Shanghai时区格式化时间
	loc, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := timestamp.In(loc)
	// 格式化日志：只包含一个时间戳
	return fmt.Sprintf("[%s] %s %s :: %s %d %s",
		shanghaiTime.Format("2006-01-02 15:04:05"),
		method,
		path,
		ip,
		statusCode,
		latency.String(),
	)
}
