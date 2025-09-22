package main

import (
	"database/sql"
	"fmt"
	"log"
	"rpa-middleware/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 加载配置
	cfg := config.DefaultConfig()

	// 连接数据库（不指定数据库名）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=%t&loc=%s",
		cfg.MySQL.Username,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Charset,
		cfg.MySQL.ParseTime,
		cfg.MySQL.Loc,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}

	fmt.Println("Connected to MySQL successfully")

	// 创建数据库
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.MySQL.Database))
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	fmt.Printf("Database %s created successfully\n", cfg.MySQL.Database)

	// 使用数据库
	_, err = db.Exec(fmt.Sprintf("USE %s", cfg.MySQL.Database))
	if err != nil {
		log.Fatalf("Failed to use database: %v", err)
	}

	// 创建表结构
	fmt.Println("Creating tables...")

	// 创建采购请求表
	createPurchaseRequestsTable := `
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

	_, err = db.Exec(createPurchaseRequestsTable)
	if err != nil {
		fmt.Printf("Error creating purchase_requests table: %v\n", err)
	} else {
		fmt.Println("purchase_requests table created successfully")
	}

	// 创建采购订单表
	createPurchaseOrdersTable := `
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
	`

	_, err = db.Exec(createPurchaseOrdersTable)
	if err != nil {
		fmt.Printf("Error creating purchase_orders table: %v\n", err)
	} else {
		fmt.Println("purchase_orders table created successfully")
	}

	// 创建供应商表
	createSuppliersTable := `
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
	`

	_, err = db.Exec(createSuppliersTable)
	if err != nil {
		fmt.Printf("Error creating suppliers table: %v\n", err)
	} else {
		fmt.Println("suppliers table created successfully")
	}

	// 创建流程监控表
	createProcessMonitorTable := `
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
	`

	_, err = db.Exec(createProcessMonitorTable)
	if err != nil {
		fmt.Printf("Error creating process_monitor table: %v\n", err)
	} else {
		fmt.Println("process_monitor table created successfully")
	}

	// 创建文件上传表
	createFileUploadsTable := `
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
	`

	_, err = db.Exec(createFileUploadsTable)
	if err != nil {
		fmt.Printf("Error creating file_uploads table: %v\n", err)
	} else {
		fmt.Println("file_uploads table created successfully")
	}

	// 创建AI对话表
	createAIConversationsTable := `
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
	`

	_, err = db.Exec(createAIConversationsTable)
	if err != nil {
		fmt.Printf("Error creating ai_conversations table: %v\n", err)
	} else {
		fmt.Println("ai_conversations table created successfully")
	}

	// 创建AI消息表
	createAIMessagesTable := `
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
	`

	_, err = db.Exec(createAIMessagesTable)
	if err != nil {
		fmt.Printf("Error creating ai_messages table: %v\n", err)
	} else {
		fmt.Println("ai_messages table created successfully")
	}

	// 插入模拟数据
	fmt.Println("Inserting mock data...")

	// 插入采购请求数据
	insertPurchaseRequests := `
	INSERT IGNORE INTO purchase_requests (
		request_id, doc_type, plant, quantity, unit_price, material, delivery_date,
		vendor_code, short_text, material_group, unit_type, requester,
		purchase_organization, currency, total_amount, priority, status, urgency,
		comments, retry_count, created_at, updated_at
	) VALUES
	('PR-20231201-001', 'NB', '1000', 100, 25.50, 'MAT001', '2023-12-15 00:00:00',
	 'VENDOR001', 'Office Supplies Procurement', 'MG001', 'EA', 'John Smith',
	 'PO001', 'CNY', 2550.00, 2, 'pending', 'normal',
	 'Urgent office supplies needed', 0, NOW(), NOW()),
	('PR-20231201-002', 'NB', '1000', 50, 100.00, 'MAT002', '2023-12-20 00:00:00',
	 'VENDOR002', 'Computer Equipment', 'MG002', 'EA', 'Alice Johnson',
	 'PO001', 'CNY', 5000.00, 3, 'approved', 'urgent',
	 'New computer equipment for office', 0, NOW(), NOW()),
	('PR-20231201-003', 'UB', '2000', 200, 15.75, 'MAT003', '2023-12-25 00:00:00',
	 'VENDOR003', 'Raw Materials', 'MG003', 'KG', 'Bob Wilson',
	 'PO002', 'CNY', 3150.00, 1, 'processing', 'normal',
	 'Raw materials for production', 0, NOW(), NOW());
	`

	_, err = db.Exec(insertPurchaseRequests)
	if err != nil {
		fmt.Printf("Error inserting purchase requests: %v\n", err)
	} else {
		fmt.Println("Purchase requests data inserted successfully")
	}

	// 插入采购订单数据
	insertPurchaseOrders := `
	INSERT IGNORE INTO purchase_orders (
		po_number, status, supplier_name, supplier_code, total_amount, currency,
		created_by, created_at, sent_at, confirmed_at, notes
	) VALUES
	('PO-20231201-001', 'sent', 'ABC Manufacturing Co., Ltd.', 'SUP001', 50000.00, 'CNY',
	 'John Smith', NOW(), NOW(), NULL, 'Urgent procurement'),
	('PO-20231201-002', 'confirmed', 'XYZ Electronics Inc.', 'SUP002', 35000.00, 'CNY',
	 'Alice Johnson', NOW(), NOW(), NOW(), 'Electronics equipment order'),
	('PO-20231201-003', 'delivered', 'DEF Materials Corp.', 'SUP003', 20000.00, 'CNY',
	 'Bob Wilson', NOW(), NOW(), NOW(), 'Raw materials delivered');
	`

	_, err = db.Exec(insertPurchaseOrders)
	if err != nil {
		fmt.Printf("Error inserting purchase orders: %v\n", err)
	} else {
		fmt.Println("Purchase orders data inserted successfully")
	}

	// 插入供应商数据
	insertSuppliers := `
	INSERT IGNORE INTO suppliers (
		id, code, name, type, status, contact_person, phone, email, address,
		categories, notes, created_at
	) VALUES
	('SUP_001', 'SUP001', 'ABC Manufacturing Co., Ltd.', 'manufacturer', 'active',
	 'John Smith', '13800138001', 'john.smith@abc.com', 'No.123, Industrial Road, Chaoyang District, Beijing',
	 '["Raw Materials", "Equipment"]', 'Long-term partner', NOW()),
	('SUP_002', 'SUP002', 'XYZ Electronics Inc.', 'manufacturer', 'active',
	 'Alice Johnson', '13800138002', 'alice.johnson@xyz.com', 'No.456, Tech Street, Haidian District, Beijing',
	 '["Electronics", "Computer Equipment"]', 'Reliable electronics supplier', NOW()),
	('SUP_003', 'SUP003', 'DEF Materials Corp.', 'distributor', 'active',
	 'Bob Wilson', '13800138003', 'bob.wilson@def.com', 'No.789, Material Avenue, Fengtai District, Beijing',
	 '["Raw Materials", "Chemicals"]', 'Material distribution specialist', NOW());
	`

	_, err = db.Exec(insertSuppliers)
	if err != nil {
		fmt.Printf("Error inserting suppliers: %v\n", err)
	} else {
		fmt.Println("Suppliers data inserted successfully")
	}

	fmt.Println("Database setup completed!")
}
