package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"home-vue-go/internal/api"
	"home-vue-go/internal/config"
	"home-vue-go/internal/database"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var distFS embed.FS

func main() {
	// 设置亚洲/上海时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatal("无法加载时区:", err)
	}
	time.Local = loc

	// 获取可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("无法获取可执行文件路径:", err)
	}
	exeDir := filepath.Dir(exePath)

	// 创建data目录
	dataDir := filepath.Join(exeDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("无法创建data目录:", err)
	}

	// 初始化配置
	cfg := config.New(dataDir)

	// 初始化数据库
	db, err := database.Init(cfg.DatabasePath)
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}
	defer db.Close()

	// 设置Gin模式
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin路由（不使用Default以避免重复日志）
	r := gin.New()

	// 添加恢复中间件
	r.Use(gin.Recovery())

	// 静态文件服务 - 从文件系统提供前端构建文件
	// 优先使用可执行文件同目录下的dist目录
	distPath := filepath.Join(exeDir, "dist")
	if _, err := os.Stat(distPath); err == nil {
		// 使用可执行文件同目录下的dist
		r.Static("/static", filepath.Join(distPath, "static"))
		r.StaticFile("/favicon.ico", filepath.Join(distPath, "favicon.ico"))
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// API和上传路径不处理
			if path == "/api" || strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/uploads/") {
				c.Status(404)
				return
			}
			// 其他路径返回index.html（SPA路由支持）
			c.File(filepath.Join(distPath, "index.html"))
		})
		log.Printf("前端文件目录: %s", distPath)
	} else {
		// 回退到当前工作目录的dist（开发模式）
		if _, err := os.Stat("./dist"); err == nil {
			r.Static("/static", "./dist/static")
			r.StaticFile("/favicon.ico", "./dist/favicon.ico")
			r.NoRoute(func(c *gin.Context) {
				if c.Request.URL.Path != "/api" && !strings.HasPrefix(c.Request.URL.Path, "/api/") && !strings.HasPrefix(c.Request.URL.Path, "/uploads/") {
					c.File("./dist/index.html")
				}
			})
			log.Println("使用当前目录的dist文件夹（开发模式）")
		} else {
			log.Println("警告: 未找到前端文件目录，前端功能不可用")
			log.Println("提示: 请将前端构建文件放在可执行文件同目录的dist文件夹中")
		}
	}

	// 配置CORS
	r.Use(corsMiddleware())

	// 初始化API路由
	api.SetupRoutes(r, db, cfg)

	// 创建前端服务器（1552端口）- 从嵌入的文件系统提供前端文件
	frontendRouter := gin.New()
	frontendRouter.Use(gin.Recovery())
	frontendRouter.Use(corsMiddleware())

	// 从嵌入的文件系统加载前端文件
	distRoot, err := fs.Sub(distFS, "dist")
	if err == nil {
		// 使用嵌入的文件系统
		frontendRouter.StaticFS("/static", http.FS(distRoot))
		
		// 提供favicon
		frontendRouter.GET("/favicon.ico", func(c *gin.Context) {
			data, err := distRoot.Open("favicon.ico")
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			defer data.Close()
			content, err := io.ReadAll(data)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
			c.Data(http.StatusOK, "image/x-icon", content)
		})
		
		// API代理：将/api请求代理到1551端口
		frontendRouter.Any("/api/*path", func(c *gin.Context) {
			client := &http.Client{
				Timeout: 30 * time.Second,
			}
			
			targetURL := "http://localhost:1551" + c.Request.URL.Path
			if c.Request.URL.RawQuery != "" {
				targetURL += "?" + c.Request.URL.RawQuery
			}
			
			req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, targetURL, c.Request.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建代理请求失败"})
				return
			}
			
			for key, values := range c.Request.Header {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
			
			resp, err := client.Do(req)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": "代理请求失败: " + err.Error()})
				return
			}
			defer resp.Body.Close()
			
			for key, values := range resp.Header {
				for _, value := range values {
					c.Writer.Header().Add(key, value)
				}
			}
			
			// 复制响应状态码和内容
			c.Status(resp.StatusCode)
			c.Header("Content-Type", resp.Header.Get("Content-Type"))
			io.Copy(c.Writer, resp.Body)
		})
		
		// /uploads代理
		frontendRouter.Static("/uploads", cfg.UploadDir)
		
		// SPA路由支持
		frontendRouter.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/uploads/") {
				c.Status(404)
				return
			}
			
			// 尝试打开文件
			filePath := strings.TrimPrefix(path, "/")
			if filePath == "" {
				filePath = "index.html"
			}
			file, err := distRoot.Open(filePath)
			if err == nil {
				defer file.Close()
				stat, err := file.Stat()
				if err == nil && !stat.IsDir() {
					content, err := io.ReadAll(file)
					if err == nil {
						contentType := "text/html"
						if strings.HasSuffix(filePath, ".css") {
							contentType = "text/css"
						} else if strings.HasSuffix(filePath, ".js") {
							contentType = "application/javascript"
						} else if strings.HasSuffix(filePath, ".json") {
							contentType = "application/json"
						} else if strings.HasSuffix(filePath, ".ico") {
							contentType = "image/x-icon"
						} else if strings.HasSuffix(filePath, ".png") {
							contentType = "image/png"
						} else if strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg") {
							contentType = "image/jpeg"
						} else if strings.HasSuffix(filePath, ".svg") {
							contentType = "image/svg+xml"
						}
						c.Data(http.StatusOK, contentType, content)
						return
					}
				}
			}
			
			// 返回index.html（SPA路由）
			indexFile, err := distRoot.Open("index.html")
			if err == nil {
				defer indexFile.Close()
				content, err := io.ReadAll(indexFile)
				if err == nil {
					c.Data(http.StatusOK, "text/html; charset=utf-8", content)
				} else {
					c.Status(http.StatusNotFound)
				}
			} else {
				c.Status(http.StatusNotFound)
			}
		})
		
		log.Println("使用嵌入的前端文件（单一可执行文件模式）")
	} else {
		// 回退到文件系统（开发模式）
		log.Println("警告: 无法加载嵌入的前端文件，尝试从文件系统加载")
		distPath := filepath.Join(exeDir, "dist")
		if _, err := os.Stat(distPath); err == nil {
			frontendRouter.Static("/static", filepath.Join(distPath, "static"))
			frontendRouter.StaticFile("/favicon.ico", filepath.Join(distPath, "favicon.ico"))
			
			frontendRouter.Any("/api/*path", func(c *gin.Context) {
				client := &http.Client{Timeout: 30 * time.Second}
				targetURL := "http://localhost:1551" + c.Request.URL.Path
				if c.Request.URL.RawQuery != "" {
					targetURL += "?" + c.Request.URL.RawQuery
				}
				req, _ := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, targetURL, c.Request.Body)
				for key, values := range c.Request.Header {
					for _, value := range values {
						req.Header.Add(key, value)
					}
				}
				resp, err := client.Do(req)
				if err != nil {
					c.JSON(http.StatusBadGateway, gin.H{"error": "代理请求失败"})
					return
				}
				defer resp.Body.Close()
				for key, values := range resp.Header {
					for _, value := range values {
						c.Writer.Header().Add(key, value)
					}
				}
				// 复制响应状态码和内容
				c.Status(resp.StatusCode)
				c.Header("Content-Type", resp.Header.Get("Content-Type"))
				io.Copy(c.Writer, resp.Body)
			})
			
			frontendRouter.Static("/uploads", cfg.UploadDir)
			frontendRouter.NoRoute(func(c *gin.Context) {
				if !strings.HasPrefix(c.Request.URL.Path, "/api/") && !strings.HasPrefix(c.Request.URL.Path, "/uploads/") {
					c.File(filepath.Join(distPath, "index.html"))
				}
			})
			log.Printf("使用文件系统前端文件: %s", distPath)
		} else if _, err := os.Stat("./dist"); err == nil {
			frontendRouter.Static("/static", "./dist/static")
			frontendRouter.StaticFile("/favicon.ico", "./dist/favicon.ico")
			frontendRouter.Any("/api/*path", func(c *gin.Context) {
				client := &http.Client{Timeout: 30 * time.Second}
				targetURL := "http://localhost:1551" + c.Request.URL.Path
				if c.Request.URL.RawQuery != "" {
					targetURL += "?" + c.Request.URL.RawQuery
				}
				req, _ := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, targetURL, c.Request.Body)
				for key, values := range c.Request.Header {
					for _, value := range values {
						req.Header.Add(key, value)
					}
				}
				resp, err := client.Do(req)
				if err != nil {
					c.JSON(http.StatusBadGateway, gin.H{"error": "代理请求失败"})
					return
				}
				defer resp.Body.Close()
				for key, values := range resp.Header {
					for _, value := range values {
						c.Writer.Header().Add(key, value)
					}
				}
				// 复制响应状态码和内容
				c.Status(resp.StatusCode)
				c.Header("Content-Type", resp.Header.Get("Content-Type"))
				io.Copy(c.Writer, resp.Body)
			})
			frontendRouter.Static("/uploads", cfg.UploadDir)
			frontendRouter.NoRoute(func(c *gin.Context) {
				if !strings.HasPrefix(c.Request.URL.Path, "/api/") && !strings.HasPrefix(c.Request.URL.Path, "/uploads/") {
					c.File("./dist/index.html")
				}
			})
			log.Println("使用当前目录的dist文件夹（开发模式）")
		} else {
			log.Println("警告: 未找到前端文件，前端功能不可用")
		}
	}

	// 启动两个服务器
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "1551"
	}
	
	frontendPort := os.Getenv("FRONTEND_PORT")
	if frontendPort == "" {
		frontendPort = "1552"
	}

	// 创建API服务器
	apiServer := &http.Server{
		Addr:    ":" + apiPort,
		Handler:  r,
	}

	// 创建前端服务器
	frontendServer := &http.Server{
		Addr:    ":" + frontendPort,
		Handler: frontendRouter,
	}

	// 只在首次启动时显示默认账号信息
	firstRunFile := filepath.Join(dataDir, ".first_run")
	if _, err := os.Stat(firstRunFile); os.IsNotExist(err) {
		os.WriteFile(firstRunFile, []byte(""), 0644)
		log.Printf("默认管理员账号: admin, 密码: admin123")
		log.Printf("提示: 首次启动后，请及时修改默认密码以确保安全")
	}

	log.Printf("========================================")
	log.Printf("服务器启动成功！")
	log.Printf("后端API: http://localhost:%s", apiPort)
	log.Printf("前端界面: http://localhost:%s", frontendPort)
	log.Printf("========================================")

	// 在goroutine中启动前端服务器
	go func() {
		if err := frontendServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("前端服务器启动失败: %v", err)
		}
	}()

	// 在主goroutine中启动API服务器
	if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("API服务器启动失败: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
