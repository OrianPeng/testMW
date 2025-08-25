package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	fmt.Println("Testing Redis connection...")

	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("✅ Successfully connected to Redis")

	// 测试基本操作
	testKey := "test:connection"
	testValue := "Hello Redis!"

	// 设置值
	if err := client.Set(ctx, testKey, testValue, 0).Err(); err != nil {
		log.Fatalf("Failed to set key: %v", err)
	}
	fmt.Printf("✅ Set key: %s = %s\n", testKey, testValue)

	// 获取值
	val, err := client.Get(ctx, testKey).Result()
	if err != nil {
		log.Fatalf("Failed to get key: %v", err)
	}
	fmt.Printf("✅ Get key: %s = %s\n", testKey, val)

	// 测试哈希操作
	hashKey := "test:hash"
	pipe := client.Pipeline()
	pipe.HSet(ctx, hashKey, "field1", "value1")
	pipe.HSet(ctx, hashKey, "field2", "value2")
	pipe.HSet(ctx, hashKey, "field3", "value3")
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to set hash: %v", err)
	}
	fmt.Printf("✅ Set hash: %s\n", hashKey)

	// 获取哈希值
	hashVal, err := client.HGetAll(ctx, hashKey).Result()
	if err != nil {
		log.Fatalf("Failed to get hash: %v", err)
	}
	fmt.Printf("✅ Get hash: %s = %v\n", hashKey, hashVal)

	// 测试有序集合操作
	zsetKey := "test:zset"
	pipe = client.Pipeline()
	pipe.ZAdd(ctx, zsetKey, &redis.Z{Score: 1.0, Member: "item1"})
	pipe.ZAdd(ctx, zsetKey, &redis.Z{Score: 2.0, Member: "item2"})
	pipe.ZAdd(ctx, zsetKey, &redis.Z{Score: 3.0, Member: "item3"})
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to set zset: %v", err)
	}
	fmt.Printf("✅ Set zset: %s\n", zsetKey)

	// 获取有序集合
	zsetVal, err := client.ZRange(ctx, zsetKey, 0, -1).Result()
	if err != nil {
		log.Fatalf("Failed to get zset: %v", err)
	}
	fmt.Printf("✅ Get zset: %s = %v\n", zsetKey, zsetVal)

	// 清理测试数据
	pipe = client.Pipeline()
	pipe.Del(ctx, testKey)
	pipe.Del(ctx, hashKey)
	pipe.Del(ctx, zsetKey)
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("Warning: Failed to cleanup test data: %v", err)
	} else {
		fmt.Println("✅ Cleaned up test data")
	}

	// 关闭连接
	if err := client.Close(); err != nil {
		log.Printf("Warning: Failed to close Redis connection: %v", err)
	} else {
		fmt.Println("✅ Closed Redis connection")
	}

	fmt.Println("🎉 All Redis tests passed!")
}
