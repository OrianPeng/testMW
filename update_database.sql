-- 数据库更新脚本：迁移到新的PR数据结构
-- 执行前请备份数据库

USE middleware;

-- 1. 创建新表结构
CREATE TABLE IF NOT EXISTS purchase_requests_new (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    request_id VARCHAR(100) NOT NULL UNIQUE,
    doc_type VARCHAR(50) NOT NULL DEFAULT 'PR',
    plant VARCHAR(50) NOT NULL DEFAULT '1000',
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    material VARCHAR(100) NOT NULL,
    delivery_date DATETIME NOT NULL,
    vendor_code VARCHAR(50) NOT NULL DEFAULT 'VENDOR-001',
    short_text TEXT NOT NULL,
    material_group VARCHAR(50) NOT NULL DEFAULT 'GENERAL',
    unit_type VARCHAR(20) NOT NULL DEFAULT 'EA',
    requester VARCHAR(100) NOT NULL,
    purchase_organization VARCHAR(50) NOT NULL DEFAULT 'PURCHASE-ORG-001',
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

-- 2. 迁移现有数据（如果有的话）
-- 注意：这里假设旧表存在，如果不存在可以跳过这部分
INSERT INTO purchase_requests_new (
    request_id, doc_type, plant, quantity, unit_price, material,
    delivery_date, vendor_code, short_text, material_group, unit_type,
    requester, purchase_organization, currency, total_amount,
    priority, status, urgency, approver_id, approver_name, approved_at,
    rejection_reason, comments, attachments, metadata, callback,
    retry_count, error_msg, processed_at, created_at, updated_at
)
SELECT 
    request_id,
    'PR' as doc_type,
    COALESCE(department, '1000') as plant,
    quantity,
    unit_price,
    COALESCE(item_name, 'UNKNOWN') as material,
    COALESCE(expected_date, NOW()) as delivery_date,
    'VENDOR-001' as vendor_code,
    COALESCE(item_description, 'No description') as short_text,
    'GENERAL' as material_group,
    'EA' as unit_type,
    COALESCE(requester_name, requester_id) as requester,
    'PURCHASE-ORG-001' as purchase_organization,
    'CNY' as currency,
    total_amount,
    priority,
    status,
    urgency,
    approver_id,
    approver_name,
    approved_at,
    rejection_reason,
    comments,
    attachments,
    metadata,
    '' as callback,
    0 as retry_count,
    '' as error_msg,
    NULL as processed_at,
    created_at,
    updated_at
FROM purchase_requests
WHERE NOT EXISTS (SELECT 1 FROM purchase_requests_new WHERE purchase_requests_new.request_id = purchase_requests.request_id);

-- 3. 删除旧表
DROP TABLE IF EXISTS purchase_requests;

-- 4. 重命名新表
RENAME TABLE purchase_requests_new TO purchase_requests;

-- 5. 验证表结构
DESCRIBE purchase_requests;

-- 6. 验证数据迁移
SELECT COUNT(*) as total_records FROM purchase_requests;
SELECT * FROM purchase_requests LIMIT 5; 