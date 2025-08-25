-- 创建数据库
CREATE DATABASE IF NOT EXISTS middleware CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE middleware;

-- 创建采购请求表
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