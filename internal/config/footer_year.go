package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// FooterYearConfig 用于存储底部版权年份配置
type FooterYearConfig struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// footerYearConfigPath 返回年份配置文件路径
func (c *Config) footerYearConfigPath() string {
	return filepath.Join(c.DataDir, "footer_year.json")
}

// LoadFooterYear 读取底部年份配置（文件不存在时返回空配置而不是报错）
func (c *Config) LoadFooterYear() (*FooterYearConfig, error) {
	path := c.footerYearConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 没有配置文件时返回空配置
			return &FooterYearConfig{}, nil
		}
		return nil, err
	}

	var cfg FooterYearConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		// 解析失败时返回空配置，但把错误也返回方便日志排查
		return &FooterYearConfig{}, err
	}
	return &cfg, nil
}

// SaveFooterYear 保存底部年份配置
func (c *Config) SaveFooterYear(cfg *FooterYearConfig) error {
	if cfg == nil {
		cfg = &FooterYearConfig{}
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	path := c.footerYearConfigPath()
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

