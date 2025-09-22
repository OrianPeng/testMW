package config

import (
	"time"
)

// Config 配置结构
type Config struct {
	Server ServerConfig `yaml:"server"`
	RPA    RPAConfig    `yaml:"rpa"`
	UiPath UiPathConfig `yaml:"uipath"`
	Queue  QueueConfig  `yaml:"queue"`
	Redis  RedisConfig  `yaml:"redis"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Logger LoggerConfig `yaml:"logger"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// RPAConfig RPA配置
type RPAConfig struct {
	BaseURL       string        `yaml:"base_url"`
	Timeout       time.Duration `yaml:"timeout"`
	CheckInterval time.Duration `yaml:"check_interval"`
}

// UiPathConfig UiPath配置
type UiPathConfig struct {
	OrchBaseURL string        `yaml:"orch_base_url"`
	TenancyName string        `yaml:"tenancy_name"`
	Username    string        `yaml:"username"`
	Password    string        `yaml:"password"`
	FolderID    int           `yaml:"folder_id"`
	QueueName   string        `yaml:"queue_name"`
	VerifySSL   bool          `yaml:"verify_ssl"`
	Timeout     time.Duration `yaml:"timeout"`
}

// QueueConfig 队列配置
type QueueConfig struct {
	Type          string        `yaml:"type"`
	CheckInterval time.Duration `yaml:"check_interval"`
	MaxRetries    int           `yaml:"max_retries"`
	RetryInterval time.Duration `yaml:"retry_interval"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	Database        string        `yaml:"database"`
	Charset         string        `yaml:"charset"`
	ParseTime       bool          `yaml:"parse_time"`
	Loc             string        `yaml:"loc"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
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
			Port:         8082,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		RPA: RPAConfig{
			BaseURL:       "http://localhost:8084",
			Timeout:       60 * time.Second,
			CheckInterval: 60 * time.Second,
		},
		UiPath: UiPathConfig{
			OrchBaseURL: "https://fmvrpaorchp01.sg.lan",
			TenancyName: "Default",
			Username:    "API",
			Password:    "Dash_Api_12345678",
			FolderID:    31,
			QueueName:   "CreationPR",
			VerifySSL:   false,
			Timeout:     60 * time.Second,
		},
		Queue: QueueConfig{
			Type:          "memory",
			CheckInterval: 5 * time.Second,
			MaxRetries:    3,
			RetryInterval: 10 * time.Second,
		},
		Redis: RedisConfig{
			Addr:         "localhost:6379",
			Password:     "",
			DB:           0,
			PoolSize:     10,
			MinIdleConns: 5,
		},
		MySQL: MySQLConfig{
			Host:            "localhost",
			Port:            3306,
			Username:        "root",
			Password:        "Test123ls",
			Database:        "middleware",
			Charset:         "utf8mb4",
			ParseTime:       true,
			Loc:             "Local",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 300 * time.Second,
		},
		Logger: LoggerConfig{
			Level:  "info",
			Format: "json",
		},
	}
}
