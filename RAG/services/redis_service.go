package services

import (
	"RAG/utils"
	"context"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis 初始化 Redis 客户端连接
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     utils.Config.GetString("redis.addr"),
		Password: utils.Config.GetString("redis.password"),
		DB:       utils.Config.GetInt("redis.db"),
	})

	// 检查连接
	err := RedisClient.Ping(context.Background()).Err()
	if err != nil {
		utils.Logger.Fatal("Failed to connect to Redis:", err)
	}
	utils.Logger.Println("Connected to Redis")
}

// CloseRedisClient 关闭 Redis 客户端连接
func CloseRedisClient() {
	err := RemoveOnlineUsers()
	if err != nil {
		utils.Logger.Println("Error removing online users:", err)
	}
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			utils.Logger.Printf("Failed to close Redis connection: %v", err)
			return
		}
		utils.Logger.Println("Redis connection closed.")
	} else {
		utils.Logger.Println("Redis client is not initialized.")
	}
}

func AddOnlineUser(userID string) error {
	return RedisClient.Set(ctx, "user:"+userID, "online", 0).Err()
}

func RemoveOnlineUser(userID string) error {
	return RedisClient.Del(ctx, "user:"+userID).Err()
}
func RemoveOnlineUsers() error {
	userIDs, err := GetOnlineUsers()
	if err != nil {
		utils.Logger.Println("RemoveOnlineUsers Function :Error getting online users:", err)
		return err
	}
	for _, userID := range userIDs {
		if err := RemoveOnlineUser(userID); err != nil {
			return err
		}
	}
	utils.Logger.Println("All online users removed successfully.")
	return nil
}
func GetOnlineUsers() ([]string, error) {
	keys, err := RedisClient.Keys(ctx, "user:*").Result()
	if err != nil {
		utils.Logger.Println("Error getting online users:", err)
		return nil, err
	}
	// Extract user IDs from keys
	var users []string
	for _, key := range keys {
		users = append(users, key[5:]) // Remove "user:" prefix
	}
	return users, nil
}
