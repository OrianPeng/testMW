package main

import (
	"database/sql"
	"fmt"
	"rpa-middleware/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 加载配置
	cfg := config.DefaultConfig()

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.MySQL.Username,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database,
		cfg.MySQL.Charset,
		cfg.MySQL.ParseTime,
		cfg.MySQL.Loc,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Failed to connect to MySQL: %v\n", err)
		return
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		fmt.Printf("Failed to ping MySQL: %v\n", err)
		return
	}

	fmt.Println("Connected to MySQL successfully")

	// 添加updated_at字段到purchase_orders表
	alterQuery := `
	ALTER TABLE purchase_orders 
	ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	`

	_, err = db.Exec(alterQuery)
	if err != nil {
		fmt.Printf("Error adding updated_at column: %v\n", err)
	} else {
		fmt.Println("updated_at column added successfully to purchase_orders table")
	}

	// 添加updated_at字段到suppliers表
	alterQuery2 := `
	ALTER TABLE suppliers 
	ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	`

	_, err = db.Exec(alterQuery2)
	if err != nil {
		fmt.Printf("Error adding updated_at column to suppliers: %v\n", err)
	} else {
		fmt.Println("updated_at column added successfully to suppliers table")
	}

	fmt.Println("Database fixes completed!")
}
