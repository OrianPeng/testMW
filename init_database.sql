-- 创建数据库
CREATE DATABASE IF NOT EXISTS middleware CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE middleware;

-- 创建采购请求表
CREATE TABLE IF NOT EXISTS purchase_requests (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    request_id VARCHAR(100) NOT NULL UNIQUE,
    requester_id VARCHAR(100) NOT NULL,
    requester_name VARCHAR(200) NOT NULL,
    department VARCHAR(100) NOT NULL,
    item_name VARCHAR(200) NOT NULL,
    item_description TEXT,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    priority INT DEFAULT 2,
    status VARCHAR(20) DEFAULT 'pending',
    reason TEXT NOT NULL,
    urgency VARCHAR(20) DEFAULT 'normal',
    expected_date DATETIME NULL,
    approver_id VARCHAR(100) NULL,
    approver_name VARCHAR(200) NULL,
    approved_at DATETIME NULL,
    rejection_reason TEXT NULL,
    comments TEXT,
    attachments TEXT,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_request_id (request_id),
    INDEX idx_requester_id (requester_id),
    INDEX idx_department (department),
    INDEX idx_status (status),
    INDEX idx_priority (priority),
    INDEX idx_created_at (created_at),
    INDEX idx_status_priority (status, priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入一些测试数据
INSERT INTO purchase_requests (
    request_id, requester_id, requester_name, department,
    item_name, item_description, quantity, unit_price, total_amount,
    priority, status, reason, urgency, expected_date, comments
) VALUES 
('PR-TEST-001', 'USER-001', '张三', '技术部', '笔记本电脑', '高性能开发用笔记本电脑', 2, 5999.00, 11998.00, 3, 'pending', '开发团队需要新的开发设备', 'normal', DATE_ADD(NOW(), INTERVAL 7 DAY), '请尽快处理'),
('PR-TEST-002', 'USER-002', '李四', '市场部', '投影仪', '会议室用投影仪', 1, 2999.00, 2999.00, 2, 'pending', '会议室设备更新', 'urgent', DATE_ADD(NOW(), INTERVAL 3 DAY), '紧急需求'),
('PR-TEST-003', 'USER-003', '王五', '人事部', '办公椅', '人体工学办公椅', 5, 800.00, 4000.00, 1, 'approved', '改善员工办公环境', 'normal', DATE_ADD(NOW(), INTERVAL 14 DAY), '已批准'),
('PR-TEST-004', 'USER-004', '赵六', '财务部', '打印机', '彩色激光打印机', 1, 2500.00, 2500.00, 2, 'rejected', '财务部门打印需求', 'normal', DATE_ADD(NOW(), INTERVAL 5 DAY), '预算不足');

-- 显示表结构
DESCRIBE purchase_requests;

-- 显示测试数据
SELECT * FROM purchase_requests; 