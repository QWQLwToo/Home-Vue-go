package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DataDir      string
	DatabasePath string
	UploadDir    string
	JWTSecret    string
}

func New(dataDir string) *Config {
	// 确保data目录存在
	os.MkdirAll(dataDir, 0755)

	// 创建上传目录
	uploadDir := filepath.Join(dataDir, "uploads")
	os.MkdirAll(uploadDir, 0755)

	// 从环境变量获取JWT密钥，如果没有则使用默认值
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	return &Config{
		DataDir:      dataDir,
		DatabasePath: filepath.Join(dataDir, "home.db"),
		UploadDir:    uploadDir,
		JWTSecret:    jwtSecret,
	}
}
