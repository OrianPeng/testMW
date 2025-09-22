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

	// 查询采购请求表
	rows, err := db.Query("SELECT COUNT(*) FROM purchase_requests")
	if err != nil {
		fmt.Printf("Error querying purchase_requests: %v\n", err)
		return
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			fmt.Printf("Error scanning count: %v\n", err)
			return
		}
	}

	fmt.Printf("Purchase requests count: %d\n", count)

	// 查询前几条记录
	rows2, err := db.Query("SELECT id, request_id, status FROM purchase_requests LIMIT 5")
	if err != nil {
		fmt.Printf("Error querying purchase_requests details: %v\n", err)
		return
	}
	defer rows2.Close()

	fmt.Println("Purchase requests:")
	for rows2.Next() {
		var id int64
		var requestID, status string
		err = rows2.Scan(&id, &requestID, &status)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}
		fmt.Printf("ID: %d, RequestID: %s, Status: %s\n", id, requestID, status)
	}
}
