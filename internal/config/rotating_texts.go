package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// RotatingTextsConfig 轮换文本配置
type RotatingTextsConfig struct {
	Texts []string `json:"texts"`
}

// rotatingTextsConfigPath 返回轮换文本配置文件路径
func (c *Config) rotatingTextsConfigPath() string {
	return filepath.Join(c.DataDir, "rotating_texts.json")
}

// LoadRotatingTexts 加载轮换文本配置
func (c *Config) LoadRotatingTexts() (*RotatingTextsConfig, error) {
	path := c.rotatingTextsConfigPath()
	
	// 如果文件不存在，返回默认值
	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultTexts := &RotatingTextsConfig{
			Texts: []string{
				"你好鸭，欢迎来到我的主页！！",
				"随时可以联系我，期待与你交流。",
				"愿你历尽千帆，归来仍是少年。",
				"梦想还是要有的，万一实现了呢？",
				"I hope you have a happy day every day.",
			},
		}
		return defaultTexts, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg RotatingTextsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// 确保至少有默认文本
	if len(cfg.Texts) == 0 {
		cfg.Texts = []string{
			"你好鸭，欢迎来到我的主页！！",
			"随时可以联系我，期待与你交流。",
			"愿你历尽千帆，归来仍是少年。",
			"梦想还是要有的，万一实现了呢？",
			"I hope you have a happy day every day.",
		}
	}

	return &cfg, nil
}

// SaveRotatingTexts 保存轮换文本配置
func (c *Config) SaveRotatingTexts(cfg *RotatingTextsConfig) error {
	path := c.rotatingTextsConfigPath()
	
	// 限制最多8条文本
	if len(cfg.Texts) > 8 {
		cfg.Texts = cfg.Texts[:8]
	}
	
	// 过滤空文本
	filteredTexts := make([]string, 0, len(cfg.Texts))
	for _, text := range cfg.Texts {
		if text != "" {
			filteredTexts = append(filteredTexts, text)
		}
	}
	cfg.Texts = filteredTexts

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
