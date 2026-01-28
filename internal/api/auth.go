package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"home-vue-go/internal/config"
	"home-vue-go/internal/database"
	"home-vue-go/internal/ent/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(db *database.Database, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
			return
		}

		ctx := c.Request.Context()
		user, err := db.Client.User.Query().Where(user.UsernameEQ(req.Username)).First(ctx)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}

		// 验证密码
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}

		// 获取登录IP和用户代理
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		// 查询IP地理位置
		location := queryIPLocation(ip)

		// 记录登录历史到内存
		addLoginHistory(user.Username, ip, location, userAgent, true)

		// 生成JWT token
		claims := Claims{
			Username: user.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
	}
}

func ChangePassword(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			OldPassword string `json:"oldPassword" binding:"required"`
			NewPassword string `json:"newPassword" binding:"required,min=8"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "旧密码和新密码不能为空，且新密码至少8位"})
			return
		}

		// 从JWT中获取用户名
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
			return
		}

		usernameStr := username.(string)
		ctx := c.Request.Context()

		// 查询用户
		user, err := db.Client.User.Query().Where(user.UsernameEQ(usernameStr)).First(ctx)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			return
		}

		// 验证旧密码
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "旧密码错误"})
			return
		}

		// 加密新密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}

		// 更新密码到数据库
		updatedUser, err := db.Client.User.UpdateOneID(user.ID).SetPassword(string(hashedPassword)).Save(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码更新失败: " + err.Error()})
			return
		}

		// 验证密码已保存（可选，用于调试）
		if updatedUser == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码更新失败: 未返回更新后的用户"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
	}
}

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证token"})
			c.Abort()
			return
		}

		// 移除 "Bearer " 前缀
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			// Token失效时，记录登录IP（用于统计）
			ip := c.ClientIP()
			userAgent := c.GetHeader("User-Agent")
			location := queryIPLocation(ip)
			// 记录为失败的登录尝试（token失效）
			addLoginHistory("", ip, location, userAgent, false)
			
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

// 查询IP地理位置（使用多个公共API）
func queryIPLocation(ip string) string {
	if ip == "" || ip == "::1" || ip == "127.0.0.1" {
		return "本地"
	}

	// 尝试多个IP查询服务
	apis := []struct {
		name string
		url  string
	}{
		{"ipapi", fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", ip)},
		{"ipapi.co", fmt.Sprintf("https://ipapi.co/%s/json/", ip)},
		{"ip.sb", fmt.Sprintf("https://api.ip.sb/geoip/%s", ip)},
	}

	for _, api := range apis {
		if location := queryIPFromAPI(api.url, api.name); location != "" {
			return location
		}
	}

	return "未知"
}

func queryIPFromAPI(url, apiName string) string {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return ""
	}

	switch apiName {
	case "ipapi":
		if country, ok := result["country"].(string); ok {
			region := ""
			if r, ok := result["regionName"].(string); ok {
				region = r
			}
			city := ""
			if c, ok := result["city"].(string); ok {
				city = c
			}
			parts := []string{country}
			if region != "" {
				parts = append(parts, region)
			}
			if city != "" {
				parts = append(parts, city)
			}
			return strings.Join(parts, " ")
		}
	case "ipapi.co":
		if country, ok := result["country_name"].(string); ok {
			city := ""
			if c, ok := result["city"].(string); ok {
				city = c
			}
			if city != "" {
				return fmt.Sprintf("%s %s", country, city)
			}
			return country
		}
	case "ip.sb":
		if country, ok := result["country"].(string); ok {
			city := ""
			if c, ok := result["city"].(string); ok {
				city = c
			}
			if city != "" {
				return fmt.Sprintf("%s %s", country, city)
			}
			return country
		}
	}

	return ""
}

// 获取登录历史
func GetLoginHistory(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := c.DefaultQuery("limit", "20")
		limitInt := 20
		if n, err := parseInt(limit); err == nil && n > 0 {
			limitInt = n
			if limitInt > 100 {
				limitInt = 100
			}
		}

		// 从内存获取登录历史
		histories := getLoginHistoryRecords(limitInt)
		
		result := make([]gin.H, 0, len(histories))
		for _, h := range histories {
			result = append(result, gin.H{
				"username":  h.Username,
				"ip":        h.IP,
				"location":  h.Location,
				"userAgent": h.UserAgent,
				"loginTime": h.LoginTime.Format("2006-01-02 15:04:05"),
				"success":   h.Success,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  result,
			"count": len(result),
		})
	}
}

// 简单的整数解析
func parseInt(s string) (int, error) {
	var n int
	for _, char := range s {
		if char >= '0' && char <= '9' {
			n = n*10 + int(char-'0')
		} else {
			return 0, fmt.Errorf("invalid number")
		}
	}
	return n, nil
}
