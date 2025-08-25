package database

import (
	"database/sql"
	"fmt"
	"rpa-middleware/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// MySQLManager MySQL数据库管理器
type MySQLManager struct {
	db     *sql.DB
	config *config.MySQLConfig
	logger *logrus.Logger
}

// NewMySQLManager 创建新的MySQL管理器
func NewMySQLManager(cfg *config.MySQLConfig, logger *logrus.Logger) (*MySQLManager, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("MySQL database connected successfully")

	return &MySQLManager{
		db:     db,
		config: cfg,
		logger: logger,
	}, nil
}

// GetDB 获取数据库连接
func (m *MySQLManager) GetDB() *sql.DB {
	return m.db
}

// Close 关闭数据库连接
func (m *MySQLManager) Close() error {
	if m.db != nil {
		m.logger.Info("Closing MySQL database connection")
		return m.db.Close()
	}
	return nil
}

// InitTables 初始化数据库表
func (m *MySQLManager) InitTables() error {
	m.logger.Info("Initializing database tables")

	// 创建采购请求表
	createPurchaseRequestTable := `
	CREATE TABLE IF NOT EXISTS purchase_requests (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		request_id VARCHAR(100) NOT NULL UNIQUE,
		doc_type VARCHAR(50) NOT NULL,
		plant VARCHAR(50) NOT NULL,
		quantity INT NOT NULL,
		unit_price DECIMAL(10,2) NOT NULL,
		material VARCHAR(100) NOT NULL,
		delivery_date DATETIME NOT NULL,
		vendor_code VARCHAR(50) NOT NULL,
		short_text TEXT NOT NULL,
		material_group VARCHAR(50) NOT NULL,
		unit_type VARCHAR(20) NOT NULL,
		requester VARCHAR(100) NOT NULL,
		purchase_organization VARCHAR(50) NOT NULL,
		currency VARCHAR(10) NOT NULL DEFAULT 'CNY',
		total_amount DECIMAL(10,2) NOT NULL,
		priority INT DEFAULT 2,
		status VARCHAR(20) DEFAULT 'pending',
		urgency VARCHAR(20) DEFAULT 'normal',
		approver_id VARCHAR(100) NULL,
		approver_name VARCHAR(200) NULL,
		approved_at DATETIME NULL,
		rejection_reason TEXT NULL,
		comments TEXT,
		attachments TEXT,
		metadata JSON,
		callback VARCHAR(500),
		retry_count INT DEFAULT 0,
		error_msg TEXT,
		processed_at DATETIME NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_request_id (request_id),
		INDEX idx_requester (requester),
		INDEX idx_plant (plant),
		INDEX idx_status (status),
		INDEX idx_priority (priority),
		INDEX idx_created_at (created_at),
		INDEX idx_status_priority (status, priority),
		INDEX idx_vendor_code (vendor_code),
		INDEX idx_material_group (material_group)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err := m.db.Exec(createPurchaseRequestTable)
	if err != nil {
		return fmt.Errorf("failed to create purchase_requests table: %w", err)
	}

	m.logger.Info("Database tables initialized successfully")
	return nil
}

// HealthCheck 健康检查
func (m *MySQLManager) HealthCheck() error {
	return m.db.Ping()
}
