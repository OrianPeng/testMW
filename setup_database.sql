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

-- 创建采购订单表
CREATE TABLE IF NOT EXISTS purchase_orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    po_number VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) DEFAULT 'draft',
    supplier_name VARCHAR(200) NOT NULL,
    supplier_code VARCHAR(50) NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'CNY',
    created_by VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sent_at DATETIME NULL,
    confirmed_at DATETIME NULL,
    notes TEXT,
    INDEX idx_po_number (po_number),
    INDEX idx_status (status),
    INDEX idx_supplier_code (supplier_code),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建供应商表
CREATE TABLE IF NOT EXISTS suppliers (
    id VARCHAR(50) PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(20) DEFAULT 'manufacturer',
    status VARCHAR(20) DEFAULT 'active',
    contact_person VARCHAR(100),
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    categories JSON,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_code (code),
    INDEX idx_name (name),
    INDEX idx_type (type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建流程监控表
CREATE TABLE IF NOT EXISTS process_monitor (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    request_id VARCHAR(100) NOT NULL,
    process_type VARCHAR(10) NOT NULL,
    current_step VARCHAR(100) NOT NULL,
    responsible VARCHAR(100) NOT NULL,
    department VARCHAR(100) NOT NULL,
    blocked_duration VARCHAR(50),
    priority INT DEFAULT 2,
    reason TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_request_id (request_id),
    INDEX idx_process_type (process_type),
    INDEX idx_status (status),
    INDEX idx_priority (priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建文件上传表
CREATE TABLE IF NOT EXISTS file_uploads (
    id VARCHAR(50) PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size BIGINT NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    uploader VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_file_type (file_type),
    INDEX idx_uploader (uploader),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建AI对话表
CREATE TABLE IF NOT EXISTS ai_conversations (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    department VARCHAR(100),
    context JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建AI消息表
CREATE TABLE IF NOT EXISTS ai_messages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    conversation_id VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    response TEXT,
    response_type VARCHAR(20),
    suggestions JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_conversation_id (conversation_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (conversation_id) REFERENCES ai_conversations(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
