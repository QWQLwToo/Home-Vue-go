package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"home-vue-go/internal/config"

	"github.com/gin-gonic/gin"
)

var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".avif": true,
	".svg":  true,
	".bmp":  true,
	".ico":  true,
}

func UploadFile(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败: " + err.Error()})
			return
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedImageExtensions[ext] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("不支持的文件格式。支持的格式: %v", getKeys(allowedImageExtensions)),
			})
			return
		}

		// 生成唯一文件名
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dst := filepath.Join(cfg.UploadDir, filename)

		// 保存文件
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败: " + err.Error()})
			return
		}

		// 返回文件URL（相对于uploads目录）
		fileURL := fmt.Sprintf("/uploads/%s", filename)
		c.JSON(http.StatusOK, gin.H{
			"url":  fileURL,
			"path": fileURL,
		})
	}
}

func getKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
