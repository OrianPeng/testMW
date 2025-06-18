package config

import (
	"time"
)

// Config 应用配置
type Config struct {
	Server ServerConfig `yaml:"server"`
	RPA    RPAConfig    `yaml:"rpa"`
	Queue  QueueConfig  `yaml:"queue"`
	Logger LoggerConfig `yaml:"logger"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// RPAConfig RPA 系统配置
type RPAConfig struct {
	BaseURL        string        `yaml:"base_url"`
	Timeout        time.Duration `yaml:"timeout"`
	CheckInterval  time.Duration `yaml:"check_interval"`
}

// QueueConfig 队列配置
type QueueConfig struct {
	CheckInterval time.Duration `yaml:"check_interval"`
	MaxRetries    int           `yaml:"max_retries"`
	RetryInterval time.Duration `yaml:"retry_interval"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		RPA: RPAConfig{
			BaseURL:       "http://localhost:8081",
			Timeout:       60 * time.Second,
			CheckInterval: 5 * time.Second,
		},
		Queue: QueueConfig{
			CheckInterval: 5 * time.Second,
			MaxRetries:    3,
			RetryInterval: 10 * time.Second,
		},
		Logger: LoggerConfig{
			Level:  "info",
			Format: "json",
		},
	}
}