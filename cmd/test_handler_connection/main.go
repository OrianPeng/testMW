package main

import (
	"context"
	"database/sql"
	"fmt"
	"rpa-middleware/internal/api"
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

	// 创建Handler
	handler := api.NewPurchaseRequestHandler(repo, logger)

	// 测试查询
	query := &models.PurchaseRequestQuery{
		Page:     1,
		PageSize: 10,
	}

	fmt.Println("Testing Repository List method...")
	result, err := repo.List(context.Background(), query)
	if err != nil {
		fmt.Printf("Repository List error: %v\n", err)
		return
	}

	fmt.Printf("Repository List successful! Found %d requests\n", len(result.Requests))

	// 测试Handler的数据库连接
	fmt.Println("Testing Handler database connection...")

	// 模拟一个简单的查询来测试Handler的数据库连接
	testQuery := &models.PurchaseRequestQuery{
		Page:     1,
		PageSize: 5,
	}

	handlerResult, err := repo.List(context.Background(), testQuery)
	if err != nil {
		fmt.Printf("Handler database connection error: %v\n", err)
		return
	}

	fmt.Printf("Handler database connection successful! Found %d requests\n", len(handlerResult.Requests))
	fmt.Printf("Handler created successfully: %v\n", handler != nil)
}
