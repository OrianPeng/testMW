package main

import (
	"context"
	"database/sql"
	"fmt"
	"rpa-middleware/internal/config"
	"rpa-middleware/internal/models"
	"rpa-middleware/internal/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化日志
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

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

	// 创建Repository
	repo := repository.NewPurchaseRequestRepository(db, logger)

	// 测试查询
	query := &models.PurchaseRequestQuery{
		Page:     1,
		PageSize: 10,
	}

	result, err := repo.List(context.Background(), query)
	if err != nil {
		fmt.Printf("Error querying purchase requests: %v\n", err)
		return
	}

	fmt.Printf("Query successful! Found %d requests\n", len(result.Requests))
	for i, req := range result.Requests {
		if i < 3 { // 只显示前3条
			fmt.Printf("Request %d: ID=%d, RequestID=%s, Status=%s\n",
				i+1, req.ID, req.RequestID, req.Status)
		}
	}
}
