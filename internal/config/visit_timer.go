package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// VisitTimerConfig 用于存储访问时间显示配置
type VisitTimerConfig struct {
	ShowVisitTimer bool `json:"showVisitTimer"`
}

// visitTimerConfigPath 返回访问时间配置文件路径
func (c *Config) visitTimerConfigPath() string {
	return filepath.Join(c.DataDir, "visit_timer.json")
}

// LoadVisitTimer 读取访问时间显示配置（文件不存在时返回默认配置）
func (c *Config) LoadVisitTimer() (*VisitTimerConfig, error) {
	path := c.visitTimerConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 没有配置文件时返回默认配置（显示）
			return &VisitTimerConfig{ShowVisitTimer: true}, nil
		}
		return nil, err
	}

	var cfg VisitTimerConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		// 解析失败时返回默认配置
		return &VisitTimerConfig{ShowVisitTimer: true}, err
	}
	return &cfg, nil
}

// SaveVisitTimer 保存访问时间显示配置
func (c *Config) SaveVisitTimer(cfg *VisitTimerConfig) error {
	if cfg == nil {
		cfg = &VisitTimerConfig{ShowVisitTimer: true}
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	path := c.visitTimerConfigPath()
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}
